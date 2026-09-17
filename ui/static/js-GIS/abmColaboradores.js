let colaboradores = [];
let idColaboradorABajar = null;

const cardsGrid = document.getElementById('cardsGrid');
const emptyState = document.getElementById('emptyState');
const inputBuscar = document.getElementById('inputBuscar');
const selectOrden = document.getElementById('selectOrden');
const formLogo = document.getElementById('formLogo');
const logoPreview = document.getElementById('logoPreview');
const logoPreviewPlaceholder = document.getElementById('logoPreviewPlaceholder');

const drawerEl = document.getElementById('drawerColaborador');
const drawerBS = new bootstrap.Offcanvas(drawerEl);
const modalEliminarBS = new bootstrap.Modal(document.getElementById('modalEliminar'));
const toastBS = new bootstrap.Toast(document.getElementById('liveToast'));

document.addEventListener('DOMContentLoaded', () => {
    cargarDesdeAPI();
    configurarEventos();
});

async function cargarDesdeAPI() {
    try {
        const res = await fetch('/api/v1/admin/colaboradores-todos');
        if (res.status === 401 || res.status === 403) {
            window.location.href = '/login';
            return;
        }
        if (!res.ok) throw new Error('Error al cargar datos');
        colaboradores = await res.json() || [];
        renderizarCards();
    } catch (err) {
        mostrarToast('Error al conectar con el servidor', 'danger');
    }
}

function mostrarToast(mensaje, tipo = 'success') {
    const toastEl = document.getElementById('liveToast');
    toastEl.className = `toast align-items-center text-bg-${tipo} border-0 shadow`;
    document.getElementById('toastMessage').textContent = mensaje;
    toastBS.show();
}

function obtenerColaboradoresFiltrados() {
    const texto = (inputBuscar.value || '').toLowerCase();
    const orden = selectOrden.value;
    let resultado = colaboradores.filter(c => {
        return (c.descripcion || '').toLowerCase().includes(texto);
    });
    resultado.sort((a, b) => {
        if (orden === 'recientes') return b.id - a.id;
        if (orden === 'antiguos') return a.id - b.id;
        return 0;
    });
    return resultado;
}

function renderizarCards() {
    const listado = obtenerColaboradoresFiltrados();
    cardsGrid.innerHTML = '';
    if (listado.length === 0) {
        emptyState.classList.remove('d-none');
        return;
    }
    emptyState.classList.add('d-none');

    listado.forEach(c => {
        const estaActivo = c.activo !== false;
        const col = document.createElement('div');
        col.className = 'col-12 col-sm-6 col-md-4 col-lg-3';

        const botonesAccion = estaActivo
            ? `<button class="btn btn-light border btn-sm shadow-sm" onclick="abrirDrawerEditar(${c.id})" title="Editar">
                    <i class="bi bi-pencil text-primary"></i>
               </button>
               <button class="btn btn-light border text-danger btn-sm shadow-sm" onclick="confirmarBaja(${c.id})" title="Dar de baja">
                    <i class="bi bi-eye-slash"></i>
               </button>`
            : `<button class="btn btn-light border btn-sm text-success shadow-sm" onclick="restaurar(${c.id})" title="Restaurar colaborador">
                    <i class="bi bi-arrow-counterclockwise"></i> Restaurar
               </button>`;

        const logoHtml = c.logo_url
            ? `<img src="${c.logo_url}" alt="Logo Colaborador" class="img-fluid" onerror="this.onerror=null; this.src='/static/assets/img-GIS/favicon-text.png';">`
            : `<div class="logo-placeholder"><i class="bi bi-building"></i></div>`;

        col.innerHTML = `
            <div class="colaborador-card-wrapper ${!estaActivo ? 'inactivo' : ''}">
                <div class="card-header-actions">
                    <span class="${estaActivo ? 'badge-vigente' : 'badge-inactivo'}">
                        <i class="bi ${estaActivo ? 'bi-check-circle-fill' : 'bi-eye-slash-fill'}"></i>
                        ${estaActivo ? 'Activo' : 'Inactivo'}
                    </span>
                    <div class="d-flex gap-1">
                        ${botonesAccion}
                    </div>
                </div>

                <!-- Tarjeta Flip interactiva al estilo de proyectos.html -->
                <div class="card-flip" title="Pasa el cursor o haz clic para girar">
                    <div class="card-flip-inner">
                        <!-- Frente: Logo -->
                        <div class="card-flip-front">
                            ${logoHtml}
                            <span class="flip-hint"><i class="bi bi-arrow-repeat"></i> Girar para ver reseña</span>
                        </div>
                        <!-- Dorso: Reseña / Descripción en gradiente azul -->
                        <div class="card-flip-back">
                            <p>${c.descripcion || 'Sin descripción disponible.'}</p>
                        </div>
                    </div>
                </div>
            </div>
        `;
        cardsGrid.appendChild(col);
    });
}

