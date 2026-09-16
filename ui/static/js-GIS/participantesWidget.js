function crearWidgetParticipantes(config) {
    var seleccionados = (config.seleccionInicial || []).slice();
    var integrantesActivos = config.integrantesActivos || [];

    var buscador = document.getElementById(config.buscadorId);
    var dropdown = document.getElementById(config.dropdownId);
    var chipsBox = document.getElementById(config.chipsId);
    var hiddenBox = document.getElementById(config.hiddenId);

    function iniciales(nombre) {
        return nombre.split(' ').map(function (p) { return p[0]; }).slice(0, 2).join('').toUpperCase();
    }

    function yaSeleccionado(nombre) {
        return seleccionados.some(function (p) { return p.nombre.toLowerCase() === nombre.toLowerCase(); });
    }

    function chipIntegranteHTML(p) {
        return '<span class="chip-integrante"><span class="ini">' + iniciales(p.nombre) + '</span>' + p.nombre + '</span>';
    }

    function renderChips() {
        chipsBox.innerHTML = seleccionados.map(function (p, idx) {
            if (p.tipo === 'integrante') {
                return '<span class="chip-integrante"><span class="ini">' + iniciales(p.nombre) + '</span>' + p.nombre +
                    '<button type="button" class="chip-remove" data-idx="' + idx + '" aria-label="Quitar ' + p.nombre + '">&times;</button></span>';
            }
            return '<span class="chip-externo"><span class="ini">' + iniciales(p.nombre) + '</span>' + p.nombre +
                '<span class="tag-externo">externo</span>' +
                '<button type="button" class="chip-remove" data-idx="' + idx + '" aria-label="Quitar ' + p.nombre + '">&times;</button></span>';
        }).join('');

        hiddenBox.innerHTML = seleccionados.map(function (p) {
            if (p.tipo === 'integrante') {
                return '<input type="hidden" name="integrantes" value="' + p.id + '">';
            }
            return '<input type="hidden" name="externos" value="' + p.nombre.replace(/"/g, '&quot;') + '">';
        }).join('');

        chipsBox.querySelectorAll('.chip-remove').forEach(function (btn) {
            btn.addEventListener('click', function () {
                seleccionados.splice(parseInt(btn.getAttribute('data-idx'), 10), 1);
                renderChips();
            });
        });
    }

    function cerrarDropdown() { dropdown.classList.add('d-none'); dropdown.innerHTML = ''; }

    function buscar() {
        var q = buscador.value.trim();
        if (!q) { cerrarDropdown(); return; }
        var ql = q.toLowerCase();
        var coincidencias = integrantesActivos.filter(function (p) {
            return p.nombre.toLowerCase().indexOf(ql) !== -1 && !yaSeleccionado(p.nombre);
        });

        if (coincidencias.length) {
            dropdown.innerHTML = coincidencias.slice(0, 5).map(function (p) {
                return '<div class="participant-search-result" data-id="' + p.id + '" data-nombre="' + p.nombre + '">' + chipIntegranteHTML(p) + '</div>';
            }).join('');
        } else {
            dropdown.innerHTML =
                '<div class="participant-search-empty">' +
                'Ningún integrante activo coincide con &quot;' + q + '&quot;.' +
                '<div class="participant-search-add-externo" data-externo="' + q + '">' +
                '<i class="bi bi-person-plus-fill"></i> Agregar &quot;' + q + '&quot; como participante externo' +
                '</div>' +
                '</div>';
        }
        dropdown.classList.remove('d-none');
    }

    buscador.addEventListener('input', buscar);
    buscador.addEventListener('focus', buscar);

    dropdown.addEventListener('click', function (e) {
        var resultado = e.target.closest('.participant-search-result');
        if (resultado) {
            seleccionados.push({ id: resultado.getAttribute('data-id'), nombre: resultado.getAttribute('data-nombre'), tipo: 'integrante' });
            buscador.value = ''; cerrarDropdown(); renderChips();
            return;
        }
        var externo = e.target.closest('.participant-search-add-externo');
        if (externo) {
            var nombre = externo.getAttribute('data-externo');
            if (nombre && !yaSeleccionado(nombre)) seleccionados.push({ nombre: nombre, tipo: 'externo' });
            buscador.value = ''; cerrarDropdown(); renderChips();
        }
    });

    document.addEventListener('click', function (e) {
        if (!e.target.closest('#' + config.wrapId)) cerrarDropdown();
    });

    var btnManual = config.btnManualId ? document.getElementById(config.btnManualId) : null;
    var wrapManual = config.manualWrapId ? document.getElementById(config.manualWrapId) : null;
    var inputManual = config.manualInputId ? document.getElementById(config.manualInputId) : null;
    var btnConfirmarManual = config.manualConfirmId ? document.getElementById(config.manualConfirmId) : null;

    if (btnManual && wrapManual && inputManual && btnConfirmarManual) {
        btnManual.addEventListener('click', function () {
            wrapManual.classList.remove('d-none');
            inputManual.focus();
        });

        function confirmarManual() {
            var nombre = inputManual.value.trim();
            if (nombre && !yaSeleccionado(nombre)) {
                seleccionados.push({ nombre: nombre, tipo: 'externo' });
                renderChips();
            }
            inputManual.value = '';
            wrapManual.classList.add('d-none');
        }

        btnConfirmarManual.addEventListener('click', confirmarManual);
        inputManual.addEventListener('keydown', function (e) {
            if (e.key === 'Enter') { e.preventDefault(); confirmarManual(); }
        });
    }

    renderChips();
}
