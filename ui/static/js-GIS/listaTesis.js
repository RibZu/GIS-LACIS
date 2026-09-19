document.addEventListener('DOMContentLoaded', function () {
    const buscador = document.getElementById('buscadorTesisInput');
    const filterPills = document.querySelectorAll('.filter-pill-tesis');
    const carreraBtns = document.querySelectorAll('.btn-carrera-filter');
    const rows = Array.from(document.querySelectorAll('.tesis-row'));
    const noResultsMsg = document.getElementById('noResultsMsg');
    const contadorRegistros = document.getElementById('contadorRegistros');
    const paginationControls = document.getElementById('paginationControls');

    const ITEMS_PER_PAGE = 10;
    let currentPage = 1;
    let activeNivel = 'all';
    let activeCarrera = 'all';

   
    const urlParams = new URLSearchParams(window.location.search);
    const status = urlParams.get('status') || (document.body ? document.body.dataset.status : '');
    if (status && status.trim() !== '') {
        mostrarModalNotificacionTesis(status.trim());
        if (document.body) document.body.removeAttribute('data-status');
        window.history.replaceState({}, document.title, window.location.pathname);
    }

    function renderPagination(totalFiltered, totalPages) {
        if (!paginationControls) return;
        paginationControls.innerHTML = '';
        if (totalFiltered === 0 || totalPages <= 1) {
            if (totalPages <= 1 && totalFiltered > 0) {
                paginationControls.innerHTML = `<li class="page-item active"><span class="page-link">1</span></li>`;
            }
            return;
        }

        const prevLi = document.createElement('li');
        prevLi.className = `page-item ${currentPage === 1 ? 'disabled' : ''}`;
        prevLi.innerHTML = `<a class="page-link" href="#" aria-label="Anterior"><i class="bi bi-chevron-left"></i> Anterior</a>`;
        prevLi.addEventListener('click', function (e) {
            e.preventDefault();
            if (currentPage > 1) {
                currentPage--;
                applyFiltersAndPagination();
            }
        });
        paginationControls.appendChild(prevLi);

        for (let i = 1; i <= totalPages; i++) {
            const pageLi = document.createElement('li');
            pageLi.className = `page-item ${i === currentPage ? 'active' : ''}`;
            pageLi.innerHTML = `<a class="page-link" href="#">${i}</a>`;
            pageLi.addEventListener('click', function (e) {
                e.preventDefault();
                currentPage = i;
                applyFiltersAndPagination();
            });
            paginationControls.appendChild(pageLi);
        }

        const nextLi = document.createElement('li');
        nextLi.className = `page-item ${currentPage === totalPages ? 'disabled' : ''}`;
        nextLi.innerHTML = `<a class="page-link" href="#" aria-label="Siguiente">Siguiente <i class="bi bi-chevron-right"></i></a>`;
        nextLi.addEventListener('click', function (e) {
            e.preventDefault();
            if (currentPage < totalPages) {
                currentPage++;
                applyFiltersAndPagination();
            }
        });
        paginationControls.appendChild(nextLi);
    }

    function applyFiltersAndPagination() {
        if (!buscador) return;
        const searchText = buscador.value.toLowerCase().trim();

        const matchingRows = rows.filter(row => {
            const rowSearch = row.getAttribute('data-search') ? row.getAttribute('data-search').toLowerCase() : '';
            const rowNivel = (row.getAttribute('data-nivel') || '').toLowerCase();
            const rowCarrera = (row.getAttribute('data-carrera') || '').toLowerCase();

            const matchesSearch = searchText === '' || rowSearch.includes(searchText);

            let matchesNivel = true;
            if (activeNivel !== 'all') {
                matchesNivel = rowNivel.includes(activeNivel.toLowerCase());
            }

            let matchesCarrera = true;
            if (activeCarrera !== 'all') {
                matchesCarrera = rowCarrera.includes(activeCarrera.toLowerCase());
            }

            return matchesSearch && matchesNivel && matchesCarrera;
        });

        const totalFiltered = matchingRows.length;
        const totalPages = Math.ceil(totalFiltered / ITEMS_PER_PAGE) || 1;

        if (currentPage > totalPages) {
            currentPage = totalPages;
        }

        const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
        const endIndex = startIndex + ITEMS_PER_PAGE;

        rows.forEach(row => {
            row.style.display = 'none';
        });

        matchingRows.forEach((row, index) => {
            if (index >= startIndex && index < endIndex) {
                row.style.display = '';
            } else {
                row.style.display = 'none';
            }
        });

        if (totalFiltered === 0) {
            if (noResultsMsg) noResultsMsg.classList.remove('d-none');
            if (contadorRegistros) contadorRegistros.textContent = '0 tesis encontradas';
        } else {
            if (noResultsMsg) noResultsMsg.classList.add('d-none');
            const fromNum = startIndex + 1;
            const toNum = Math.min(endIndex, totalFiltered);
            if (contadorRegistros) {
                contadorRegistros.textContent = `Mostrando ${fromNum} a ${toNum} de ${totalFiltered} tesis (Página ${currentPage} de ${totalPages})`;
            }
        }

        renderPagination(totalFiltered, totalPages);
    }

    if (buscador) {
        buscador.addEventListener('input', function () {
            currentPage = 1;
            applyFiltersAndPagination();
        });
    }

    filterPills.forEach(pill => {
        pill.addEventListener('click', function () {
            filterPills.forEach(p => p.classList.remove('active'));
            this.classList.add('active');
            activeNivel = this.getAttribute('data-filter') || 'all';
            currentPage = 1;
            applyFiltersAndPagination();
        });
    });

    carreraBtns.forEach(btn => {
        btn.addEventListener('click', function () {
            carreraBtns.forEach(b => b.classList.remove('active'));
            this.classList.add('active');
            activeCarrera = this.getAttribute('data-carrera') || 'all';
            currentPage = 1;
            applyFiltersAndPagination();
        });
    });

    document.querySelectorAll('.btn-action-ver-tesis').forEach(btn => {
        btn.addEventListener('click', function () {
            const titulo = this.getAttribute('data-titulo') || '';
            const nivel = this.getAttribute('data-nivel') || '';
            const carrera = this.getAttribute('data-carrera') || '';
            const anio = this.getAttribute('data-anio') || '';
            const autor = this.getAttribute('data-autor') || '';
            const director = this.getAttribute('data-director') || '';
            const coodirector = this.getAttribute('data-coodirector') || '';
            const palabrasClave = this.getAttribute('data-palabras') || '';
            const resumen = this.getAttribute('data-resumen') || '';
            const pdf = this.getAttribute('data-pdf') || '';

            verDetallesTesis(titulo, nivel, carrera, anio, autor, director, coodirector, palabrasClave, resumen, pdf);
        });
    });

    applyFiltersAndPagination();
});

