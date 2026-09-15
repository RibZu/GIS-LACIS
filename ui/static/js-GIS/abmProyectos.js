let proyectos = [];
let idProyectoAEliminar = null;
let todosLosIntegrantes = [];
let proyectoIdActual = null;
let widgetEquipo = null;

const ROLES_PROYECTO = ['Director / Co-Director', 'Investigador', 'Asesor Externo', 'Estudiante / Becario'];

const cardsGrid = document.getElementById('cardsGrid');
const emptyState = document.getElementById('emptyState');
const inputBuscar = document.getElementById('inputBuscar');
const selectOrden = document.getElementById('selectOrden');

const drawerEl = document.getElementById('drawerProyecto');
const drawerBS = new bootstrap.Offcanvas(drawerEl);
const modalEliminarBS = new bootstrap.Modal(document.getElementById('modalEliminar'));
const modalMiniIntegranteBS = new bootstrap.Modal(document.getElementById('modalMiniIntegrante'));
const toastBS = new bootstrap.Toast(document.getElementById('liveToast'));

document.addEventListener('DOMContentLoaded', () => {
    cargarDesdeAPI();
    cargarIntegrantes();
    configurarEventos();
});

// ==================== PROYECTOS ====================

async function cargarDesdeAPI() {
    try {
        const res = await fetch('/api/v1/admin/proyectos-todos');
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
        cargarDesdeAPI();
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

    let resultado = proyectos.filter(p => {
        const matchesText = (p.titulo || '').toLowerCase().includes(texto) ||
            (p.descripcion || '').toLowerCase().includes(texto) ||
            (p.equipo_historico || '').toLowerCase().includes(texto);
        return matchesText;
    });

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
        const esVigente = (p.anio_fin || 0) >= anioActual;
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

                    <div class="pt-3 mt-auto d-flex justify-content-between align-items-center">
                        <small class="text-muted text-truncate me-2" style="max-width: 180px;">
                            
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

// ==================== DRAWER CREAR / EDITAR ====================

function abrirDrawerCrear() {
    document.getElementById('formProyecto').reset();
    document.getElementById('formId').value = '';
    proyectoIdActual = null;
    if (widgetEquipo) widgetEquipo.resetear([]);
    document.getElementById('drawerTitulo').innerHTML = '<i class="bi bi-plus-circle text-primary"></i> <span>Nuevo Proyecto</span>';
    document.getElementById('formAnioInicio').value = new Date().getFullYear();
    document.getElementById('formAnioFin').value = new Date().getFullYear() + 2;
    drawerBS.show();
}

function abrirDrawerEditar(id) {
    const p = proyectos.find(item => item.id === id);
    if (!p) return;

    proyectoIdActual = p.id;
    document.getElementById('formId').value = p.id;
    document.getElementById('formTitulo').value = p.titulo || '';
    document.getElementById('formAnioInicio').value = p.anio_inicio || 2024;
    document.getElementById('formAnioFin').value = p.anio_fin || 2026;
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
        descripcion: document.getElementById('formDescripcion').value.trim(),
        enlace: document.getElementById('formEnlace').value.trim()
    };

    try {
        if (!id) {
            // Crear proyecto nuevo
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

            // Asignar el equipo armado en memoria mientras se completaba el form
            const equipoAAsignar = widgetEquipo ? widgetEquipo.obtenerSeleccionados() : [];
            for (const m of equipoAAsignar) {
                await fetch(`/api/v1/admin/proyectos/${nuevoProyecto.id}/equipo`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ integrante_id: parseInt(m.id), rol_en_proyecto: m.rol })
                });
            }

            drawerBS.hide();
            mostrarToast('Proyecto creado con su equipo con éxito.');
            cargarDesdeAPI();
        } else {
            // Actualizar proyecto existente (el equipo se maneja aparte, en vivo)
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
        }
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
});

// ==================== ELIMINAR PROYECTO ====================

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

// ==================== EQUIPO DEL PROYECTO (buscador estilo Desarrollos) ====================

