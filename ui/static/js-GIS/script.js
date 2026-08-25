document.addEventListener('DOMContentLoaded', function() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('is-visible');
            }
        });
    }, { threshold: 0.1 });

    document.querySelectorAll('.animacion-al-scroll').forEach(el => {
        observer.observe(el);
    });

    const principalDiv = document.getElementById('principal');

    if (principalDiv) {

        const sentinel = document.createElement('div');
        sentinel.style.position = 'absolute';
        sentinel.style.top = '0';
        sentinel.style.height = '1px';
        sentinel.style.width = '100%';
        principalDiv.appendChild(sentinel);

        const bgObserver = new IntersectionObserver(([entry]) => {

            if (!entry.isIntersecting && entry.boundingClientRect.top < 0) {
                principalDiv.classList.add('fondo-azul');
            } else {
                principalDiv.classList.remove('fondo-azul');
            }
        });
        bgObserver.observe(sentinel);
    }

    const tarjetas = document.querySelectorAll('.tarjeta-click');

    tarjetas.forEach(tarjeta => {
        tarjeta.addEventListener('click', function(e) {
            // Si el usuario hace clic en el enlace de Ver Detalle O en el enlace Mailto, NO giramos la tarjeta
            if (e.target.tagName === 'A' || e.target.closest('a')) {
                return; 
            }
            // Si hace clic en cualquier otra parte del cuerpo trasero/frontal, se voltea
            this.classList.toggle('flipped');
        });
    });
});