document.addEventListener('DOMContentLoaded', () => {
    if (typeof AOS !== 'undefined') {
        AOS.init({ duration: 1000, once: false });
    }
    CargarLista();
    CargarLogros();
    listaResponsive();
});

let cacheProyectos = [];

function mostrarDetalleProyecto(p) {
    const titulo = document.getElementById("titulo-proyecto");
    const periodo = document.getElementById("periodo-proyecto");
    const info = document.getElementById("info-proyecto");
    if (titulo && periodo && info) {
        [titulo, periodo, info].forEach(el => {
            el.classList.remove("blur-in");
            void el.offsetWidth;
            el.classList.add("blur-in");
        });
        titulo.innerHTML = p.titulo || "";
        const anios = (p.anio_inicio && p.anio_fin && p.anio_inicio !== p.anio_fin)
            ? `${p.anio_inicio} - ${p.anio_fin}`
            : (p.anio_inicio || p.anio_fin || "");
        periodo.innerHTML = anios;
        info.innerHTML = (p.descripcion || "").replace(/\n/g, '<br>');
    }
}

function obtenerNombres(id) {
    if (cacheProyectos.length > 0) {
        const found = cacheProyectos.find(p => p.id === id);
        if (found) {
            mostrarDetalleProyecto(found);
            return;
        }
    }
    fetch('/api/v1/proyectos')
        .then(response => response.json())
        .then(data => {
            cacheProyectos = data || [];
            const found = cacheProyectos.find(p => p.id === id);
            if (found) {
                mostrarDetalleProyecto(found);
            }
        })
        .catch(err => console.error("Error al obtener proyecto:", err));
}

function CargarLista() {
    let lista = document.getElementById("lista-proyecto");
    if (!lista) return;

    fetch('/api/v1/proyectos')
        .then(response => response.json())
        .then(data => {
            cacheProyectos = data || [];
            lista.innerHTML = "";
            if (cacheProyectos.length === 0) {
                const vacio = document.createElement("li");
                vacio.textContent = "No hay proyectos disponibles actualmente.";
                vacio.className = "text-muted p-2";
                lista.appendChild(vacio);
                return;
            }

            cacheProyectos.forEach((p, index) => {
                const nuevoNom = document.createElement("li");
                nuevoNom.textContent = p.titulo;
                nuevoNom.addEventListener("click", () => {
                    mostrarDetalleProyecto(p);
                });
                lista.appendChild(nuevoNom);
            });

            // Mostrar el primer proyecto por defecto
            if (cacheProyectos.length > 0) {
                mostrarDetalleProyecto(cacheProyectos[0]);
            }
        })
        .catch(err => console.error("Error al cargar lista de proyectos:", err));
}

function CargarLogros() {
    const logrosContainer = document.getElementById("logros");
    if (!logrosContainer) return;

    fetch('/api/v1/reconocimientos')
        .then(response => response.json())
        .then(data => {
            logrosContainer.innerHTML = "";
            const premios = data || [];
            if (premios.length === 0) {
                return;
            }

            premios.forEach(p => {
                const container = document.createElement("div");
                container.className = "col-12 col-sm-6 col-md-4 col-lg-3 mb-4";
                container.setAttribute('data-aos', 'fade-up');

                const cfp = document.createElement("div");
                cfp.className = "card-flip mx-auto";

                const cfi = document.createElement("div");
                cfi.className = "card-flip-inner";

                const cff = document.createElement("div");
                cff.className = "card-flip-front text-center";

                const icono = document.createElement("i");
                icono.className = "fa-solid fa-medal";

                const titulo = document.createElement("p");

                const cfb = document.createElement("div");
                cfb.className = "card-flip-back";

                const decript = document.createElement("p");

                // Añadir contenido de texto
                titulo.innerHTML = p.titulo;
                decript.innerHTML = (p.descripcion || "").replace(/\n/g, '<br>');

                // Anidar correctamente
                cff.appendChild(icono);
                cff.appendChild(titulo);
                cfb.appendChild(decript);

                cfi.appendChild(cff);
                cfi.appendChild(cfb);

                cfp.appendChild(cfi);
                container.appendChild(cfp);

                logrosContainer.appendChild(container);
            });
        })
        .catch(err => console.error("Error al cargar logros:", err));
}

