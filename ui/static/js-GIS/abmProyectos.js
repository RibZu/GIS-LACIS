let proyectos = [];
        let idProyectoAEliminar = null;

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
            document.getElementById('drawerTitulo').innerHTML = '<i class="bi bi-plus-circle text-primary"></i> <span>Nuevo Proyecto</span>';
            document.getElementById('formAnioInicio').value = new Date().getFullYear();
            document.getElementById('formAnioFin').value = new Date().getFullYear() + 2;
            drawerBS.show();
        }

        function abrirDrawerEditar(id) {
            const p = proyectos.find(item => item.id === id);
            if (!p) return;

            document.getElementById('formId').value = p.id;
            document.getElementById('formTitulo').value = p.titulo || '';
            document.getElementById('formAnioInicio').value = p.anio_inicio || 2024;
            document.getElementById('formAnioFin').value = p.anio_fin || 2026;
            document.getElementById('formEquipo').value = p.equipo_historico || '';
            document.getElementById('formDescripcion').value = p.descripcion || '';
            document.getElementById('formEnlace').value = p.enlace || '';

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
                let res;
                if (!id) {
                    res = await fetch('/api/v1/admin/proyectos', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload)
                    });
                } else {
                    res = await fetch(`/api/v1/admin/proyectos/${id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload)
                    });
                }

                if (!res.ok) {
                    const err = await res.json();
                    throw new Error(err.error || 'Error al guardar');
                }

                drawerBS.hide();
                mostrarToast(id ? 'Proyecto actualizado con éxito.' : 'Proyecto creado con éxito.');
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
            // Ya no escuchamos a los radio buttons de vigencia, solo al buscador y al selector de orden
            inputBuscar.addEventListener('input', () => renderizarCards());
            selectOrden.addEventListener('change', () => renderizarCards());
        }