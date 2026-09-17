/*
 * Buscador de "Equipo del Proyecto" — un solo clic agrega a la persona
 * (igual que el buscador de Desarrollos), y el rol se elige/edita
 * directamente en la fila de la lista con un <select> inline.
 *
 * Uso: crearWidgetEquipoProyecto({ ...ids..., integrantesActivos, roles, seleccionInicial, onAgregar, onQuitar, onRolCambiado, onSinCoincidencias })
 * Devuelve: { agregarDirecto(id, nombre, rol), resetear(seleccionInicial), obtenerSeleccionados() }
 */
function crearWidgetEquipoProyecto(config) {
    var seleccionados = (config.seleccionInicial || []).slice();
    var integrantesActivos = config.integrantesActivos || [];
    var roles = config.roles || [];

    var buscador = document.getElementById(config.buscadorId);
    var dropdown = document.getElementById(config.dropdownId);
    var listaBox = document.getElementById(config.listaId);

    function iniciales(nombre) {
        return nombre.split(' ').map(function (p) { return p[0] || ''; }).slice(0, 2).join('').toUpperCase();
    }

    function yaSeleccionado(id) {
        return seleccionados.some(function (p) { return String(p.id) === String(id); });
    }

    function opcionesRolHTML(rolActual) {
        var html = '<option value="">Sin rol asignado</option>';
        roles.forEach(function (r) {
            var selected = (r === rolActual) ? ' selected' : '';
            html += '<option value="' + r + '"' + selected + '>' + r + '</option>';
        });
        return html;
    }

    function renderLista() {
        if (seleccionados.length === 0) {
            listaBox.innerHTML = '<li class="list-group-item text-muted small">Aún no hay integrantes cargados.</li>';
            return;
        }
        listaBox.innerHTML = seleccionados.map(function (p, idx) {
            return '<li class="list-group-item d-flex justify-content-between align-items-center gap-2">' +
                '<span class="d-flex align-items-center gap-2 flex-grow-1">' +
                '<span class="chip-ini">' + iniciales(p.nombre) + '</span>' +
                '<span class="fw-bold small">' + p.nombre + '</span>' +
                '</span>' +
                '<select class="form-select form-select-sm w-auto" data-idx="' + idx + '" data-rol-select>' +
                opcionesRolHTML(p.rol) +
                '</select>' +
                '<button type="button" class="btn btn-sm btn-outline-danger" data-idx="' + idx + '" data-quitar>' +
                '<i class="bi bi-x-lg"></i></button>' +
                '</li>';
        }).join('');

        listaBox.querySelectorAll('[data-quitar]').forEach(function (btn) {
            btn.addEventListener('click', function () {
                var idx = parseInt(btn.getAttribute('data-idx'), 10);
                var quitado = seleccionados.splice(idx, 1)[0];
                renderLista();
                if (quitado && typeof config.onQuitar === 'function') config.onQuitar(quitado.id);
            });
        });

        listaBox.querySelectorAll('[data-rol-select]').forEach(function (sel) {
            sel.addEventListener('change', function () {
                var idx = parseInt(sel.getAttribute('data-idx'), 10);
                seleccionados[idx].rol = sel.value;
                if (typeof config.onRolCambiado === 'function') {
                    config.onRolCambiado(seleccionados[idx].id, sel.value);
                }
            });
        });
    }

    function cerrarDropdown() { dropdown.classList.add('d-none'); dropdown.innerHTML = ''; }

    function agregar(id, nombre, rol) {
        if (yaSeleccionado(id)) return;
        var persona = { id: id, nombre: nombre, rol: rol || '' };
        seleccionados.push(persona);
        renderLista();
        if (typeof config.onAgregar === 'function') config.onAgregar(persona);
    }

    function buscar() {
        var q = buscador.value.trim();
        if (!q) { cerrarDropdown(); return; }
        var ql = q.toLowerCase();
        var coincidencias = integrantesActivos.filter(function (p) {
            return p.nombre.toLowerCase().indexOf(ql) !== -1 && !yaSeleccionado(p.id);
        });

        if (coincidencias.length) {
            dropdown.innerHTML = coincidencias.slice(0, 6).map(function (p) {
                return '<div class="participant-search-result" data-id="' + p.id + '" data-nombre="' + p.nombre + '">' +
                    '<span class="chip-ini">' + iniciales(p.nombre) + '</span>' + p.nombre +
                    '</div>';
            }).join('');
        } else {
            dropdown.innerHTML =
                '<div class="participant-search-empty">' +
                'Ningún integrante activo coincide con &quot;' + q + '&quot;.' +
                '<div class="participant-search-add-externo" id="btnCrearNuevoDesdeBusqueda">' +
                '<i class="bi bi-person-plus-fill"></i> Cargar un integrante nuevo' +
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
            agregar(resultado.getAttribute('data-id'), resultado.getAttribute('data-nombre'), '');
            buscador.value = '';
            cerrarDropdown();
            return;
        }
        var crearNuevo = e.target.closest('.participant-search-add-externo');
        if (crearNuevo) {
            var query = buscador.value.trim();
            cerrarDropdown();
            if (typeof config.onSinCoincidencias === 'function') config.onSinCoincidencias(query);
        }
    });

    document.addEventListener('click', function (e) {
        if (!e.target.closest('#' + config.wrapId)) cerrarDropdown();
    });

    renderLista();

    return {
        agregarDirecto: function (id, nombre, rol) {
            agregar(id, nombre, rol);
        },
        resetear: function (seleccionInicial) {
            seleccionados = (seleccionInicial || []).slice();
            renderLista();
        },
        obtenerSeleccionados: function () {
            return seleccionados.slice();
        }
    };
}