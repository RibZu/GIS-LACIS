param(
    [string]$Destino = 'H:\GIS-LACIS-Deploy',
    [string]$DesarrollosPrueba = '',
    [string]$UsuariosPrueba = '',
    [string]$IntegrantesPrueba = '',
    [string]$ExcluirArchivos = '',
    [switch]$Reemplazar
)

$ErrorActionPreference = 'Stop'
$utf8SinBom = New-Object System.Text.UTF8Encoding($false)
[Console]::OutputEncoding = $utf8SinBom
$OutputEncoding = $utf8SinBom
$env:PGCLIENTENCODING = 'UTF8'

$Raiz = Split-Path -Parent $PSScriptRoot
$BaseOrigen = 'lacis'
$BaseTemporal = 'lacis_deploy'
$HashAdmin = '$2a$10$0GSeTb6bSZBT1/kZbrf41.giDbz3GaU2ABeL20lew0NdX5dDw2cSK'
$Temporales = Join-Path $env:TEMP 'gis-lacis-deploy'
$Tablas = @('integrante', 'desarrollo', 'desarrollo_integrantes', 'participante_externo', 'usuario_gestor', 'tesis', 'proyecto', 'reconocimientos', 'colaboradores')
$ColumnasArchivo = @(@('integrante', 'imagen_url'), @('integrante', 'cv_url'), @('tesis', 'archivo_pdf'), @('colaboradores', 'logo_url'))
$Conexion = @('-h', 'localhost', '-p', '5433', '-U', 'postgres')
$Estatico = Join-Path $Raiz 'ui\static'
$Advertencias = New-Object System.Collections.Generic.List[string]

function Fallar([string]$mensaje) {
    Write-Host "ERROR: $mensaje" -ForegroundColor Red
    exit 1
}

function Paso([string]$mensaje) {
    Write-Host ''
    Write-Host "== $mensaje" -ForegroundColor Cyan
}

function Convertir-Ids([string]$texto, [string]$nombre) {
    $ids = @()
    foreach ($parte in ($texto -split '[,\s]+')) {
        if ($parte -eq '') { continue }
        $n = 0
        if (-not [int]::TryParse($parte, [ref]$n) -or $n -le 0) {
            Fallar "El valor '$parte' de $nombre no es un ID válido."
        }
        $ids += $n
    }
    return ,$ids
}

function Lista-Sql([int[]]$ids) {
    if ($ids.Count -eq 0) { return 'ARRAY[]::int[]' }
    return 'ARRAY[' + ($ids -join ',') + ']::int[]'
}

$carpetaPg = Get-ChildItem 'C:\Program Files\PostgreSQL' -Directory -ErrorAction SilentlyContinue |
    Where-Object { Test-Path (Join-Path $_.FullName 'bin\psql.exe') } |
    Sort-Object { [int]($_.Name -replace '\D', '0') } |
    Select-Object -Last 1
if (-not $carpetaPg) { Fallar 'No se encontró PostgreSQL en C:\Program Files\PostgreSQL.' }
$Bin = Join-Path $carpetaPg.FullName 'bin'

function Consultar([string]$base, [string]$sql) {
    $filas = & (Join-Path $Bin 'psql.exe') -X -q -A -t -F "`t" @Conexion -d $base -c $sql
    if ($LASTEXITCODE -ne 0) { Fallar "La consulta falló en la base $base." }
    return @($filas | Where-Object { $_ -and "$_".Trim() -ne '' })
}

function Ejecutar-Archivo([string]$base, [string]$ruta) {
    & (Join-Path $Bin 'psql.exe') -X -q -v ON_ERROR_STOP=1 -1 @Conexion -d $base -f $ruta | Out-Null
    if ($LASTEXITCODE -ne 0) { Fallar "psql no pudo ejecutar $ruta" }
}

function Contar-Tablas([string]$base) {
    $sql = ($Tablas | ForEach-Object { "SELECT '$_', count(*) FROM $_" }) -join ' UNION ALL '
    $resultado = [ordered]@{}
    foreach ($fila in (Consultar $base $sql)) {
        $partes = "$fila" -split "`t"
        $resultado[$partes[0]] = [int]$partes[1]
    }
    return $resultado
}