async function cargarIntegrantes() {
    try {
        const res = await fetch('/api/v1/integrantes');
        const data = await res.json() || [];
        const activos = data
            .filter(i => i.activo !== false)
            .map(i => ({ id: i.id, nombre: `${i.nombre} ${i.apellido}` }));
        // Mutamos el mismo array (no reasignamos) para que el widget vea los cambios.
        todosLosIntegrantes.length = 0;
        activos.forEach(i => todosLosIntegrantes.push(i));
    } catch (err) {
        console.error('No se pudo cargar la lista de integrantes', err);
    }
    if (!widgetEquipo) inicializarWidgetEquipo();
}

function inicializarWidgetEquipo() {
    widgetEquipo = crearWidgetEquipoProyecto({
        wrapId: 'equipoProyectoWrap',
        buscadorId: 'inputBuscarIntegrante',
        dropdownId: 'resultadosBusqueda',
        listaId: 'listaEquipoActual',
        integrantesActivos: todosLosIntegrantes,
        roles: ROLES_PROYECTO,
        seleccionInicial: [],
        onAgregar: (persona) => {
            if (proyectoIdActual) agregarMiembroAlServidor(persona.id, persona.rol);
        },
        onRolCambiado: (integranteId, nuevoRol) => {
            if (proyectoIdActual) agregarMiembroAlServidor(integranteId, nuevoRol);
        },
        onQuitar: (integranteId) => {
            if (proyectoIdActual) quitarMiembroDelServidor(integranteId);
        },
        onSinCoincidencias: (query) => {
            const partes = query.split(' ');
            document.getElementById('miniNombre').value = partes[0] || '';
            document.getElementById('miniApellido').value = partes.slice(1).join(' ') || '';
            modalMiniIntegranteBS.show();
        }
    });
}

async function agregarMiembroAlServidor(integranteId, rol) {
    try {
        const res = await fetch(`/api/v1/admin/proyectos/${proyectoIdActual}/equipo`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ integrante_id: parseInt(integranteId), rol_en_proyecto: rol })
        });
        if (!res.ok) throw new Error('No se pudo agregar el integrante');
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
}

async function quitarMiembroDelServidor(integranteId) {
    try {
        const res = await fetch(`/api/v1/admin/proyectos/${proyectoIdActual}/equipo/${integranteId}`, { method: 'DELETE' });
        if (!res.ok) throw new Error('No se pudo quitar el integrante');
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
}

async function cargarEquipoDeProyecto(id) {
    let equipo = [];
    try {
        const res = await fetch(`/api/v1/proyectos/${id}/equipo`);
        equipo = await res.json() || [];
    } catch (err) {
        equipo = [];
    }
    const seleccionInicial = equipo.map(m => ({
        id: m.integrante_id,
        nombre: `${m.nombre} ${m.apellido}`,
        rol: m.rol_en_proyecto || ''
    }));
    if (widgetEquipo) widgetEquipo.resetear(seleccionInicial);
}

// ==================== CARGA RÁPIDA DE INTEGRANTE ====================

document.getElementById('linkCrearIntegranteRapido').addEventListener('click', (e) => {
    e.preventDefault();
    document.getElementById('miniNombre').value = '';
    document.getElementById('miniApellido').value = '';
    document.getElementById('miniPerteneceLacis').checked = false;
    modalMiniIntegranteBS.show();
});

document.getElementById('btnGuardarMiniIntegrante').addEventListener('click', async () => {
    const rolSelect = document.getElementById('miniRolId');
    const payload = {
        nombre: document.getElementById('miniNombre').value.trim(),
        apellido: document.getElementById('miniApellido').value.trim(),
        rol_id: parseInt(rolSelect.value),
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
        const nombreCompleto = `${nuevo.nombre} ${nuevo.apellido}`;
        todosLosIntegrantes.push({ id: nuevo.id, nombre: nombreCompleto });

        const rolPorDefecto = rolSelect.options[rolSelect.selectedIndex].text;
        if (widgetEquipo) widgetEquipo.agregarDirecto(nuevo.id, nombreCompleto, rolPorDefecto);

        modalMiniIntegranteBS.hide();
        mostrarToast('Integrante creado y agregado al equipo.');
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
});

// ==================== FILTROS ====================

function configurarEventos() {
    inputBuscar.addEventListener('input', () => renderizarCards());
    selectOrden.addEventListener('change', () => renderizarCards());
}