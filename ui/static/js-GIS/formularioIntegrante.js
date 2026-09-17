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
