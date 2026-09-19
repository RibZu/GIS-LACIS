let claseTituloIntegrantes = document.querySelector(".titulo-integrantes");
let claseTituloDirector = document.querySelector(".titulo-director");

function AjustarColumnasDirectoresDeLinea(cantidad) {
  const contenedor = document.getElementById("contenedor-directores-deLinea");
  if (!contenedor) return;
  contenedor.classList.remove("row-cols-lg-2", "row-cols-lg-3");
  contenedor.classList.add(cantidad >= 3 ? "row-cols-lg-3" : "row-cols-lg-2");
}

let grupoActual = "todas";

function filtrarCartas(clase) {
  if (clase === "todas" || clase === "lacis") {
    grupoActual = clase;
  }

  if (typeof contenedorPrincipalDeDirector !== "undefined" && contenedorPrincipalDeDirector) {
    contenedorPrincipalDeDirector.innerHTML = "";
  }
  if (typeof contenedorPrincipalDeDirectoresDeLinea !== "undefined" && contenedorPrincipalDeDirectoresDeLinea) {
    contenedorPrincipalDeDirectoresDeLinea.innerHTML = "";
  }
  if (typeof contenedorPrincipalDeIntegrantes !== "undefined" && contenedorPrincipalDeIntegrantes) {
    contenedorPrincipalDeIntegrantes.innerHTML = "";
  }

  let filtradosIntegrantes = [];
  let filtradosDirectores = [];
  let filtradosDirector = [];

  if (clase === "todas") {
    filtradosIntegrantes = todos_los_integrantes;
    filtradosDirectores = todos_los_directores;
    filtradosDirector = director;
  } else if (clase === "lacis") {
    filtradosDirector = directorLacis;
    filtradosDirectores = todos_los_directoresLacis;
    filtradosIntegrantes = todos_los_integrantes_lacis;
  } else if (clase === "director") {
    if (grupoActual === "lacis") {
      filtradosDirector = directorLacis;
      filtradosDirectores = todos_los_directoresLacis;
    } else {
      filtradosDirector = director;
      filtradosDirectores = todos_los_directores;
    }
  } else {
    let baseIntegrantes = (grupoActual === "lacis") ? todos_los_integrantes_lacis : todos_los_integrantes;
    filtradosIntegrantes = baseIntegrantes.filter(p => p.clase_rol === clase);

    if (clase === "estudiante") {
      let extra = baseIntegrantes.filter(p => p.clase_rol === "becario");
      extra.forEach(ex => {
        if (!filtradosIntegrantes.some(f => f.id === ex.id || f.nombre === ex.nombre)) {
          filtradosIntegrantes.push(ex);
        }
      });
    }

    let baseDirectores = (grupoActual === "lacis") ? todos_los_directoresLacis : todos_los_directores;
    filtradosDirectores = baseDirectores.filter(p => p.clase_rol === clase);
    let baseDirector = (grupoActual === "lacis") ? directorLacis : director;
    filtradosDirector = baseDirector.filter(p => p.clase_rol === clase);
  }

  CrearCartaIntegrantes(filtradosIntegrantes);
  CrearCartaDirectores(filtradosDirectores);
  CrearCartaDirector(filtradosDirector);

  AjustarColumnasDirectoresDeLinea(filtradosDirectores.length);

  if (typeof AOS !== 'undefined') {
    AOS.refreshHard();
  }
}

let filtroBoton = document.querySelectorAll(".filtro-btn");

filtroBoton.forEach(boton => {
  boton.addEventListener("click", () => {
    let id = boton.id;
    let clase = "";

    filtroBoton.forEach(b => b.classList.remove("active"));
    boton.classList.add("active");

    switch (id) {
      case "todos":
        clase = "todas";
        if (claseTituloIntegrantes) claseTituloIntegrantes.style.display = "block";
        if (claseTituloDirector) claseTituloDirector.style.display = "block";
        break;
      case "lacis":
        clase = "lacis";
        if (claseTituloIntegrantes) claseTituloIntegrantes.style.display = "block";
        if (claseTituloDirector) claseTituloDirector.style.display = "block";
        break;
      case "director":
        clase = "director";
        if (claseTituloIntegrantes) claseTituloIntegrantes.style.display = "none";
        if (claseTituloDirector) claseTituloDirector.style.display = "block";
        break;
      case "investigador":
        clase = "investigador";
        if (claseTituloDirector) claseTituloDirector.style.display = "none";
        if (claseTituloIntegrantes) claseTituloIntegrantes.style.display = "block";
        break;
      case "estudiante":
        clase = "estudiante";
        if (claseTituloDirector) claseTituloDirector.style.display = "none";
        if (claseTituloIntegrantes) claseTituloIntegrantes.style.display = "block";
        break;
      case "asesorExterno":
        clase = "asesor-externo";
        if (claseTituloDirector) claseTituloDirector.style.display = "none";
        if (claseTituloIntegrantes) claseTituloIntegrantes.style.display = "block";
        break;
    }

    let offcanvas = document.getElementById('offcanvasWithBothOptions');
    if (offcanvas && typeof bootstrap !== 'undefined') {
      bootstrap.Offcanvas.getInstance(offcanvas)?.hide();
    }

    filtrarCartas(clase);
  });
});

document.addEventListener('DOMContentLoaded', () => {
  if (typeof AOS !== 'undefined') {
    AOS.init();
  }
});
