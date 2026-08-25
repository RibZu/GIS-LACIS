
function despliegue_cartas(icono) {
  
  let carta=icono.closest(".carta-director, .carta-integrantes");

  if(carta.classList.contains("expandir")){
    carta.classList.remove("expandir");
    icono.classList.remove("bi-dash-circle");
    icono.classList.add("bi-plus-circle");
  }else{
    carta.classList.add("expandir");
    icono.classList.remove("bi-plus-circle");
    icono.classList.add("bi-dash-circle");
  }

}

