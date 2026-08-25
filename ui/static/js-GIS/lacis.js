let boton_frontal=document.querySelectorAll('.flip-btn');

boton_frontal.forEach(boton =>{
    boton.addEventListener('click',() =>{

        let cambiar_perspectiva=boton.closest('.carta-flipped');
        if(cambiar_perspectiva.classList.contains('flipped')){
            cambiar_perspectiva.classList.remove('flipped');
        }else{
            cambiar_perspectiva.classList.add("flipped");
        }

    })
})


