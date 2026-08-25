

/* arreglos para colocar los datos de los miembros del grupo */
let director=[];
let todos_los_integrantes=[];
let todos_los_directores=[];

let directorLacis=[];
let todos_los_directoresLacis=[];
let  todos_los_integrantes_lacis=[];



/* contenedores padres para colocar las cartas dentro */
let contenedorPrincipalDeDirector=document.getElementById('contenedor-director');
let contenedorPrincipalDeIntegrantes=document.getElementById('contenedor-cartas-integrantes');
let contenedorPrincipalDeDirectoresDeLinea=document.getElementById('contenedor-directores-deLinea');





 /* funcion para crear cartas de director */

function CrearCartaDirector(datos){

  if(datos.length>0){

    datos.forEach(director=>{

      let contenedor_director=document.createElement('div');                      /*creacion de un contenedor div para colocar las cartas */
      contenedor_director.setAttribute("class", "col-md-12 mx-auto text-center"); /* asignacion de atributos (clases) al contenedor */

      if(director.mail){

            contenedor_director.innerHTML+=`
      
          ${ director.rol==="CO-DIRECTOR" ? ` <div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">CO-DIRECTOR</h2>
            </div> `
            :
             `<div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">DIRECTOR</h2>
            </div>` }
             
            <div class=" card text-center mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 carta-director ${director.clase_rol}"
              data-aos="fade-up" data-aos-duration="500">
              
              <h4 class=" titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${director.rol}</h4>
              <div class="my-2">
                ${!director.imagen ? `<img src="../../static/assets/img-GIS/IMGIntegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="..."> ` 
                  :
                    ` <img src="${director.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">` }
               
              </div>

              <div class="card-body">
                <h5 class="nombre card-title">${director.nombre} </h5>
                <a href="${director.cv}" target="_blank"
                  class="  ms-3 me-3 fs-4"><i class="bi bi-file-earmark-person "></i></a>
                <a class="mail ms-3 me-3 fs-4" onclick="email_despliegue(this)"><i class="bi bi-envelope-fill"></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${director.especializacion}</h6>
                  <p class="card-text">${director.descripcion}</p>
                </div>
                <span class="email-text">${director.mail}</span>
              </div>
            </div>


          


      `

      }else if(director.linkedin){

         contenedor_director.innerHTML+=`
      
       
             <div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">DIRECTOR</h2>
            </div>
            <div class=" card text-center mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 carta-director ${director.clase_rol}"
              data-aos="fade-up" data-aos-duration="500">
              
              <h4 class=" titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${director.rol}</h4>
              <div class="my-2">
                 ${!director.imagen ? `<img src="../../static/assets/img-GIS/IMGIntegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="..."> ` 
                  :
                    ` <img src="${director.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">` }
              </div>

              <div class="card-body">
                <h5 class="nombre card-title">${director.nombre} </h5>
                <a href="${director.cv}" target="_blank"
                  class="linkedin  ms-3 me-3 fs-4"><i class="bi bi-file-earmark-person "></i></a>
                <a href="${director.linkedin}"target="_blank" class="mail ms-3 me-3 fs-4" ><i class="bi bi-linkedin"></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${director.especializacion}</h6>
                  <p class="card-text">${director.descripcion}</p>
                </div>
              </div>
            </div>


          


      `

      }

     
      
      contenedorPrincipalDeDirector.appendChild(contenedor_director); /* Este método inserta el elemento como hijo de un elemento existente padre */
    

    })


    

  }

}

/* funcion para crear cartas de los directores en linea */

function CrearCartaDirectores(datos){

if(datos.length>0){
  

  datos.forEach(directores=>{

    let contenedorCartaDirectoresDeLinea=document.createElement('div'); /*creacion de un contenedor div para colocar las cartas */
    contenedorCartaDirectoresDeLinea.setAttribute("class", "col-sm-7  mb-3 mb-sm-0 text-center");/* asignacion de atributos (clases) al contenedor */
     
    if(directores.grupo==="lacis"){

      
    if(directores.mail){

          contenedorCartaDirectoresDeLinea.innerHTML+=`
      
             <div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">MIEMBRO DE COMITE</h2>
            </div>
            <div class="card mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 h-5 carta-integrantes ${directores.clase_rol}"
              data-aos="fade-up" data-aos-duration="1000">
              <h5 class="titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${directores.rol}</h5>
              <div class="my-2">
                ${!directores.imagen ? `<img src="../../static/assets/img-GIS/IMGintegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                   : 
                   `<img src="${directores.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                }
                
              </div>
              <div class="card-body">
                <h5 class="nombre card-title">${directores.nombre}</h5>
                <a href="${directores.cv}" target="_blank"
                  class="  ms-3 me-3 fs-4"><i class="bi bi-file-earmark-person"></i></a>
                <a class="mail ms-3 me-3 fs-4" onclick="email_despliegue(this)"><i class=" bi bi-envelope-fill" ></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${directores.especializacion}</h6>
                  <p class="card-text">${directores.descripcion}
                  </p>
                </div>
                <span class="email-text">${directores.mail}</span>
              </div>
            </div>


          `
      

    }else if(directores.linkedin){
      contenedorCartaDirectoresDeLinea.innerHTML+=`
      
             <div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">MIEMBRO DE COMITE</h2>
            </div>
            <div class="card mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 h-5 carta-integrantes ${directores.clase_rol}"
              data-aos="fade-up" data-aos-duration="1000">
              <h5 class="titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${directores.rol}</h5>
              <div class="my-2">
                 ${!directores.imagen ? `<img src="../../static/assets/img-GIS/IMGintegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                   : 
                   `<img src="${directores.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                }
              </div>
              <div class="card-body">
                <h5 class="nombre card-title">${directores.nombre}</h5>
                <a href="${directores.cv}" target="_blank"
                  class="  ms-3 me-3 fs-4"><i class="bi bi-file-earmark-person"></i></a>
                <a href="${directores.linkedin}" target="_blank"
                  class="mail ms-3 me-3 fs-4" ><i class="bi bi-linkedin"></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${directores.especializacion}</h6>
                  <p class="card-text">${directores.descripcion}
                  </p>
                </div>
              </div>
            </div>


      `
      

    }

     
    
    
    contenedorPrincipalDeDirectoresDeLinea.appendChild(contenedorCartaDirectoresDeLinea);/* Este método inserta el elemento como hijo de un elemento existente padre */
   
    
     

    }else if(directores.grupo!=="lacis"){
     

    if(directores.mail){

          contenedorCartaDirectoresDeLinea.innerHTML+=`
      
             <div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">DIRECTOR DE LINEA </h2>
            </div>
            <div class="card mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 h-5 carta-integrantes ${directores.clase_rol}"
              data-aos="fade-up" data-aos-duration="1000">
              <h5 class="titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${directores.rol}</h5>
              <div class="my-2">
                ${!directores.imagen ? `<img src="../../static/assets/img-GIS/IMGintegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                   : 
                   `<img src="${directores.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                }
                
              </div>
              <div class="card-body">
                <h5 class="nombre card-title">${directores.nombre}</h5>
                <a href="${directores.cv}" target="_blank"
                  class="  ms-3 me-3 fs-4"><i class="bi bi-file-earmark-person"></i></a>
                <a class="mail ms-3 me-3 fs-4" onclick="email_despliegue(this)"><i class=" bi bi-envelope-fill" ></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${directores.especializacion}</h6>
                  <p class="card-text">${directores.descripcion}
                  </p>
                </div>
                <span class="email-text">${directores.mail}</span>
              </div>
            </div>


          `
      

    }else if(directores.linkedin){
      contenedorCartaDirectoresDeLinea.innerHTML+=`
      
             <div class="container my-5" data-aos="fade-right" data-aos-duration="1000">
              <h2 class="directores text-center">DIRECTOR DE LINEA </h2>
            </div>
            <div class="card mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 h-5 carta-integrantes ${directores.clase_rol}"
              data-aos="fade-up" data-aos-duration="1000">
              <h5 class="titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${directores.rol}</h5>
              <div class="my-2">
                 ${!directores.imagen ? `<img src="../../static/assets/img-GIS/IMGintegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                   : 
                   `<img src="${directores.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
                }
              </div>
              <div class="card-body">
                <h5 class="nombre card-title">${directores.nombre}</h5>
                <a href="${directores.cv}" target="_blank"
                  class="  ms-3 me-3 fs-4"><i class="bi bi-file-earmark-person"></i></a>
                <a href="${directores.linkedin}" target="_blank"
                  class="mail ms-3 me-3 fs-4" ><i class="bi bi-linkedin"></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${directores.especializacion}</h6>
                  <p class="card-text">${directores.descripcion}
                  </p>
                </div>
              </div>
            </div>


      `
      

    }

     
    
    
    contenedorPrincipalDeDirectoresDeLinea.appendChild(contenedorCartaDirectoresDeLinea);/* Este método inserta el elemento como hijo de un elemento existente padre */
   
    

    }
   
  });
 
    
  

}


}
 

/* funcion para crear cartas de los integrantes */

function CrearCartaIntegrantes(datos){

    

    if(datos.length>0){

      datos.forEach(integrante => {
        let contenedorCartaIntegrantes=document.createElement('div');/*creacion de un contenedor div para colocar las cartas (Este método crea un nuevo objeto del tipo especificado por el atributo pasado por parametro)*/
        contenedorCartaIntegrantes.setAttribute("class","col my-4");/* asignacion de atributos (clases) al contenedor */

        if(integrante.linkedin){


          contenedorCartaIntegrantes.innerHTML +=`
        
        <div class="card mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 carta-integrantes ${integrante.clase_rol}"
              data-aos="fade-up" data-aos-duration="10000">
              <h5 class="titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${integrante.rol}
              </h5>
              <div class="my-2">
              ${!integrante.imagen ? `<img src="../../static/assets/img-GIS/IMGintegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="..."></img> `
                :
                ` <img src="${integrante.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
              }
              </div>
              <div class="card-body ">
                <h5 class="nombre card-title">${integrante.nombre}</h5>
                <a href="${integrante.cv}" target="_blank" class="  ms-3 me-3 fs-4"><i
                    class="bi bi-file-earmark-person"></i></a>
                <a href="${integrante.linkedin}" target="_blank"
                  class="mail ms-3 me-3 fs-4" ><i class="bi bi-linkedin"></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${integrante.especializacion} </h6>
                  <p class="card-text">${integrante.descripcion}
                  </p>
                </div>
              </div>
          </div>

        
        `
        

        }else if(integrante.mail){

          contenedorCartaIntegrantes.innerHTML +=`
          <div class="card mx-auto shadow p-3 mb-5 bg-body-tertiary rounded-4 carta-integrantes ${integrante.clase_rol}"
              data-aos="fade-up" data-aos-duration="10000">
              <h5 class="titulo-principal border-primary-subtle card-title border-bottom border-2 py-2">${integrante.rol}
              </h5>
              <div class="my-2">
                 ${!integrante.imagen ? `<img src="../../static/assets/img-GIS/IMGintegrantes/default.png"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="..."></img> `
                :
                ` <img src="${integrante.imagen}"
                  class="card-img-top  rounded-circle border border-white border-5 imagen-integrantes " alt="...">`
              }
              </div>
              <div class="card-body ">
                <h5 class="nombre card-title">${integrante.nombre}</h5>
                <a href="${integrante.cv}" target="_blank" class="  ms-3 me-3 fs-4"><i
                    class="bi bi-file-earmark-person"></i></a>
                <a class="mail ms-3 me-3 fs-4" onclick="email_despliegue(this)"><i class=" bi bi-envelope-fill"></i></a>
                <div class="container">
                  <i class="bi bi-plus-circle icono-mas-info" onclick="despliegue_cartas(this)"></i>
                </div>
                <div class="texto-oculto">
                  <h6 class="subtitulo-cartas-cuerpo card-title my-2">${integrante.especializacion} </h6>
                  <p class="card-text">${integrante.descripcion}
                  </p>
                </div>
                <span class="email-text">${integrante.mail}</span>
              </div>
          </div>

          
        `
        

        }

       
        contenedorPrincipalDeIntegrantes.appendChild(contenedorCartaIntegrantes);/* Este método inserta el elemento como hijo de un elemento existente padre */
        
       
        
        
        
      
  

      });
      
    
    
    }

    

}
 


 /* Codigo para cargar datos desde el JSON a JS */

fetch('../../static/json-GIS/director.json')
  .then(response => response.json())
  .then(data => {
    director=data;
    CrearCartaDirector(director); /* invocacion de funcion a cargar carta director */

  })

  fetch('../../static/json-GIS/directorLacis.json')
  .then(response => response.json())
  .then(data => {
    directorLacis = data; 
  })

   fetch('../../static/json-GIS/directoresLacis.json')
  .then(response => response.json())
  .then(data => {
    todos_los_directoresLacis = data; 
  })


 fetch('../../static/json-GIS/directores.json')
  .then(response => response.json())
  .then(data => {
    todos_los_directores=data;
    CrearCartaDirectores(todos_los_directores); /* invocacion de funcion a cargar carta directores de linea */
    
  })

    

  fetch('../../static/json-GIS/integrantes.json')
  .then(response => response.json())
  .then(data => {
    todos_los_integrantes=data;
    CrearCartaIntegrantes(todos_los_integrantes); /* invocacion de funcion a cargar carta integrantes */
    
  })

    fetch('../../static/json-GIS/integrantesLacis.json')
  .then(response => response.json())
  .then(data => {
    todos_los_integrantes_lacis=data;
    
  })



 

