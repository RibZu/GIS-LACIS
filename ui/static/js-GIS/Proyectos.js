document.addEventListener('DOMContentLoaded', () => {
    if (typeof AOS !== 'undefined') {
        AOS.init({ duration: 1000, once: false });
    }
    CargarLista();
    CargarLogros();
    CargarColaboradores();
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
    const lista = document.getElementById("lista-proyecto");
    const carruselInner = document.getElementById("carruselProyectosInner");

    fetch('/api/v1/proyectos')
        .then(response => response.json())
        .then(data => {
            cacheProyectos = data || [];

            // 1. Cargar vista Desktop si existe el contenedor
            if (lista) {
                lista.innerHTML = "";
                if (cacheProyectos.length === 0) {
                    const vacio = document.createElement("li");
                    vacio.textContent = "No hay proyectos disponibles actualmente.";
                    vacio.className = "text-muted p-2";
                    lista.appendChild(vacio);
                } else {
                    cacheProyectos.forEach((p, index) => {
                        const nuevoNom = document.createElement("li");
                        nuevoNom.textContent = p.titulo;
                        if (index === 0) nuevoNom.classList.add("active");

                        nuevoNom.addEventListener("click", () => {
                            lista.querySelectorAll("li").forEach(el => el.classList.remove("active"));
                            nuevoNom.classList.add("active");
                            mostrarDetalleProyecto(p);
                        });
                        lista.appendChild(nuevoNom);
                    });

                    mostrarDetalleProyecto(cacheProyectos[0]);
                }
            }

            // 2. Cargar vista Mobile (Carrusel) si existe el contenedor
            if (carruselInner) {
                renderCarruselProyectos(cacheProyectos);
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

                titulo.innerHTML = p.titulo;
                decript.innerHTML = (p.descripcion || "").replace(/\n/g, '<br>');

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

function CargarColaboradores() {
    const contenedor = document.getElementById("colaboradoresContainer");
    if (!contenedor) return;

    fetch('/api/v1/colaboradores')
        .then(response => response.json())
        .then(data => {
            contenedor.innerHTML = "";
            const colaboradores = data || [];
            if (colaboradores.length === 0) {
                return;
            }

            colaboradores.forEach(c => {
                const container = document.createElement("div");
                container.className = "col-12 col-sm-6 col-md-4 col-lg-3 mb-4";
                container.setAttribute('data-aos', 'fade-up');

                const cfp = document.createElement("div");
                cfp.className = "card-flip card-colaborador-chica mx-auto";

                const cfi = document.createElement("div");
                cfi.className = "card-flip-inner";

                const cff = document.createElement("div");
                cff.className = "card-flip-front text-center";

                const img = document.createElement("img");
                img.src = c.logo_url;
                img.alt = "Logo " + (c.descripcion || "colaborador");

                const cfb = document.createElement("div");
                cfb.className = "card-flip-back";

                const desc = document.createElement("p");
                desc.innerHTML = (c.descripcion || "").replace(/\n/g, '<br>');

                cff.appendChild(img);
                cfb.appendChild(desc);

                cfi.appendChild(cff);
                cfi.appendChild(cfb);

                cfp.appendChild(cfi);
                container.appendChild(cfp);

                contenedor.appendChild(container);
            });
        })
        .catch(err => console.error("Error al cargar colaboradores:", err));
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

function renderCarruselProyectos(proyectos) {
    const inner = document.getElementById("carruselProyectosInner");
    const indicador = document.getElementById("carruselIndicadorTexto");
    const carruselEl = document.getElementById("carouselExampleControls");
    if (!inner) return;

    inner.innerHTML = "";

    if (!proyectos || proyectos.length === 0) {
        inner.innerHTML = `
            <div class="carousel-item active">
                <div class="carrusel-item-contenido text-center py-4 text-muted">
                    No hay proyectos disponibles actualmente.
                </div>
            </div>`;
        if (indicador) indicador.textContent = "";
        return;
    }

    proyectos.forEach((item, index) => {
        const itemDiv = document.createElement("div");
        itemDiv.className = index === 0 ? "carousel-item active" : "carousel-item";

        const anios = (item.anio_inicio && item.anio_fin && item.anio_inicio !== item.anio_fin)
            ? `${item.anio_inicio} - ${item.anio_fin}`
            : (item.anio_inicio || item.anio_fin || "");

        const descHtml = (item.descripcion || "").replace(/\n/g, '<br>');

        itemDiv.innerHTML = `
            <div class="carrusel-item-contenido">
                ${anios ? `<div class="mb-2"><span class="carrusel-periodo-badge">${anios}</span></div>` : ""}
                <h4 class="carrusel-titulo-proyecto">${item.titulo || ""}</h4>
                <div class="carrusel-descripcion-proyecto mt-3">${descHtml}</div>
            </div>
        `;
        inner.appendChild(itemDiv);
    });

    if (indicador) {
        indicador.textContent = `Proyecto 1 de ${proyectos.length}`;
    }

    if (carruselEl) {
        if (window.bootstrap && bootstrap.Carousel) {
            bootstrap.Carousel.getOrCreateInstance(carruselEl, {
                interval: false,
                touch: true,
                wrap: true
            });
        }

        // Remover listeners anteriores si los hubiera y escuchar cambios de slide
        carruselEl.removeEventListener('slide.bs.carousel', carruselEl._slideHandler);
        carruselEl._slideHandler = (event) => {
            if (indicador) {
                const nextIndex = typeof event.to === 'number' ? event.to : 0;
                indicador.textContent = `Proyecto ${nextIndex + 1} de ${proyectos.length}`;
            }
        };
        carruselEl.addEventListener('slide.bs.carousel', carruselEl._slideHandler);
    }
}