let activo = false;

function mostrarEmail() {
    const email = document.getElementById('emailTexto');
    if (!email) return;

    if (!activo) {
        email.classList.add('activo');
        activo = true;
    } else {
        email.classList.remove('activo');
        activo = false;
    }
}

function listaResponsive() {
    const mediaQuery = window.matchMedia("(max-width: 900px)");
    if (mediaQuery.matches) {
        const nom1 = document.getElementById("listaNom1");
        const nom2 = document.getElementById("listaNom2");
        if (nom1) nom1.remove();
        if (nom2) nom2.remove();

        const listaContainer = document.getElementById("Lista");
        if (!listaContainer) return;

        // Crear contenedores del carrusel
        const cs = document.createElement("div");
        cs.className = "carousel slide";
        cs.setAttribute('data-bs-ride', 'carousel');
        cs.id = "carouselExampleControls";

        const ci = document.createElement("div");
        ci.className = "carousel-inner";

        // Crear botones
        const btnPrev = document.createElement("button");
        btnPrev.className = "carousel-control-prev";
        btnPrev.setAttribute('data-bs-target', '#carouselExampleControls');
        btnPrev.setAttribute('data-bs-slide', 'prev');
        btnPrev.setAttribute('type', 'button');

        const sp1 = document.createElement("span");
        sp1.className = "letraAzul fs-3";

        const i1 = document.createElement("i");
        i1.className = "fa-solid fa-arrow-left";

        const btnNext = document.createElement("button");
        btnNext.className = "carousel-control-next";
        btnNext.setAttribute('data-bs-target', '#carouselExampleControls');
        btnNext.setAttribute('data-bs-slide', 'next');
        btnNext.setAttribute('type', 'button');

        const sp2 = document.createElement("span");
        sp2.className = "letraAzul fs-3";
        const i2 = document.createElement("i");
        i2.className = "fa-solid fa-arrow-right";

        // Añadir elementos
        listaContainer.appendChild(cs);
        cs.appendChild(ci);
        cs.appendChild(btnPrev);
        btnPrev.appendChild(sp1);
        sp1.appendChild(i1);

        cs.appendChild(btnNext);
        btnNext.appendChild(sp2);
        sp2.appendChild(i2);

        // Crear elementos dinámicos desde API
        fetch('/api/v1/proyectos')
            .then(response => response.json())
            .then(data => {
                const proyectos = data || [];
                if (proyectos.length === 0) return;

                proyectos.forEach((item, index) => {
                    const nuevoNomR = document.createElement("div");
                    nuevoNomR.className = index === 0 ? "carousel-item active" : "carousel-item";

                    const cic = document.createElement("div");
                    cic.className = "p-4";

                    const titulo = document.createElement("h5");
                    titulo.className = "letraAzul text-center fw-bold fs-7 mt-3";
                    const periodo = document.createElement("p");
                    periodo.className = "text-muted text-center";
                    const info = document.createElement("p");
                    info.className = "mt-3 fs-6";

                    titulo.innerHTML = item.titulo;
                    const anios = (item.anio_inicio && item.anio_fin && item.anio_inicio !== item.anio_fin)
                        ? `${item.anio_inicio} - ${item.anio_fin}`
                        : (item.anio_inicio || item.anio_fin || "");
                    periodo.innerHTML = anios;
                    info.innerHTML = (item.descripcion || "").replace(/\n/g, '<br>');

                    nuevoNomR.appendChild(cic);
                    cic.appendChild(titulo);
                    cic.appendChild(periodo);
                    cic.appendChild(info);
                    ci.appendChild(nuevoNomR);
                });
            })
            .catch(err => console.error("Error en carrusel responsive de proyectos:", err));
    }
}

window.addEventListener("resize", () => {
    if (window.innerWidth > 900) {
        location.reload();
    } else {
        listaResponsive();
    }
});
