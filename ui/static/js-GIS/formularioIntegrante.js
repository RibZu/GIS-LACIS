document.addEventListener('DOMContentLoaded', () => {
    initGrupoToggle('lacis', 'pertenece_lacis', 'card_lacis', 'contenedor_rol_lacis', 'rol_lacis_id');
    initGrupoToggle('software', 'pertenece_grupo_software', 'card_software', 'contenedor_rol_software', 'rol_software_id');
});

function initGrupoToggle(nombreGrupo, checkId, cardId, contenedorId, selectId) {
    const chk = document.getElementById(checkId);
    const card = document.getElementById(cardId);
    const cont = document.getElementById(contenedorId);
    const sel = document.getElementById(selectId);

    if (!chk || !card || !cont || !sel) {
        return;
    }

    const actualizarEstado = () => {
        if (chk.checked) {
            card.classList.add('active-group');
            card.style.opacity = '1';
            cont.style.display = 'block';
            sel.disabled = false;
            sel.required = true;
        } else {
            card.classList.remove('active-group');
            card.style.opacity = '0.65';
            cont.style.display = 'none';
            sel.disabled = true;
            sel.required = false;
        }
    };

    chk.addEventListener('change', actualizarEstado);

    actualizarEstado();
}

document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');
    if (form) {
        form.addEventListener('submit', () => {
            const btn = form.querySelector('button[type="submit"]');
            if (btn) {
                btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>Guardando...';
                btn.style.pointerEvents = 'none';
            }
        });
    }
});