function Resolver-Exacto([string]$relativo) {
    $actual = $Estatico
    $partes = @()
    foreach ($parte in ($relativo -split '/')) {
        if ($parte -eq '') { continue }
        if (-not (Test-Path -LiteralPath $actual -PathType Container)) { return $null }
        $hijos = @(Get-ChildItem -LiteralPath $actual -Force)
        $hijo = $hijos | Where-Object { $_.Name -ceq $parte } | Select-Object -First 1
        if (-not $hijo) {
            $forma = $parte.Normalize([Text.NormalizationForm]::FormC).ToLowerInvariant()
            $hijo = $hijos | Where-Object { $_.Name.Normalize([Text.NormalizationForm]::FormC).ToLowerInvariant() -eq $forma } | Select-Object -First 1
        }
        if (-not $hijo) { return $null }
        $partes += $hijo.Name
        $actual = $hijo.FullName
    }
    return ($partes -join '/')
}

function Copiar-Carpeta([string]$origenRelativo, [string]$destinoRelativo, [string[]]$excluir) {
    $origen = Join-Path $Raiz $origenRelativo
    $destino = Join-Path $script:DestinoCompleto $destinoRelativo
    $opciones = @($origen, $destino, '/E', '/NFL', '/NDL', '/NJH', '/NJS', '/NP')
    if ($excluir.Count -gt 0) {
        $opciones += '/XF'
        $opciones += $excluir
    }
    & robocopy @opciones | Out-Null
    if ($LASTEXITCODE -ge 8) { Fallar "robocopy falló al copiar $origenRelativo (código $LASTEXITCODE)." }
}

$desarrollos = Convertir-Ids $DesarrollosPrueba 'DesarrollosPrueba'
$usuarios = Convertir-Ids $UsuariosPrueba 'UsuariosPrueba'
$integrantes = Convertir-Ids $IntegrantesPrueba 'IntegrantesPrueba'
$excluir = @($ExcluirArchivos -split ',' | ForEach-Object { $_.Trim() } | Where-Object { $_ })

Paso 'Comprobando requisitos'
& (Join-Path $Bin 'pg_isready.exe') -h localhost -p 5433 | Out-Null
if ($LASTEXITCODE -ne 0) { Fallar 'La base de datos local no responde en el puerto 5433. Levantala con iniciar.bat y probá de nuevo.' }