function verDetallesTesis(titulo, nivel, carrera, anio, autor, director, coodirector, palabrasClave, resumen, pdf) {
    const modalTitulo = document.getElementById('modalTesisTitulo');
    const modalNivelBadge = document.getElementById('modalTesisNivelBadge');
    const modalCarrera = document.getElementById('modalTesisCarrera');
    const modalAnio = document.getElementById('modalTesisAnio');
    const modalAutor = document.getElementById('modalTesisAutor');
    const modalDirector = document.getElementById('modalTesisDirector');
    const modalCoodirector = document.getElementById('modalTesisCoodirector');
    const modalCoodirectorWrapper = document.getElementById('modalCoodirectorWrapper');
    const modalPalabrasContainer = document.getElementById('modalTesisPalabras');
    const modalResumen = document.getElementById('modalTesisResumen');
    const modalPdfWrapper = document.getElementById('modalTesisPdfWrapper');
    const modalPdfLink = document.getElementById('modalTesisPdfLink');

    if (modalTitulo) modalTitulo.textContent = titulo || 'Sin título';
    if (modalCarrera) modalCarrera.textContent = carrera || 'No especificada';
    if (modalAnio) modalAnio.textContent = (anio && anio !== '0') ? anio : 'Año no especificado';
    if (modalAutor) modalAutor.textContent = autor || 'No especificado';
    if (modalDirector) modalDirector.textContent = director || 'No especificado';

    if (modalCoodirectorWrapper && modalCoodirector) {
        if (coodirector && coodirector.trim() !== '') {
            modalCoodirectorWrapper.style.display = 'block';
            modalCoodirector.textContent = coodirector;
        } else {
            modalCoodirectorWrapper.style.display = 'none';
        }
    }

    if (modalResumen) {
        modalResumen.textContent = resumen || 'No se ha cargado un resumen o abstract para este trabajo final.';
    }

    if (modalNivelBadge) {
        modalNivelBadge.className = 'badge';
        const nivelNorm = (nivel || '').toLowerCase();
        if (nivelNorm.includes('doc')) {
            modalNivelBadge.classList.add('badge-rol-director');
        } else if (nivelNorm.includes('mae')) {
            modalNivelBadge.classList.add('badge-rol-investigador');
        } else if (nivelNorm.includes('esp')) {
            modalNivelBadge.classList.add('badge-rol-asesor');
        } else {
            modalNivelBadge.classList.add('badge-rol-estudiante');
        }
        modalNivelBadge.textContent = nivel || 'Posgrado';
    }

    if (modalPalabrasContainer) {
        modalPalabrasContainer.innerHTML = '';
        if (palabrasClave && palabrasClave.trim() !== '') {
            const palabras = palabrasClave.split(/[,;]+/);
            palabras.forEach(p => {
                const tag = p.trim();
                if (tag) {
                    const span = document.createElement('span');
                    span.className = 'badge-pertenencia me-1 mb-1';
                    span.textContent = tag;
                    modalPalabrasContainer.appendChild(span);
                }
            });
        } else {
            modalPalabrasContainer.innerHTML = '<span class="text-muted small">Sin palabras clave registradas</span>';
        }
    }

    if (modalPdfWrapper && modalPdfLink) {
        if (pdf && pdf.trim() !== '') {
            modalPdfWrapper.style.display = 'block';
            modalPdfLink.href = pdf;
        } else {
            modalPdfWrapper.style.display = 'none';
        }
    }

    const modalElem = document.getElementById('modalDetalleTesis');
    if (modalElem && typeof bootstrap !== 'undefined') {
        const modal = new bootstrap.Modal(modalElem);
        modal.show();
    }
}