function abrirDrawerCrear() {
    document.getElementById('formColaborador').reset();
    document.getElementById('formId').value = '';
    actualizarPreviewLogo(null);
    document.getElementById('drawerTitulo').innerHTML =
        '<i class="bi bi-plus-circle text-primary"></i> <span>Nuevo Colaborador</span>';
    drawerBS.show();
}

function abrirDrawerEditar(id) {
    const c = colaboradores.find(item => item.id === id);
    if (!c) return;
    document.getElementById('formId').value = c.id;
    document.getElementById('formDescripcion').value = c.descripcion || '';

    // Clear file input
    document.getElementById('formLogo').value = '';

    // Show current logo in preview if available
    if (c.logo_url) {
        logoPreview.src = c.logo_url;
        logoPreview.classList.remove('d-none');
        logoPreviewPlaceholder.classList.add('d-none');
    } else {
        actualizarPreviewLogo(null);
    }

    document.getElementById('drawerTitulo').innerHTML =
        `<i class="bi bi-pencil-square text-primary"></i> <span>Editar: #${c.id}</span>`;
    drawerBS.show();
}

function actualizarPreviewLogo(file) {
    if (file) {
        const reader = new FileReader();
        reader.onload = function (e) {
            logoPreview.src = e.target.result;
            logoPreview.classList.remove('d-none');
            logoPreviewPlaceholder.classList.add('d-none');
        };
        reader.readAsDataURL(file);
    } else {
        logoPreview.src = '';
        logoPreview.classList.add('d-none');
        logoPreviewPlaceholder.textContent = 'Selecciona una imagen para ver la previsualización';
        logoPreviewPlaceholder.classList.remove('d-none');
    }
}

document.getElementById('formColaborador').addEventListener('submit', async (e) => {
    e.preventDefault();
    const id = document.getElementById('formId').value;

    const formData = new FormData();
    formData.append('descripcion', document.getElementById('formDescripcion').value.trim());

    const logoFile = document.getElementById('formLogo').files[0];
    if (logoFile) {
        formData.append('logo', logoFile);
    }

    try {
        let res;
        if (!id) {
            if (!logoFile) {
                mostrarToast('Debe seleccionar una imagen para el colaborador.', 'danger');
                return;
            }
            res = await fetch('/api/v1/admin/colaboradores', {
                method: 'POST',
                body: formData
            });
        } else {
            res = await fetch(`/api/v1/admin/colaboradores/${id}`, {
                method: 'PUT',
                body: formData
            });
        }
        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.error || 'Error al guardar');
        }
        drawerBS.hide();
        mostrarToast(id ? 'Colaborador actualizado correctamente.' : 'Colaborador creado exitosamente.');
        cargarDesdeAPI();
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
});

function confirmarBaja(id) {
    idColaboradorABajar = id;
    modalEliminarBS.show();
}

document.getElementById('btnConfirmarEliminar').addEventListener('click', async () => {
    if (idColaboradorABajar !== null) {
        try {
            const res = await fetch(`/api/v1/admin/colaboradores/${idColaboradorABajar}`, { method: 'DELETE' });
            if (!res.ok) throw new Error('No se pudo dar de baja el colaborador');
            modalEliminarBS.hide();
            mostrarToast('Colaborador dado de baja.', 'danger');
            cargarDesdeAPI();
        } catch (err) {
            mostrarToast(err.message, 'danger');
        } finally {
            idColaboradorABajar = null;
        }
    }
});

async function restaurar(id) {
    try {
        const res = await fetch(`/api/v1/admin/colaboradores/${id}/restaurar`, { method: 'PATCH' });
        if (!res.ok) throw new Error('No se pudo restaurar el colaborador');
        mostrarToast('Colaborador restaurado exitosamente.');
        cargarDesdeAPI();
    } catch (err) {
        mostrarToast(err.message, 'danger');
    }
}

function configurarEventos() {
    inputBuscar.addEventListener('input', () => renderizarCards());
    selectOrden.addEventListener('change', () => renderizarCards());
    
    formLogo.addEventListener('change', (e) => {
        if (e.target.files && e.target.files[0]) {
            actualizarPreviewLogo(e.target.files[0]);
        } else {
            actualizarPreviewLogo(null);
        }
    });
}
