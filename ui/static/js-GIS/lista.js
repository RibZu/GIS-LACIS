/**
 * lista.js - Control de búsqueda, filtrado, paginación en tiempo real y modal de detalles
 * para la vista de Gestión de Integrantes en el Panel de Administración.
 */

document.addEventListener('DOMContentLoaded', function () {
    const buscador = document.getElementById('buscadorInput');
    const filterPills = document.querySelectorAll('.filter-pill');
    const pertenenciaBtns = document.querySelectorAll('.btn-pertenencia');
    const rows = Array.from(document.querySelectorAll('.integrante-row'));
    const noResultsMsg = document.getElementById('noResultsMsg');
    const contadorRegistros = document.getElementById('contadorRegistros');
    const paginationControls = document.getElementById('paginationControls');

    const ITEMS_PER_PAGE = 10;
    let currentPage = 1;
    let activeRole = 'all';
    let activePertenencia = 'all';

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
            const rowName = row.getAttribute('data-nombre') ? row.getAttribute('data-nombre').toLowerCase() : '';
            const tituloElem = row.querySelector('.titulo-integrante');
            const rowTitle = tituloElem ? tituloElem.textContent.toLowerCase() : '';
            const rowRole = row.getAttribute('data-rol');
            const isLacis = row.getAttribute('data-lacis') === 'true';
            const isSoftware = row.getAttribute('data-software') === 'true';

            const matchesSearch = rowName.includes(searchText) || rowTitle.includes(searchText);
            const matchesRole = (activeRole === 'all') || (rowRole === activeRole);

            let matchesPertenencia = true;
            if (activePertenencia === 'lacis') matchesPertenencia = isLacis;
            if (activePertenencia === 'software') matchesPertenencia = isSoftware;

            return matchesSearch && matchesRole && matchesPertenencia;
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
            if (contadorRegistros) contadorRegistros.textContent = '0 integrantes encontrados';
        } else {
            if (noResultsMsg) noResultsMsg.classList.add('d-none');
            const fromNum = startIndex + 1;
            const toNum = Math.min(endIndex, totalFiltered);
            if (contadorRegistros) {
                contadorRegistros.textContent = `Mostrando ${fromNum} a ${toNum} de ${totalFiltered} integrantes (Página ${currentPage} de ${totalPages})`;
            }
        }

        // 4. Renderizar controles
        renderPagination(totalFiltered, totalPages);
    }

    // Eventos
    if (buscador) {
        buscador.addEventListener('input', function() {
            currentPage = 1;
            applyFiltersAndPagination();
        });
    }

    filterPills.forEach(pill => {
        pill.addEventListener('click', function () {
            filterPills.forEach(p => p.classList.remove('active'));
            this.classList.add('active');
            activeRole = this.getAttribute('data-filter');
            currentPage = 1;
            applyFiltersAndPagination();
        });
    });

    pertenenciaBtns.forEach(btn => {
        btn.addEventListener('click', function () {
            pertenenciaBtns.forEach(b => b.classList.remove('active'));
            this.classList.add('active');
            activePertenencia = this.getAttribute('data-pertenencia');
            currentPage = 1;
            applyFiltersAndPagination();
        });
    });

    // Inicializar al cargar
    applyFiltersAndPagination();
});

// Función para abrir el modal con los detalles completos del integrante
function verDetalles(nombre, titulo, rol, descripcion, email, linkedin, pertenencia) {
    const modalNombre = document.getElementById('modalNombre');
    const modalTitulo = document.getElementById('modalTituloJerarquia');
    const modalDescripcion = document.getElementById('modalDescripcion');
    const modalEmail = document.getElementById('modalEmail');
    const modalPertenencia = document.getElementById('modalPertenencia');

    if (modalNombre) modalNombre.textContent = nombre || '';
    if (modalTitulo) modalTitulo.textContent = titulo || '';
    if (modalDescripcion) modalDescripcion.textContent = descripcion || '';
    if (modalEmail) modalEmail.textContent = email || '';
    if (modalPertenencia) modalPertenencia.textContent = pertenencia || '';

    // Iniciales avatar
    const avatarElem = document.getElementById('modalAvatar');
    if (avatarElem && nombre) {
        const partes = nombre.replace(/^(Dr\.|Dra\.|Lic\.|Ing\.|Est\.|Mg\.|Esp\.)\s+/, '').split(' ');
        let iniciales = '';
        if (partes.length >= 2) {
            iniciales = (partes[0][0] + partes[1][0]).toUpperCase();
        } else if (partes.length === 1) {
            iniciales = partes[0].substring(0, 2).toUpperCase();
        }
        avatarElem.textContent = iniciales;
    }

    // Badge Rol
    const rolBadge = document.getElementById('modalRolBadge');
    if (rolBadge) {
        rolBadge.className = 'badge';
        if (rol === 'DIRECTOR') rolBadge.classList.add('badge-rol-director');
        else if (rol === 'INVESTIGADOR') rolBadge.classList.add('badge-rol-investigador');
        else if (rol === 'ASESOR') rolBadge.classList.add('badge-rol-asesor');
        else rolBadge.classList.add('badge-rol-estudiante');
        rolBadge.textContent = rol;
    }

    // LinkedIn
    const linkedinWrapper = document.getElementById('modalLinkedinWrapper');
    const linkedinLink = document.getElementById('modalLinkedin');
    if (linkedinWrapper && linkedinLink) {
        if (linkedin && linkedin.trim() !== '') {
            linkedinWrapper.style.display = 'flex';
            linkedinLink.href = linkedin;
        } else {
            linkedinWrapper.style.display = 'none';
        }
    }

    const modalElem = document.getElementById('modalDetalleIntegrante');
    if (modalElem && typeof bootstrap !== 'undefined') {
        const modal = new bootstrap.Modal(modalElem);
        modal.show();
    }
}

// Exponer en window por compatibilidad con llamadas onclick
window.verDetalles = verDetalles;
