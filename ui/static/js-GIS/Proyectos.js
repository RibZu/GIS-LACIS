document.addEventListener('DOMContentLoaded', () => {
    if (typeof AOS !== 'undefined') {
        AOS.init({ duration: 1000, once: false });
    }
    CargarLista();
    CargarLogros();
    listaResponsive();
});
function obtenerNombres(i){
    fetch('../../static/json-GIS/Nombres.json')
        .then(response =>response.json())
        .then(data=>{
                const titulo = document.getElementById("titulo-proyecto");
                const periodo = document.getElementById("periodo-proyecto");
                const info = document.getElementById("info-proyecto");
                [titulo, periodo, info].forEach(el => {el.classList.remove("blur-in"); void el.offsetWidth; el.classList.add("blur-in");})
            data.forEach(p=> {
                if(p.index===i){
                titulo.innerHTML=p.nom;
                periodo.innerHTML=p.periodo;
                info.innerHTML=p.info;
                }
            });
});
}
function CargarLista(){
    let lista = document.getElementById("lista-proyecto");
    
    fetch('../../static/json-GIS/Nombres.json')
        .then(response => response.json())
        .then(data => {
            let maxIndex = -1; 

            data.forEach(p => {
                const nuevoNom = document.createElement("li");
                nuevoNom.textContent = p.nom;
                nuevoNom.addEventListener("click", () => {
                    obtenerNombres(p.index);
                });
                lista.appendChild(nuevoNom);

                if (p.index > maxIndex) {
                    maxIndex = p.index;
                }
            });

            if (maxIndex !== -1) {
                obtenerNombres(maxIndex);
            }
        });
}
function CargarLogros(){
    fetch('../../static/json-GIS/Premios.json')
        .then(response =>response.json())
        .then(data=>{
            data.forEach(p => {
                const container = document.createElement("div");
                container.className = "col-12 col-sm-6 col-md-4 col-lg-3 mb-4";
                container.setAttribute('data-aos','fade-up');

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
                titulo.innerHTML = p.prem;
                decript.innerHTML = p.desc.replace(/\n/g, '<br>');

                // Anidar correctamente
                cff.appendChild(icono);
                cff.appendChild(titulo);
                cfb.appendChild(decript);

                cfi.appendChild(cff);
                cfi.appendChild(cfb);

                cfp.appendChild(cfi);
                container.appendChild(cfp);

                document.getElementById("logros").appendChild(container);
            });

});
}
let activo = false;

function mostrarEmail() {
  const email = document.getElementById('emailTexto');
  
  if (!activo) {
    email.classList.add('activo');
    activo = true;
  } else {
    email.classList.remove('activo');
    activo = false;
  }
}
function listaResponsive(){
    const mediaQuery = window.matchMedia("(max-width: 900px)");
    if (mediaQuery.matches) {
        //remover elementos
        document.getElementById("listaNom1").remove();
        document.getElementById("listaNom2").remove();
        //crear contenedores del carrucel
        const cs = document.createElement("div");
        cs.className = "carousel slide";
        cs.setAttribute('data-bs-ride', 'carousel');
        cs.id="carouselExampleControls";

        const ci = document.createElement("div");
        ci.className = "carousel-inner";

        //crear botones
        const btnPrev = document.createElement("button");
        btnPrev.className="carousel-control-prev";
        btnPrev.setAttribute('data-bs-target','#carouselExampleControls');
        btnPrev.setAttribute('data-bs-slide','prev');
        btnPrev.setAttribute('type','button');

        const sp1 = document.createElement("span");
        sp1.className="letraAzul fs-3";

        const i1 = document.createElement("i");
        i1.className="fa-solid fa-arrow-left";

        const btnNext = document.createElement("button");
        btnNext.className="carousel-control-next";
        btnNext.setAttribute('data-bs-target','#carouselExampleControls');
        btnNext.setAttribute('data-bs-slide','next');
        btnNext.setAttribute('type','button');
        
        const sp2 = document.createElement("span");
        sp2.className="letraAzul fs-3";
        const i2 = document.createElement("i");
        i2.className="fa-solid fa-arrow-right";

        //añadir elementos
        document.getElementById("Lista").appendChild(cs);
        cs.appendChild(ci);
        cs.appendChild(btnPrev);
         btnPrev.appendChild(sp1);
          sp1.appendChild(i1);

        cs.appendChild(btnNext);
         btnNext.appendChild(sp2);
          sp2.appendChild(i2);
        //crear elementos dinamicos
        fetch('../../static/json-GIS/Nombres.json')
            .then(response =>response.json())
            .then(data=>{
                data.forEach((data, index)=> {

                const nuevoNomR=document.createElement("div");
                if(index==0){
                    nuevoNomR.className="carousel-item active";
                }else{
                    nuevoNomR.className="carousel-item";
                }
                

                const cic=document.createElement("div");
                cic.className="p-4";

                const titulo=document.createElement("h5");
                titulo.className="letraAzul text-center fw-bold fs-7 mt-3";
                const periodo=document.createElement("p");
                periodo.className="text-muted text-center";
                const info=document.createElement("p");
                info.className="mt-3 fs-6";

                titulo.innerHTML=data.nom;
                periodo.innerHTML=data.periodo;
                info.innerHTML=data.info;
                //añadir elementos dinamicos
                nuevoNomR.appendChild(cic);
                cic.appendChild(titulo);
                cic.appendChild(periodo);
                cic.appendChild(info);
                ci.appendChild(nuevoNomR);
                });
            });
    }
}
window.addEventListener("resize", () => {
  if (window.innerWidth > 900) {
    location.reload();
  } else {
    listaResponsive(); 
  }
});

