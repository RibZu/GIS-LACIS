/**
 * formularioIntegrante.js
 * Control de interactividad para los formularios de Crear y Editar Integrante en el panel Admin.
 * Gestiona el toggle dinámico de las tarjetas de pertenencia (LaCIS y Grupo Software)
 * y la obligatoriedad de los selectores de rol.
 */

document.addEventListener('DOMContentLoaded', () => {
    initGrupoToggle('lacis', 'pertenece_lacis', 'card_lacis', 'contenedor_rol_lacis', 'rol_lacis_id');
    initGrupoToggle('software', 'pertenece_grupo_software', 'card_software', 'contenedor_rol_software', 'rol_software_id');
});

/**
 * Inicializa los eventos y estado inicial para una tarjeta de grupo.
 */
function initGrupoToggle(nombreGrupo, checkId, cardId, contenedorId, selectId) {
    const chk = document.getElementById(checkId);
    const card = document.getElementById(cardId);
    const cont = document.getElementById(contenedorId);
    const sel = document.getElementById(selectId);

    if (!chk || !card || !cont || !sel) {
        return;
    }

    // Función que actualiza la vista y atributos según el estado del switch
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

    // Escuchar cambios en el switch
    chk.addEventListener('change', actualizarEstado);

    // Ejecutar al cargar la página
    actualizarEstado();
}
