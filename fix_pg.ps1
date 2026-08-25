$confPath = "C:\Program Files\PostgreSQL\18\data\pg_hba.conf"
$lines = Get-Content $confPath
$newLines = @()
foreach ($line in $lines) {
    if ($line -match "^host.*127\.0\.0\.1/32.*" -or $line -match "^host.*::1/128.*" -or $line -match "^local.*all.*all.*") {
        $line = $line -replace "scram-sha-256", "trust"
        $line = $line -replace "md5", "trust"
    }
    $newLines += $line
}
$newLines | Set-Content $confPath
Restart-Service -Name "postgresql-x64-18" -Force
