
let claseTituloIntegrantes = document.querySelector(".titulo-integrantes");
let claseTituloDirector = document.querySelector(".titulo-director");


function AjustarColumnasDirectoresDeLinea(cantidad) {
  const contenedor = document.getElementById("contenedor-directores-deLinea");
  contenedor.classList.remove("row-cols-lg-2", "row-cols-lg-3");
  contenedor.classList.add(cantidad >= 3 ? "row-cols-lg-3" : "row-cols-lg-2");
}

function filtrarCartas(clase) {
  
  contenedorPrincipalDeDirector.innerHTML = "";
  contenedorPrincipalDeDirectoresDeLinea.innerHTML = "";
  contenedorPrincipalDeIntegrantes.innerHTML = "";

  let filtradosIntegrantes;
  let filtradosDirectores;
  let filtradosDirector;


  if(clase==="todas"){
     filtradosIntegrantes=todos_los_integrantes;
     filtradosDirectores=todos_los_directores;
    filtradosDirector=director;
  }else if(clase==="lacis"){

    filtradosDirector = directorLacis;
    filtradosDirectores = todos_los_directoresLacis;
    filtradosIntegrantes= todos_los_integrantes_lacis;

  
   } else{


    filtradosIntegrantes= todos_los_integrantes.filter(p=> p.clase_rol===clase);
    filtradosDirectores=todos_los_directores.filter(p=> p.clase_rol===clase);
    filtradosDirector=director.filter(p => p.clase_rol===clase);

  }

  
  

 
  CrearCartaIntegrantes(filtradosIntegrantes);
  CrearCartaDirectores(filtradosDirectores);
  CrearCartaDirector(filtradosDirector);
  
  AjustarColumnasDirectoresDeLinea(filtradosDirectores.length); 
 
  AOS.refreshHard(); /* se actualizan las animaciones cuando se filtra */
}


let filtroBoton=document.querySelectorAll(".filtro-btn");

filtroBoton.forEach(boton => {

  boton.addEventListener("click", () => {

    let id = boton.id;
    let clase = "";
    

    switch (id) {
      case "todos":
        clase = "todas";
        claseTituloIntegrantes.style.display="block";
        claseTituloDirector.style.display="block";
        break;
      case "lacis":
        clase = "lacis";
        claseTituloIntegrantes.style.display="block";
        claseTituloDirector.style.display="block";
        break;
      case "director":
        clase = "director";
        claseTituloIntegrantes.style.display="none";
        claseTituloDirector.style.display="block";
       
        break;
      case "investigador":
        clase = "investigador";
        claseTituloDirector.style.display="none";
        claseTituloIntegrantes.style.display="block";
       
        break;
      case "estudiante":
        clase = "estudiante";
        claseTituloDirector.style.display="none";
        claseTituloIntegrantes.style.display="block";
       
        break;
      case "asesorExterno":
        clase = "asesor-externo";
        claseTituloDirector.style.display="none";
        claseTituloIntegrantes.style.display="block";
     
        break;
    }

    let offcanvas = document.getElementById('offcanvasWithBothOptions');
    bootstrap.Offcanvas.getInstance(offcanvas)?.hide();

    filtrarCartas(clase);
  });
});

document.addEventListener('DOMContentLoaded', () => {
  if (typeof AOS !== 'undefined') {
    AOS.init();
  }
});
  

  