function mostrarModalNotificacionTesis(status) {
    if (typeof Swal !== 'undefined') {
        let titulo = '¡Guardado Correctamente!';
        let texto = 'La tesis ha sido registrada exitosamente en el sistema.';
        let icono = 'success';

        if (status === 'editado') {
            titulo = '¡Editado Correctamente!';
            texto = 'Los datos de la tesis fueron actualizados con éxito.';
        } else if (status === 'eliminado') {
            titulo = '¡Eliminado Correctamente!';
            texto = 'La tesis ha sido eliminada del sistema con éxito.';
        } else if (status === 'error') {
            titulo = 'Ocurrió un error';
            texto = 'No se pudo completar la operación en el sistema.';
            icono = 'error';
        }

        Swal.fire({
            icon: icono,
            title: titulo,
            text: texto,
            confirmButtonText: '<i class="bi bi-check2-circle me-1"></i> Aceptar',
            buttonsStyling: false,
            customClass: {
                popup: 'custom-swal-popup',
                title: 'custom-swal-title',
                htmlContainer: 'custom-swal-text',
                confirmButton: 'custom-swal-confirm'
            }
        });
        return;
    }

    // Fallback nativo: Modal Bootstrap integrado
    const modalEl = document.getElementById('modalNotificacion');
    const iconEl = document.getElementById('notifIcon');
    const tituloEl = document.getElementById('notifTitulo');
    const mensajeEl = document.getElementById('notifMensaje');

    if (!modalEl) return;

    if (status === 'guardado' || status === 'creado') {
        if (iconEl) iconEl.className = 'bi bi-check-circle-fill text-success';
        if (tituloEl) tituloEl.textContent = '¡Guardado Correctamente!';
        if (mensajeEl) mensajeEl.textContent = 'La tesis ha sido registrada exitosamente en el sistema.';
    } else if (status === 'editado') {
        if (iconEl) iconEl.className = 'bi bi-check-circle-fill text-success';
        if (tituloEl) tituloEl.textContent = '¡Editado Correctamente!';
        if (mensajeEl) mensajeEl.textContent = 'Los datos de la tesis fueron actualizados con éxito.';
    } else if (status === 'eliminado') {
        if (iconEl) iconEl.className = 'bi bi-check-circle-fill text-success';
        if (tituloEl) tituloEl.textContent = '¡Eliminado Correctamente!';
        if (mensajeEl) mensajeEl.textContent = 'La tesis ha sido eliminada del sistema con éxito.';
    } else if (status === 'error') {
        if (iconEl) iconEl.className = 'bi bi-x-circle-fill text-danger';
        if (tituloEl) tituloEl.textContent = 'Ocurrió un error';
        if (mensajeEl) mensajeEl.textContent = 'No se pudo completar la operación en el sistema.';
    }

    if (typeof bootstrap !== 'undefined') {
        const bsModal = new bootstrap.Modal(modalEl);
        bsModal.show();
    }
}

