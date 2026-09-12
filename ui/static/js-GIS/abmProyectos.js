let proyectos = [];
        let idProyectoAEliminar = null;
        let todosLosIntegrantes = [];
        let equipoActual = [];       
        let equipoPendiente = [];   
        let proyectoIdActual = null;
        let integranteSeleccionado = null;
        const modalMiniIntegranteBS = new bootstrap.Modal(document.getElementById('modalMiniIntegrante'));

        const cardsGrid = document.getElementById('cardsGrid');
        const emptyState = document.getElementById('emptyState');
        const inputBuscar = document.getElementById('inputBuscar');
        const selectOrden = document.getElementById('selectOrden');

        const drawerEl = document.getElementById('drawerProyecto');
        const drawerBS = new bootstrap.Offcanvas(drawerEl);
        const modalEliminarBS = new bootstrap.Modal(document.getElementById('modalEliminar'));
        const toastBS = new bootstrap.Toast(document.getElementById('liveToast'));

        document.addEventListener('DOMContentLoaded', () => {
            cargarDesdeAPI();
            cargarIntegrantes();
            configurarEventos();
        });

        async function cargarDesdeAPI() {
            try {
                const res = await fetch('/api/v1/admin/proyectos-todos');  // antes: '/api/v1/admin/proyectos'
                if (!res.ok) throw new Error('Error de conexión');
                proyectos = await res.json() || [];
                renderizarCards();
            } catch (err) {
                mostrarToast('Error al conectar con la API', 'danger');
            }
        }
        async function restaurarProyecto(id) {
            try {
                const res = await fetch(`/api/v1/admin/proyectos/${id}/restaurar`, {
                    method: 'PATCH',
                    headers: { 'Content-Type': 'application/json' }
                });

                if (!res.ok) throw new Error('No se pudo restaurar el proyecto');

                mostrarToast('Proyecto restaurado con éxito.', 'success');
                cargarDesdeAPI(); // Recarga la lista para que vuelva a su estado normal
            } catch (err) {
                mostrarToast(err.message, 'danger');
            }
        }

        function mostrarToast(mensaje, tipo = 'success') {
            const toastEl = document.getElementById('liveToast');
            toastEl.className = `toast align-items-center text-bg-${tipo} border-0 shadow`;
            document.getElementById('toastMessage').textContent = mensaje;
            toastBS.show();
        }

        function obtenerProyectosFiltrados() {
            const texto = inputBuscar.value.toLowerCase();
            const orden = selectOrden.value;

            // Filtro solo por búsqueda de texto
            let resultado = proyectos.filter(p => {
                const matchesText = (p.titulo || '').toLowerCase().includes(texto) ||
                    (p.descripcion || '').toLowerCase().includes(texto) ||
                    (p.equipo_historico || '').toLowerCase().includes(texto);
                return matchesText;
            });

            // Ordenamiento
            resultado.sort((a, b) => {
                if (orden === 'recientes') return (b.anio_fin || 0) - (a.anio_fin || 0);
                if (orden === 'antiguos') return (a.anio_inicio || 0) - (b.anio_inicio || 0);
                if (orden === 'titulo') return (a.titulo || '').localeCompare(b.titulo || '');
            });

            return resultado;
        }

        function renderizarCards() {
            const listado = obtenerProyectosFiltrados();
            cardsGrid.innerHTML = '';

            if (listado.length === 0) {
                emptyState.classList.remove('d-none');
                return;
            }

            emptyState.classList.add('d-none');
            const anioActual = new Date().getFullYear();

            listado.forEach(p => {
                const anioActual = new Date().getFullYear();
                const esVigente = (p.anio_fin || 0) >= anioActual;
                
                // Validamos si el proyecto está activo (por defecto true si la propiedad no viene)
                const estaActivo = p.activo !== false; 
                const botonesAccion = estaActivo ? `
                    <button class="btn btn-light border btn-sm" onclick="abrirDrawerEditar(${p.id})" title="Editar"><i class="bi bi-pencil"></i></button>
                    <button class="btn btn-light border text-danger btn-sm" onclick="confirmarEliminar(${p.id})" title="Eliminar"><i class="bi bi-trash"></i></button>
                ` : `
                    <button class="btn btn-light border btn-sm" onclick="abrirDrawerEditar(${p.id})" title="Editar"><i class="bi bi-pencil"></i></button>
                    <button class="btn btn-success text-white btn-sm" onclick="restaurarProyecto(${p.id})" title="Restaurar proyecto"><i class="bi bi-arrow-counterclockwise"></i> Restaurar</button>
                `;
                const col = document.createElement('div');
                col.className = 'col-12 col-md-6 col-lg-4';

                col.innerHTML = `
                    <div class="proyecto-card ${!estaActivo ? 'bg-light border-opacity-50 opacity-75' : (esVigente ? '' : 'historico')}">
                        <div class="card-top-bar ${!estaActivo ? 'bg-secondary' : ''}"></div>
                        <div class="p-4 d-flex flex-column flex-grow-1">
                            <div class="d-flex justify-content-between align-items-center mb-2">
                                <span class="${!estaActivo ? 'badge bg-secondary text-white' : (esVigente ? 'badge-vigente' : 'badge-historico')}">
                                    <i class="bi ${!estaActivo ? 'bi-eye-slash-fill' : (esVigente ? 'bi-record-fill' : 'bi-archive-fill')}"></i>
                                    ${!estaActivo ? 'Inactivo' : `${p.anio_inicio} - ${p.anio_fin}`}
                                </span>
                                ${p.enlace ? `<a href="${p.enlace}" target="_blank" class="text-primary small text-decoration-none"><i class="bi bi-box-arrow-up-right"></i> Web</a>` : ''}
                            </div>
                            
                            <h5 class="fw-bold ${!estaActivo ? 'text-muted text-decoration-line-through' : 'text-dark'} mb-2" style="font-size: 1.05rem; line-height: 1.35;">${p.titulo}</h5>
                            
                            <p class="text-secondary small mb-3 flex-grow-1" style="display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden;">
                                ${p.descripcion || 'Sin descripción detallada.'}
                            </p>

                            <div class="pt-3 border-top mt-auto d-flex justify-content-between align-items-center">
                                <small class="text-muted text-truncate me-2" style="max-width: 180px;" title="${p.equipo_historico || ''}">
                                    <i class="bi bi-people-fill text-primary me-1"></i>${p.equipo_historico || 'Sin equipo'}
                                </small>
                                <div class="btn-group btn-group-sm">
                                    ${botonesAccion}
                                </div>
                            </div>
                        </div>
                    </div>
                `;
                cardsGrid.appendChild(col);
            });
        }

        function abrirDrawerCrear() {
            document.getElementById('formProyecto').reset();
            document.getElementById('formId').value = '';
            proyectoIdActual = null;
            equipoActual = [];
            equipoPendiente = [];
            integranteSeleccionado = null;
            document.getElementById('seleccionActual').classList.add('d-none');
            document.getElementById('btnAgregarMiembro').disabled = true;
            renderizarEquipoActual();
            document.getElementById('drawerTitulo').innerHTML = '<i class="bi bi-plus-circle text-primary"></i> <span>Nuevo Proyecto</span>';
            document.getElementById('formAnioInicio').value = new Date().getFullYear();
            document.getElementById('formAnioFin').value = new Date().getFullYear() + 2;
            drawerBS.show();
        }

        function abrirDrawerEditar(id) {
            const p = proyectos.find(item => item.id === id);
            if (!p) return;

            proyectoIdActual = p.id;
            equipoPendiente = [];
            document.getElementById('formId').value = p.id;
            document.getElementById('formTitulo').value = p.titulo || '';
            document.getElementById('formAnioInicio').value = p.anio_inicio || 2024;
            document.getElementById('formAnioFin').value = p.anio_fin || 2026;
            document.getElementById('formEquipo').value = p.equipo_historico || '';
            document.getElementById('formDescripcion').value = p.descripcion || '';
            document.getElementById('formEnlace').value = p.enlace || '';

            cargarEquipoDeProyecto(p.id);

            document.getElementById('drawerTitulo').innerHTML = '<i class="bi bi-pencil-square text-primary"></i> <span>Editar Proyecto #' + p.id + '</span>';
            drawerBS.show();
        }

        document.getElementById('formProyecto').addEventListener('submit', async (e) => {
            e.preventDefault();
            const id = document.getElementById('formId').value;
            const payload = {
                titulo: document.getElementById('formTitulo').value.trim(),
                anio_inicio: parseInt(document.getElementById('formAnioInicio').value),
                anio_fin: parseInt(document.getElementById('formAnioFin').value),
                equipo_historico: document.getElementById('formEquipo').value.trim(),
                descripcion: document.getElementById('formDescripcion').value.trim(),
                enlace: document.getElementById('formEnlace').value.trim()
            };

            try {
                if (!id) {
                    const res = await fetch('/api/v1/admin/proyectos', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload)
                    });
                    if (!res.ok) {
                        const err = await res.json();
                        throw new Error(err.error || 'Error al guardar');
                    }
                    const nuevoProyecto = await res.json();
                    proyectoIdActual = nuevoProyecto.id;
                    document.getElementById('formId').value = nuevoProyecto.id;
                    document.getElementById('drawerTitulo').innerHTML = '<i class="bi bi-pencil-square text-primary"></i> <span>Editar Proyecto #' + nuevoProyecto.id + '</span>';
                    cargarDesdeAPI();
                    mostrarToast('Proyecto creado. Ahora podés agregar el equipo abajo.');
                    return; // Dejamos el drawer abierto para cargar el equipo
                }

                const res = await fetch(`/api/v1/admin/proyectos/${id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });
                if (!res.ok) {
                    const err = await res.json();
                    throw new Error(err.error || 'Error al guardar');
                }
                drawerBS.hide();
                mostrarToast('Proyecto actualizado con éxito.');
                cargarDesdeAPI();
            } catch (err) {
                mostrarToast(err.message, 'danger');
            }
        });

        function confirmarEliminar(id) {
            idProyectoAEliminar = id;
            modalEliminarBS.show();
        }

        document.getElementById('btnConfirmarEliminar').addEventListener('click', async () => {
            if (idProyectoAEliminar !== null) {
                try {
                    const res = await fetch(`/api/v1/admin/proyectos/${idProyectoAEliminar}`, { method: 'DELETE' });
                    if (!res.ok) throw new Error('No se pudo eliminar');
                    modalEliminarBS.hide();
                    mostrarToast('Proyecto eliminado.', 'danger');
                    cargarDesdeAPI();
                } catch (err) {
                    mostrarToast(err.message, 'danger');
                } finally {
                    idProyectoAEliminar = null;
                }
            }
        });

        function configurarEventos() {
            inputBuscar.addEventListener('input', () => renderizarCards());
            selectOrden.addEventListener('change', () => renderizarCards());
        }
        async function cargarIntegrantes() {
            try {
                const res = await fetch('/api/v1/integrantes');
                todosLosIntegrantes = await res.json() || [];
                document.getElementById('listaIntegrantes').innerHTML = todosLosIntegrantes
                    .map(i => `<option data-id="${i.id}" value="${i.nombre} ${i.apellido}">`)
                    .join('');
            } catch (err) {
                console.error('No se pudo cargar la lista de integrantes', err);
            }
        }

        const inputBuscarIntegrante = document.getElementById('inputBuscarIntegrante');
        const resultadosBusqueda = document.getElementById('resultadosBusqueda');

        inputBuscarIntegrante.addEventListener('input', () => {
            const texto = inputBuscarIntegrante.value.trim().toLowerCase();
            if (!texto) {
                resultadosBusqueda.style.display = 'none';
                return;
            }

            const yaEnEquipo = new Set(obtenerEquipoVisible().map(m => m.integrante_id));
            const coincidencias = todosLosIntegrantes
                .filter(i => !yaEnEquipo.has(i.id))
                .filter(i => `${i.nombre} ${i.apellido}`.toLowerCase().includes(texto))
                .slice(0, 8);

            resultadosBusqueda.innerHTML = '';
            if (coincidencias.length === 0) {
                resultadosBusqueda.innerHTML = '<div class="list-group-item small text-muted">Sin coincidencias.</div>';
            } else {
                coincidencias.forEach(i => {
                    const item = document.createElement('button');
                    item.type = 'button';
                    item.className = 'list-group-item list-group-item-action small';
                    item.textContent = `${i.nombre} ${i.apellido}`;
                    item.addEventListener('click', () => seleccionarIntegrante(i));
                    resultadosBusqueda.appendChild(item);
                });
            }
            resultadosBusqueda.style.display = 'block';
        });

        // Cierra el dropdown al hacer clic fuera
        document.addEventListener('click', (e) => {
            if (!resultadosBusqueda.contains(e.target) && e.target !== inputBuscarIntegrante) {
                resultadosBusqueda.style.display = 'none';
            }
        });

        function seleccionarIntegrante(i) {
            integranteSeleccionado = i;
            document.getElementById('seleccionNombre').textContent = `${i.nombre} ${i.apellido}`;
            document.getElementById('seleccionActual').classList.remove('d-none');
            document.getElementById('btnAgregarMiembro').disabled = false;
            inputBuscarIntegrante.value = '';
            resultadosBusqueda.style.display = 'none';
        }

        document.getElementById('btnCancelarSeleccion').addEventListener('click', () => {
            integranteSeleccionado = null;
            document.getElementById('seleccionActual').classList.add('d-none');
            document.getElementById('btnAgregarMiembro').disabled = true;
        });

        async function cargarEquipoDeProyecto(id) {
            try {
                const res = await fetch(`/api/v1/proyectos/${id}/equipo`);
                equipoActual = await res.json() || [];
            } catch (err) {
                equipoActual = [];
            }
            renderizarEquipoActual();
        }

        function obtenerEquipoVisible() {
    return proyectoIdActual ? equipoActual : equipoPendiente;
    }

    function renderizarEquipoActual() {
        const lista = document.getElementById('listaEquipoActual');
        const equipo = obtenerEquipoVisible();
        lista.innerHTML = '';

        if (equipo.length === 0) {
            lista.innerHTML = '<li class="list-group-item text-muted small">Aún no hay integrantes cargados.</li>';
            return;
        }

        equipo.forEach((m, idx) => {
            const li = document.createElement('li');
            li.className = 'list-group-item d-flex justify-content-between align-items-center';
            li.innerHTML = `
                <span><strong>${m.nombre} ${m.apellido}</strong>${m.rol_en_proyecto ? ' — ' + m.rol_en_proyecto : ''}</span>
                <button type="button" class="btn btn-sm btn-outline-danger" data-idx="${idx}">
                    <i class="bi bi-x-lg"></i>
                </button>`;
            li.querySelector('button').addEventListener('click', () => quitarMiembroEquipo(m.integrante_id, idx));
            lista.appendChild(li);
        });
    }

        async function agregarMiembroEquipo(integranteId, rol) {
            if (!proyectoIdActual) {
                mostrarToast('Primero guardá el proyecto para poder agregar el equipo.', 'danger');
                return;
            }
            try {
                const res = await fetch(`/api/v1/admin/proyectos/${proyectoIdActual}/equipo`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ integrante_id: integranteId, rol_en_proyecto: rol })
                });
                if (!res.ok) throw new Error('No se pudo agregar el integrante');
                await cargarEquipoDeProyecto(proyectoIdActual);
            } catch (err) {
                mostrarToast(err.message, 'danger');
            }
        }

        async function quitarMiembroEquipo(integranteId, idx) {
            if (!proyectoIdActual) {
                // Todavía no se guardó nada en el servidor: solo sacarlo de la lista local
                equipoPendiente.splice(idx, 1);
                renderizarEquipoActual();
                return;
            }
            try {
                const res = await fetch(`/api/v1/admin/proyectos/${proyectoIdActual}/equipo/${integranteId}`, { method: 'DELETE' });
                if (!res.ok) throw new Error('No se pudo quitar el integrante');
                await cargarEquipoDeProyecto(proyectoIdActual);
            } catch (err) {
                mostrarToast(err.message, 'danger');
            }
        }

        document.getElementById('btnAgregarMiembro').addEventListener('click', async () => {
            if (!integranteSeleccionado) return;
            const rol = document.getElementById('inputRolEnProyecto').value.trim();

            if (proyectoIdActual) {
                // El proyecto ya existe: se persiste al toque
                await agregarMiembroEquipo(integranteSeleccionado.id, rol);
            } else {
                // Proyecto nuevo: se guarda en memoria hasta hacer clic en "Guardar Proyecto"
                equipoPendiente.push({
                    integrante_id: integranteSeleccionado.id,
                    nombre: integranteSeleccionado.nombre,
                    apellido: integranteSeleccionado.apellido,
                    rol_en_proyecto: rol
                });
                renderizarEquipoActual();
            }

            integranteSeleccionado = null;
            document.getElementById('seleccionActual').classList.add('d-none');
            document.getElementById('btnAgregarMiembro').disabled = true;
            document.getElementById('inputRolEnProyecto').value = '';
        });

        document.getElementById('linkCrearIntegranteRapido').addEventListener('click', (e) => {
            e.preventDefault();
            document.getElementById('miniNombre').value = '';
            document.getElementById('miniApellido').value = '';
            document.getElementById('miniPerteneceLacis').checked = false;
            modalMiniIntegranteBS.show();
        });

        document.getElementById('btnGuardarMiniIntegrante').addEventListener('click', async () => {
            const payload = {
                nombre: document.getElementById('miniNombre').value.trim(),
                apellido: document.getElementById('miniApellido').value.trim(),
                rol_id: parseInt(document.getElementById('miniRolId').value),
                pertenece_lacis: document.getElementById('miniPerteneceLacis').checked
            };
            if (!payload.nombre || !payload.apellido) {
                mostrarToast('Nombre y apellido son obligatorios.', 'danger');
                return;
            }
            try {
                const res = await fetch('/api/v1/admin/integrantes/mini', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });
                if (!res.ok) {
                    const err = await res.json();
                    throw new Error(err.error || 'No se pudo crear el integrante');
                }
                const nuevo = await res.json();
                await cargarIntegrantes();
                modalMiniIntegranteBS.hide();
                await agregarMiembroEquipo(nuevo.id, document.getElementById('inputRolEnProyecto').value.trim());
                mostrarToast('Integrante creado y agregado al equipo.');
            } catch (err) {
                mostrarToast(err.message, 'danger');
            }
        });