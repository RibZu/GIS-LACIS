document.addEventListener('submit', function (evento) {
    var formulario = evento.target;
    if (!formulario.matches || !formulario.matches('form[data-envio-unico]')) {
        return;
    }
    if (formulario.dataset.enviado === 'true') {
        evento.preventDefault();
        return;
    }
    formulario.dataset.enviado = 'true';
    setTimeout(function () {
        formulario.querySelectorAll('button[type="submit"]').forEach(function (boton) {
            boton.disabled = true;
        });
    }, 0);
});

window.addEventListener('pageshow', function (evento) {
    if (!evento.persisted) {
        return;
    }
    document.querySelectorAll('form[data-envio-unico]').forEach(function (formulario) {
        delete formulario.dataset.enviado;
        formulario.querySelectorAll('button[type="submit"]').forEach(function (boton) {
            boton.disabled = false;
        });
    });
});