function confirmarEliminarTesis(id, titulo) {
    if (typeof Swal !== 'undefined') {
        Swal.fire({
            title: '¿Eliminar tesis?',
            html: `¿Estás seguro de que deseas eliminar la tesis <b>"${titulo || 'este registro'}"</b>?<br><small class="text-muted">Esta acción no se puede deshacer.</small>`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: '<i class="bi bi-trash-fill me-1"></i> Sí, eliminar',
            cancelButtonText: 'Cancelar',
            buttonsStyling: false,
            customClass: {
                popup: 'custom-swal-popup',
                title: 'custom-swal-title',
                htmlContainer: 'custom-swal-text',
                confirmButton: 'custom-swal-confirm btn-danger-confirm',
                cancelButton: 'custom-swal-cancel'
            },
            reverseButtons: true,
            focusCancel: true
        }).then((result) => {
            if (result.isConfirmed) {
                window.location.href = `/admin/borrar-tesis?id=${id}`;
            }
        });
        return;
    }

    const modalEl = document.getElementById('modalConfirmarEliminar');
    const textoEl = document.getElementById('eliminarModalTexto');
    const btnAccion = document.getElementById('btnConfirmarEliminarAccion');

    if (textoEl) {
        textoEl.innerHTML = `¿Estás seguro de que deseas eliminar la tesis <b>"${titulo || 'este registro'}"</b>?<br><small class="text-muted">Esta acción no se puede deshacer.</small>`;
    }
    if (btnAccion) {
        btnAccion.href = `/admin/borrar-tesis?id=${id}`;
    }

    if (modalEl && typeof bootstrap !== 'undefined') {
        const bsModal = new bootstrap.Modal(modalEl);
        bsModal.show();
    } else if (confirm(`¿Está seguro de que desea eliminar la tesis '${titulo}'?`)) {
        window.location.href = `/admin/borrar-tesis?id=${id}`;
    }
}

// Delegación de eventos para clicks de eliminación
document.addEventListener('click', function (e) {
    const btn = e.target.closest('.btn-eliminar-item');
    if (btn && btn.dataset.tipo === 'tesis') {
        e.preventDefault();
        const id = btn.dataset.id;
        const nombre = btn.dataset.nombre;
        confirmarEliminarTesis(id, nombre);
    }
});

window.verDetallesTesis = verDetallesTesis;
window.confirmarEliminarTesis = confirmarEliminarTesis;
window.mostrarModalNotificacionTesis = mostrarModalNotificacionTesis;
