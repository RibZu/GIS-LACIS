function email_despliegue(boton) {

  
      
       let cardBody = boton.parentElement; 
        let mostrar=cardBody.querySelector('.email-text');

        if(mostrar.style.display=="block"){
            mostrar.style.display="none";
             


        }else{

            mostrar.style.display="block";
          
        }
       
    
   

  
}