$script:DestinoCompleto = [IO.Path]::GetFullPath($Destino).TrimEnd('\')
$raizCompleta = [IO.Path]::GetFullPath($Raiz).TrimEnd('\')
if ($script:DestinoCompleto.Length -le 3) { Fallar 'El destino no puede ser la raíz de una unidad.' }
if ($script:DestinoCompleto -eq $raizCompleta -or $script:DestinoCompleto.StartsWith($raizCompleta + '\') -or $raizCompleta.StartsWith($script:DestinoCompleto + '\')) {
    Fallar 'El destino no puede estar dentro del proyecto ni contener al proyecto.'
}
if ((Test-Path -LiteralPath $script:DestinoCompleto) -and -not $Reemplazar) {
    Fallar "La carpeta $($script:DestinoCompleto) ya existe. Usá -Reemplazar para rehacerla."
}

Write-Host "Proyecto original : $raizCompleta"
Write-Host "Copia para deploy : $($script:DestinoCompleto)"
Write-Host "Se borran de la exportación: desarrollos [$($desarrollos -join ', ')], usuarios [$($usuarios -join ', ')], integrantes [$($integrantes -join ', ')]"

$cantidadesOrigen = Contar-Tablas $BaseOrigen

New-Item -ItemType Directory -Force $Temporales | Out-Null
$volcadoOrigen = Join-Path $Temporales 'origen.sql'
$volcadoFinal = Join-Path $Temporales 'final.sql'
$sqlLimpieza = Join-Path $Temporales 'limpieza.sql'
$sqlReferencias = Join-Path $Temporales 'referencias.sql'

try {
    Paso "Copiando la base '$BaseOrigen' a una base temporal '$BaseTemporal' (la original no se modifica)"
    & (Join-Path $Bin 'dropdb.exe') --if-exists @Conexion $BaseTemporal
    if ($LASTEXITCODE -ne 0) { Fallar 'No se pudo preparar la base temporal.' }
    & (Join-Path $Bin 'pg_dump.exe') @Conexion -d $BaseOrigen --no-owner --no-privileges -f $volcadoOrigen
    if ($LASTEXITCODE -ne 0) { Fallar 'pg_dump falló sobre la base original.' }
    & (Join-Path $Bin 'createdb.exe') @Conexion $BaseTemporal
    if ($LASTEXITCODE -ne 0) { Fallar 'No se pudo crear la base temporal.' }
    & (Join-Path $Bin 'psql.exe') -X -q -v ON_ERROR_STOP=1 @Conexion -d $BaseTemporal -f $volcadoOrigen | Out-Null
    if ($LASTEXITCODE -ne 0) { Fallar 'No se pudo cargar la copia en la base temporal.' }

    Paso 'Limpiando la copia: datos de prueba, acentos rotos y contraseña de admin'
    $plantilla = @'
DELETE FROM desarrollo WHERE id = ANY({{DES}});
DELETE FROM integrante WHERE id = ANY({{INT}});
DELETE FROM usuario_gestor WHERE id = ANY({{USU}}) AND username <> 'admin';

CREATE FUNCTION pg_temp.arreglar_acentos(t text) RETURNS text LANGUAGE plpgsql AS $f$
DECLARE
  bytes bytea := ''::bytea;
  cp integer;
  b integer;
BEGIN
  FOR i IN 1..char_length(t) LOOP
    cp := ascii(substr(t, i, 1));
    IF cp < 256 THEN
      b := cp;
    ELSE
      b := CASE cp
        WHEN 8364 THEN 128 WHEN 8218 THEN 130 WHEN 402 THEN 131 WHEN 8222 THEN 132 WHEN 8230 THEN 133
        WHEN 8224 THEN 134 WHEN 8225 THEN 135 WHEN 710 THEN 136 WHEN 8240 THEN 137 WHEN 352 THEN 138
        WHEN 8249 THEN 139 WHEN 338 THEN 140 WHEN 381 THEN 142 WHEN 8216 THEN 145 WHEN 8217 THEN 146
        WHEN 8220 THEN 147 WHEN 8221 THEN 148 WHEN 8226 THEN 149 WHEN 8211 THEN 150 WHEN 8212 THEN 151
        WHEN 732 THEN 152 WHEN 8482 THEN 153 WHEN 353 THEN 154 WHEN 8250 THEN 155 WHEN 339 THEN 156
        WHEN 382 THEN 158 WHEN 376 THEN 159 ELSE NULL END;
      IF b IS NULL THEN
        RETURN t;
      END IF;
    END IF;
    bytes := bytes || set_byte('\x00'::bytea, 0, b);
  END LOOP;
  RETURN convert_from(bytes, 'UTF8');
EXCEPTION WHEN OTHERS THEN
  RETURN t;
END
$f$;

DO $do$
DECLARE
  r record;
  patron text := '[' || chr(195) || chr(194) || ']';
BEGIN
  FOR r IN SELECT table_name, column_name FROM information_schema.columns
           WHERE table_schema = 'public' AND data_type IN ('text', 'character varying', 'character')
             AND column_name <> 'password_hash' LOOP
    EXECUTE format('UPDATE %I SET %I = pg_temp.arreglar_acentos(%I) WHERE %I ~ %L',
                   r.table_name, r.column_name, r.column_name, r.column_name, patron);
  END LOOP;
END
$do$;

UPDATE usuario_gestor SET password_hash = '{{HASH}}' WHERE username = 'admin';

DO $v$
DECLARE
  r record;
  restantes integer;
  patron text := '[' || chr(195) || chr(194) || ']';
BEGIN
  FOR r IN SELECT table_name, column_name FROM information_schema.columns
           WHERE table_schema = 'public' AND data_type IN ('text', 'character varying', 'character')
             AND column_name <> 'password_hash' LOOP
    EXECUTE format('SELECT count(*) FROM %I WHERE %I ~ %L', r.table_name, r.column_name, patron) INTO restantes;
    IF restantes > 0 THEN
      RAISE EXCEPTION 'Quedan % celdas con acentos rotos en %.%', restantes, r.table_name, r.column_name;
    END IF;
  END LOOP;
  IF NOT EXISTS (SELECT 1 FROM usuario_gestor WHERE username = 'admin' AND rol = 'ADMIN') THEN
    RAISE EXCEPTION 'No existe el usuario admin con rol ADMIN';
  END IF;
  IF EXISTS (SELECT 1 FROM desarrollo WHERE id = ANY({{DES}}))
     OR EXISTS (SELECT 1 FROM integrante WHERE id = ANY({{INT}}))
     OR EXISTS (SELECT 1 FROM usuario_gestor WHERE id = ANY({{USU}}) AND username <> 'admin') THEN
    RAISE EXCEPTION 'Quedaron registros de prueba que debían borrarse';
  END IF;
END
$v$;
'@
    $sql = $plantilla.Replace('{{DES}}', (Lista-Sql $desarrollos)).Replace('{{INT}}', (Lista-Sql $integrantes)).Replace('{{USU}}', (Lista-Sql $usuarios)).Replace('{{HASH}}', $HashAdmin)
    [IO.File]::WriteAllText($sqlLimpieza, $sql, $utf8SinBom)
    Ejecutar-Archivo $BaseTemporal $sqlLimpieza
    Write-Host 'Datos de prueba borrados, acentos corregidos y admin verificado.'

    Paso 'Alineando las referencias a archivos con los nombres reales (mayúsculas y acentos)'
    $actualizaciones = New-Object System.Text.StringBuilder
    $cambios = 0
    $referencias = 0
    foreach ($par in $ColumnasArchivo) {
        $tabla = $par[0]
        $columna = $par[1]
        foreach ($fila in (Consultar $BaseTemporal "SELECT id, $columna FROM $tabla WHERE coalesce($columna, '') <> ''")) {
            $campos = "$fila" -split "`t", 2
            $id = $campos[0]
            $valor = $campos[1]
            $referencias++
            $normalizado = $valor.Replace('\', '/')
            $coincidencia = [regex]::Match($normalizado, '^(.*?static/)(.+)$')
            if (-not $coincidencia.Success) {
                $Advertencias.Add("$tabla.$columna id=${id}: la referencia no apunta a static/: $valor")
                continue
            }
            $exacto = Resolver-Exacto $coincidencia.Groups[2].Value
            if (-not $exacto) {
                $Advertencias.Add("$tabla.$columna id=${id}: no existe el archivo '$($coincidencia.Groups[2].Value)'")
                continue
            }
            $nuevo = $coincidencia.Groups[1].Value + $exacto
            if ($nuevo -cne $valor) {
                [void]$actualizaciones.AppendLine("UPDATE $tabla SET $columna = '" + $nuevo.Replace("'", "''") + "' WHERE id = $id;")
                $cambios++
            }
        }
    }
    if ($cambios -gt 0) {
        [IO.File]::WriteAllText($sqlReferencias, $actualizaciones.ToString(), $utf8SinBom)
        Ejecutar-Archivo $BaseTemporal $sqlReferencias
    }
    Write-Host "Referencias revisadas: $referencias. Ajustadas: $cambios. Sin archivo: $(($Advertencias | Measure-Object).Count)."

    Paso 'Cantidades de la exportación'
    $cantidadesFinal = Contar-Tablas $BaseTemporal
    foreach ($tabla in $Tablas) {
        Write-Host ('{0,-24} original {1,4}   exportación {2,4}' -f $tabla, $cantidadesOrigen[$tabla], $cantidadesFinal[$tabla])
    }

    Paso 'Exportando la base limpia'
    & (Join-Path $Bin 'pg_dump.exe') @Conexion -d $BaseTemporal --no-owner --no-privileges --encoding=UTF8 -f $volcadoFinal
    if ($LASTEXITCODE -ne 0) { Fallar 'pg_dump falló sobre la base temporal.' }
    $lineas = [IO.File]::ReadAllLines($volcadoFinal, $utf8SinBom)
    $salida = New-Object System.Collections.Generic.List[string]
    $enCopy = $false
    $anteriorEnBlanco = $true
    foreach ($linea in $lineas) {
        if ($enCopy) {
            $salida.Add($linea)
            if ($linea -eq '\.') { $enCopy = $false; $anteriorEnBlanco = $false }
            continue
        }
        if ($linea -match '^COPY .+ FROM stdin;$') {
            $enCopy = $true
            $salida.Add($linea)
            continue
        }
        if ($linea.StartsWith('--')) { continue }
        if ($linea.StartsWith('\restrict') -or $linea.StartsWith('\unrestrict')) { continue }
        if ($linea -eq 'SET transaction_timeout = 0;') { continue }
        if ($linea.Trim() -eq '') {
            if ($anteriorEnBlanco) { continue }
            $anteriorEnBlanco = $true
        } else {
            $anteriorEnBlanco = $false
        }
        $salida.Add($linea)
    }
    $textoFinal = ($salida -join "`n").TrimEnd("`n") + "`n"
}
finally {
    & (Join-Path $Bin 'dropdb.exe') --if-exists @Conexion $BaseTemporal 2>$null
    Remove-Item -LiteralPath $Temporales -Recurse -Force -ErrorAction SilentlyContinue
}

if (Test-Path -LiteralPath $script:DestinoCompleto) {
    $cantidad = (Get-ChildItem -LiteralPath $script:DestinoCompleto -Recurse -Force -File | Measure-Object).Count
    Paso "Reemplazando la carpeta existente ($cantidad archivos)"
    Remove-Item -LiteralPath $script:DestinoCompleto -Recurse -Force
}

Paso "Armando la copia en $($script:DestinoCompleto)"
New-Item -ItemType Directory -Force (Join-Path $script:DestinoCompleto 'db') | Out-Null
[IO.File]::WriteAllText((Join-Path $script:DestinoCompleto 'db\lacis.sql'), $textoFinal, $utf8SinBom)

Copiar-Carpeta 'api' 'api' @('*_test.go')
Copiar-Carpeta 'cmd' 'cmd' @('*_test.go')
Copiar-Carpeta 'internal' 'internal' @('*_test.go')
Copiar-Carpeta 'ui\html' 'ui\html' @()
Copiar-Carpeta 'ui\static\css-GIS' 'ui\static\css-GIS' @()
Copiar-Carpeta 'ui\static\js-GIS' 'ui\static\js-GIS' @()
Copiar-Carpeta 'ui\static\json-GIS' 'ui\static\json-GIS' $excluir
Copiar-Carpeta 'ui\static\assets' 'ui\static\assets' $excluir

foreach ($archivo in @('go.mod', 'go.sum', 'Dockerfile', '.dockerignore')) {
    Copy-Item -LiteralPath (Join-Path $Raiz $archivo) -Destination (Join-Path $script:DestinoCompleto $archivo)
}
Copy-Item -LiteralPath (Join-Path $Raiz 'deploy\DEPLOY.md') -Destination (Join-Path $script:DestinoCompleto 'DEPLOY.md')
[IO.File]::WriteAllText((Join-Path $script:DestinoCompleto '.gitignore'), "*.exe`n", $utf8SinBom)

Paso 'Verificando la copia'
$prohibidos = @('iniciar.bat', 'resetear_db.bat', 'fix_pg.ps1', 'iniciar_db.ps1', 'pg_data', 'tmp_scratch_pw', 'specs', '.specify', '.claude', '.idea', '.git', 'test', 'Modelo Base de datos', 'README.md', 'CLAUDE.md', 'deploy', 'main.exe')
$sobrantes = @(Get-ChildItem -LiteralPath $script:DestinoCompleto -Recurse -Force |
    Where-Object { $prohibidos -contains $_.Name -or $_.Extension -in @('.exe', '.doc') -or $_.Name -like '*.drawio*' -or $_.Name -like '*_test.go' })
if ($sobrantes.Count -gt 0) {
    Fallar ("La copia contiene archivos que no deberían estar:`n" + (($sobrantes | ForEach-Object { $_.FullName }) -join "`n"))
}
Write-Host 'No hay archivos de desarrollo en la copia.'

$extensionesTexto = @('.go', '.html', '.js', '.css', '.json', '.sql', '.md', '.ps1', '.mod', '.sum')
$archivosTexto = @(Get-ChildItem -LiteralPath $script:DestinoCompleto -Recurse -Force -File |
    Where-Object { $extensionesTexto -contains $_.Extension -or $_.Name -in @('Dockerfile', '.dockerignore', '.gitignore') })
$conClave = @($archivosTexto | Select-String -SimpleMatch 'isma_mesa22' -List)
if ($conClave.Count -gt 0) {
    Fallar ("Se encontró la contraseña de la base local en:`n" + (($conClave | ForEach-Object { $_.Path }) -join "`n"))
}
Write-Host 'Ningún archivo contiene la contraseña de la base local.'

$conAdmin = @($archivosTexto | Select-String -SimpleMatch 'admin123' -List | ForEach-Object { Split-Path -Leaf $_.Path })
Write-Host "La contraseña de prueba de admin aparece solo en: $($conAdmin -join ', ')"

Paso 'Comprobando que la base local quedó intacta'
$cantidadesDespues = Contar-Tablas $BaseOrigen
foreach ($tabla in $Tablas) {
    if ($cantidadesDespues[$tabla] -ne $cantidadesOrigen[$tabla]) {
        Fallar "La tabla $tabla de la base local cambió ($($cantidadesOrigen[$tabla]) -> $($cantidadesDespues[$tabla]))."
    }
}
Write-Host 'La base local tiene las mismas cantidades que antes.'

$archivosCopia = Get-ChildItem -LiteralPath $script:DestinoCompleto -Recurse -Force -File
$peso = ($archivosCopia | Measure-Object -Property Length -Sum).Sum
Paso 'Listo'
Write-Host ("Copia creada en {0}: {1} archivos, {2:N1} MB." -f $script:DestinoCompleto, $archivosCopia.Count, ($peso / 1MB))
if ($Advertencias.Count -gt 0) {
    Write-Host ''
    Write-Host 'Advertencias (revisar):' -ForegroundColor Yellow
    foreach ($a in $Advertencias) { Write-Host "  - $a" -ForegroundColor Yellow }
}
exit 0
