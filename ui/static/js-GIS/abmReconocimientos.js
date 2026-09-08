let reconocimientos = [];
let idReconocimientoABajar = null;

const cardsGrid = document.getElementById('cardsGrid');
const emptyState = document.getElementById('emptyState');
const inputBuscar = document.getElementById('inputBuscar');
const selectOrden = document.getElementById('selectOrden');

const drawerEl = document.getElementById('drawerReconocimiento');
const drawerBS = new bootstrap.Offcanvas(drawerEl);
const modalEliminarBS = new bootstrap.Modal(document.getElementById('modalEliminar'));
const toastBS = new bootstrap.Toast(document.getElementById('liveToast'));

document.addEventListener('DOMContentLoaded', () => {
    cargarDesdeAPI();
    configurarEventos();
});

async function cargarDesdeAPI() {
    try {
        const res = await fetch('/api/v1/admin/reconocimientos-todos');
        if (!res.ok) throw new Error('Error de conexión');
        reconocimientos = await res.json() || [];
        renderizarCards();
    } catch (err) {
        mostrarToast('Error al conectar con la API', 'danger');
    }
}

function mostrarToast(mensaje, tipo = 'success') {
    const toastEl = document.getElementById('liveToast');
    toastEl.className = `toast align-items-center text-bg-${tipo} border-0 shadow`;
    document.getElementById('toastMessage').textContent = mensaje;
    toastBS.show();
}

function obtenerReconocimientosFiltrados() {
    const texto = inputBuscar.value.toLowerCase();
    const orden = selectOrden.value;
    let resultado = reconocimientos.filter(r => {
        return (r.titulo || '').toLowerCase().includes(texto) ||
            (r.descripcion || '').toLowerCase().includes(texto);
    });
    resultado.sort((a, b) => {
        if (orden === 'recientes') return b.id - a.id;
        if (orden === 'antiguos') return a.id - b.id;
        if (orden === 'titulo') return (a.titulo || '').localeCompare(b.titulo || '');
        return 0;
    });
    return resultado;
}

function renderizarCards() {
    const listado = obtenerReconocimientosFiltrados();
    cardsGrid.innerHTML = '';
    if (listado.length === 0) {
        emptyState.classList.remove('d-none');
        return;
    }
    emptyState.classList.add('d-none');

    listado.forEach(r => {
        const estaActivo = r.activo !== false;
        const col = document.createElement('div');
        col.className = 'col-12 col-md-6 col-lg-4';

        const botonesAccion = estaActivo
            ? `<button class="btn btn-light border btn-sm" onclick="abrirDrawerEditar(${r.id})" title="Editar">
                    <i class="bi bi-pencil"></i>
               </button>
               <button class="btn btn-light border text-danger btn-sm" onclick="confirmarBaja(${r.id})" title="Dar de baja">
                    <i class="bi bi-eye-slash"></i>
               </button>`
            : `<button class="btn btn-light border btn-sm text-success" onclick="restaurar(${r.id})" title="Restaurar">
                    <i class="bi bi-arrow-counterclockwise"></i>
               </button>`;

        col.innerHTML = `
            <div class="proyecto-card ${!estaActivo ? 'bg-light border-opacity-50 opacity-75' : ''}">
                <div class="card-top-bar ${!estaActivo ? 'bg-secondary' : ''}"></div>
                <div class="p-3">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <span class="${!estaActivo ? 'badge bg-secondary text-white' : 'badge-vigente'}">
                            <i class="bi ${!estaActivo ? 'bi-eye-slash-fill' : 'bi-award-fill'}"></i>
                            ${!estaActivo ? 'Inactivo' : 'Activo'}
                        </span>
                        <div class="d-flex gap-1">${botonesAccion}</div>
                    </div>
                    <h5 class="fw-bold ${!estaActivo ? 'text-muted text-decoration-line-through' : 'text-dark'} mb-2"
                        style="font-size: 1.05rem; line-height: 1.35;">${r.titulo}</h5>
                    <p class="text-muted small mb-0">${r.descripcion || 'Sin descripción'}</p>
                </div>
            </div>
        `;
        cardsGrid.appendChild(col);
    });
}

function abrirDrawerCrear() {
    document.getElementById('formReconocimiento').reset();
    document.getElementById('formId').value = '';
    document.getElementById('drawerTitulo').innerHTML =
        '<i class="bi bi-plus-circle text-primary"></i> <span>Nuevo Reconocimiento</span>';
    drawerBS.show();
}

function abrirDrawerEditar(id) {
    const r = reconocimientos.find(item => item.id === id);
    if (!r) return;
    document.getElementById('formId').value = r.id;
    document.getElementById('formTitulo').value = r.titulo || '';
    document.getElementById('formDescripcion').value = r.descripcion || '';
    document.getElementById('drawerTitulo').innerHTML =
        `<i class="bi bi-pencil-square text-primary"></i> <span>Editar Reconocimiento #${r.id}</span>`;
    drawerBS.show();
}

document.getElementById('formReconocimiento').addEventListener('submit', async (e) => {
    e.preventDefault();
    const id = document.getElementById('formId').value;
    const payload = {
        titulo: document.getElementById('formTitulo').value.trim(),
        descripcion: document.getElementById('formDescripcion').value.trim()
    };
    try {
        let res;
        if (!id) {
            res = await fetch('/api/v1/admin/reconocimientos', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
        } else {
            res = await fetch(`/api/v1/admin/reconocimientos/${id}`, {
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
        mostrarToast(id ? 'Reconocimiento actualizado.' : 'Reconocimiento creado.');
        cargarDesdeAPI();
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
});

function confirmarBaja(id) {
    idReconocimientoABajar = id;
    modalEliminarBS.show();
}

document.getElementById('btnConfirmarEliminar').addEventListener('click', async () => {
    if (idReconocimientoABajar !== null) {
        try {
            const res = await fetch(`/api/v1/admin/reconocimientos/${idReconocimientoABajar}`, { method: 'DELETE' });
            if (!res.ok) throw new Error('No se pudo dar de baja');
            modalEliminarBS.hide();
            mostrarToast('Reconocimiento dado de baja.', 'danger');
            cargarDesdeAPI();
        } catch (err) {
            mostrarToast(err.message, 'danger');
        } finally {
            idReconocimientoABajar = null;
        }
    }
});

async function restaurar(id) {
    try {
        const res = await fetch(`/api/v1/admin/reconocimientos/${id}/restaurar`, { method: 'PATCH' });
        if (!res.ok) throw new Error('No se pudo restaurar');
        mostrarToast('Reconocimiento restaurado.');
        cargarDesdeAPI();
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
}

function configurarEventos() {
    inputBuscar.addEventListener('input', () => renderizarCards());
    selectOrden.addEventListener('change', () => renderizarCards());
}
