document.addEventListener('DOMContentLoaded', () => {
    if (typeof AOS !== 'undefined') {
        AOS.init();
    }

    const boton_frontal = document.querySelectorAll('.flip-btn');
    boton_frontal.forEach(boton => {
        boton.addEventListener('click', () => {
            const cambiar_perspectiva = boton.closest('.carta-flipped');
            if (cambiar_perspectiva) {
                if (cambiar_perspectiva.classList.contains('flipped')) {
                    cambiar_perspectiva.classList.remove('flipped');
                } else {
                    cambiar_perspectiva.classList.add("flipped");
                }
            }
        });
    });

    const boton_ver_mas_productos = document.getElementById('btn-ver-mas-productos');
    if (boton_ver_mas_productos) {
        boton_ver_mas_productos.addEventListener('click', () => {
            const ocultos = document.querySelectorAll('.producto-oculto');
            [...ocultos].slice(0, 6).forEach(carta => {
                carta.classList.remove('producto-oculto', 'd-none');
            });
            if (document.querySelectorAll('.producto-oculto').length === 0) {
                boton_ver_mas_productos.closest('.text-center').remove();
            }
            if (typeof AOS !== 'undefined') {
                AOS.refreshHard();
            }
        });
    }
});
