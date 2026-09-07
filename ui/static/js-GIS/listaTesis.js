/**
 * listaTesis.js - Control de búsqueda reactiva, filtrado por nivel/carrera,
 * paginación en tiempo real y modal de detalles para la vista de Gestión de Tesis.
 * Todo el código JS está centralizado aquí; cero JS en la vista HTML.
 */

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

    function renderPagination(totalFiltered, totalPages) {
        if (!paginationControls) return;
        paginationControls.innerHTML = '';
        if (totalFiltered === 0 || totalPages <= 1) {
            if (totalPages <= 1 && totalFiltered > 0) {
                paginationControls.innerHTML = `<li class="page-item active"><span class="page-link">1</span></li>`;
            }
            return;
        }

        // Botón Anterior
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

        // Botones Numéricos
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

        // Botón Siguiente
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

        // 1. Filtrar filas
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

        // 2. Ocultar todas y mostrar solo la página activa
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

        // 3. Texto del contador y mensaje vacío
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

        // 4. Renderizar controles de paginación
        renderPagination(totalFiltered, totalPages);
    }

    // Eventos de búsqueda y filtros
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

    // Delegación de eventos para botones de "Ver Detalles"
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

    // Inicializar al cargar
    applyFiltersAndPagination();
});

// Función para abrir el modal con la ficha completa de la Tesis
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

    // Badge Nivel
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

    // Palabras clave (chips)
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

    // Enlace a PDF
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

window.verDetallesTesis = verDetallesTesis;
