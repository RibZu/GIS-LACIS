document.addEventListener('DOMContentLoaded', function () {
    const MAX_AUTORES_GRADO = 4;

    const selectNivel = document.getElementById('nivel');
    const selectCarrera = document.getElementById('carrera_origen');
    const btnAgregarAutor = document.getElementById('btn_agregar_autor');
    const contenedorAutores = document.getElementById('contenedor_autores');
    const templateAutor = document.getElementById('template_autor_item');
    const badgeLimiteAutores = document.getElementById('badge_limite_autores');
    const contadorAutoresSpan = document.getElementById('contador_autores');

    const radAutorRegistrado = document.getElementById('autor_tipo_registrado');
    const radAutorExterno = document.getElementById('autor_tipo_externo');
    const contenedorAutorRegistrado = document.getElementById('contenedor_autor_registrado');
    const contenedorAutorExterno = document.getElementById('contenedor_autor_externo');
    const selectAutorId = document.getElementById('autor_id');
    const inputAutorHistorico = document.getElementById('autor_historico');

    function actualizarVisibilidadAutorPrincipal() {
        if (!radAutorRegistrado || !radAutorExterno) return;
        if (radAutorRegistrado.checked) {
            if (contenedorAutorRegistrado) contenedorAutorRegistrado.style.display = 'block';
            if (contenedorAutorExterno) contenedorAutorExterno.style.display = 'none';
            if (selectAutorId) selectAutorId.required = true;
            if (inputAutorHistorico) inputAutorHistorico.required = false;
        } else {
            if (contenedorAutorRegistrado) contenedorAutorRegistrado.style.display = 'none';
            if (contenedorAutorExterno) contenedorAutorExterno.style.display = 'block';
            if (selectAutorId) selectAutorId.required = false;
            if (inputAutorHistorico) inputAutorHistorico.required = true;
        }
    }

    if (radAutorRegistrado && radAutorExterno) {
        radAutorRegistrado.addEventListener('change', actualizarVisibilidadAutorPrincipal);
        radAutorExterno.addEventListener('change', actualizarVisibilidadAutorPrincipal);
        actualizarVisibilidadAutorPrincipal();
    }

    function setupCardAutorAdicional(card, numero) {
        const radReg = card.querySelector('.radio-autor-registrado');
        const radExt = card.querySelector('.radio-autor-externo');
        const lblReg = card.querySelector('.label-autor-registrado');
        const lblExt = card.querySelector('.label-autor-externo');
        const boxReg = card.querySelector('.box-autor-registrado');
        const boxExt = card.querySelector('.box-autor-externo');
        const selId = card.querySelector('.select-autor-id');
        const inHist = card.querySelector('.input-autor-externo');
        const numLabel = card.querySelector('.autor-numero-label');
        const btnQuitar = card.querySelector('.btn-quitar-autor');

        if (numLabel) {
            numLabel.textContent = 'Autor(a) ' + numero;
        }

        const radioGroupName = 'autor_tipo_' + numero;
        const regId = 'autor_tipo_registrado_' + numero;
        const extId = 'autor_tipo_externo_' + numero;

        if (radReg) {
            radReg.name = radioGroupName;
            radReg.id = regId;
        }
        if (lblReg) lblReg.setAttribute('for', regId);

        if (radExt) {
            radExt.name = radioGroupName;
            radExt.id = extId;
        }
        if (lblExt) lblExt.setAttribute('for', extId);

        if (selId) {
            selId.name = 'autor_id_' + numero;
            selId.id = 'autor_id_' + numero;
        }
        if (inHist) {
            inHist.name = 'autor_historico_' + numero;
            inHist.id = 'autor_historico_' + numero;
        }

        function toggleVisibilidad() {
            if (radReg && radReg.checked) {
                if (boxReg) boxReg.style.display = 'block';
                if (boxExt) boxExt.style.display = 'none';
            } else {
                if (boxReg) boxReg.style.display = 'none';
                if (boxExt) boxExt.style.display = 'block';
            }
        }

        if (radReg) radReg.addEventListener('change', toggleVisibilidad);
        if (radExt) radExt.addEventListener('change', toggleVisibilidad);
        toggleVisibilidad();

        if (btnQuitar) {
            btnQuitar.onclick = function () {
                card.remove();
                renumerarAutores();
                actualizarEstadoAutores();
            };
        }
    }

    function renumerarAutores() {
        if (!contenedorAutores) return;
        const cards = contenedorAutores.querySelectorAll('.autor-item');
        cards.forEach(function (card, index) {
            const numero = index + 1;
            if (numero > 1) {
                setupCardAutorAdicional(card, numero);
            }
        });
    }

    function actualizarEstadoAutores() {
        if (!selectNivel || !contenedorAutores) return;
        const val = selectNivel.value || '';
        const esGrado = val === 'Grado' || val.toLowerCase().includes('grado');
        const cards = contenedorAutores.querySelectorAll('.autor-item');
        const total = cards.length;

        if (esGrado) {
            if (btnAgregarAutor) btnAgregarAutor.classList.remove('d-none');
            if (badgeLimiteAutores) {
                badgeLimiteAutores.textContent = 'Grado: hasta 4 autores permitidos';
                badgeLimiteAutores.className = 'badge badge-autores-info badge-grado ms-2';
            }
            if (btnAgregarAutor) {
                if (total >= MAX_AUTORES_GRADO) {
                    btnAgregarAutor.disabled = true;
                    btnAgregarAutor.classList.add('disabled');
                    btnAgregarAutor.title = 'Máximo de 4 autores alcanzado';
                } else {
                    btnAgregarAutor.disabled = false;
                    btnAgregarAutor.classList.remove('disabled');
                    btnAgregarAutor.title = '';
                }
            }
        } else {

            if (btnAgregarAutor) btnAgregarAutor.classList.add('d-none');
            if (badgeLimiteAutores) {
                badgeLimiteAutores.textContent = 'Posgrado: 1 solo autor';
                badgeLimiteAutores.className = 'badge badge-autores-info badge-posgrado ms-2';
            }

            cards.forEach(function (card, index) {
                if (index > 0) {
                    card.remove();
                }
            });
        }

        const countActual = contenedorAutores.querySelectorAll('.autor-item').length;
        if (contadorAutoresSpan) {
            contadorAutoresSpan.textContent = countActual;
        }
    }

    function agregarAutor() {
        if (!selectNivel) return;
        const val = selectNivel.value || '';
        const esGrado = val === 'Grado' || val.toLowerCase().includes('grado');
        if (!esGrado) return;
        if (!contenedorAutores || !templateAutor) return;

        const currentCount = contenedorAutores.querySelectorAll('.autor-item').length;
        if (currentCount >= MAX_AUTORES_GRADO) return;

        const clone = templateAutor.content.cloneNode(true);
        const card = clone.querySelector('.autor-item');
        contenedorAutores.appendChild(card);

        renumerarAutores();
        actualizarEstadoAutores();
    }

    if (btnAgregarAutor) {
        btnAgregarAutor.addEventListener('click', agregarAutor);
    }

    function onNivelChange() {
        actualizarEstadoAutores();

        if (selectCarrera && selectCarrera.value.trim() === '') {
            const nivelVal = selectNivel ? selectNivel.value : '';
            if (nivelVal === 'Doctorado') {
                selectCarrera.value = 'Doctorado en Ingeniería en Informática';
            } else if (nivelVal === 'Especialización') {
                selectCarrera.value = 'Especialización en Ingeniería de Software';
            } else if (nivelVal === 'Grado' || nivelVal.toLowerCase().includes('grado')) {
                selectCarrera.value = 'Ingeniería en Informática';
            }
        }
    }

    if (selectNivel) {
        selectNivel.addEventListener('change', onNivelChange);
        selectNivel.addEventListener('input', onNivelChange);
        renumerarAutores();
        actualizarEstadoAutores();
    }

    const radDirRegistrado = document.getElementById('director_tipo_registrado');
    const radDirExterno = document.getElementById('director_tipo_externo');
    const contenedorDirRegistrado = document.getElementById('contenedor_director_registrado');
    const contenedorDirExterno = document.getElementById('contenedor_director_externo');

    function actualizarVisibilidadDirector() {
        if (!radDirRegistrado || !radDirExterno) return;
        if (radDirRegistrado.checked) {
            if (contenedorDirRegistrado) contenedorDirRegistrado.style.display = 'block';
            if (contenedorDirExterno) contenedorDirExterno.style.display = 'none';
        } else {
            if (contenedorDirRegistrado) contenedorDirRegistrado.style.display = 'none';
            if (contenedorDirExterno) contenedorDirExterno.style.display = 'block';
        }
    }

    if (radDirRegistrado && radDirExterno) {
        radDirRegistrado.addEventListener('change', actualizarVisibilidadDirector);
        radDirExterno.addEventListener('change', actualizarVisibilidadDirector);
        actualizarVisibilidadDirector();
    }

    const radCoodirNinguno = document.getElementById('coodirector_tipo_ninguno');
    const radCoodirRegistrado = document.getElementById('coodirector_tipo_registrado');
    const radCoodirExterno = document.getElementById('coodirector_tipo_externo');
    const contenedorCoodirRegistrado = document.getElementById('contenedor_coodirector_registrado');
    const contenedorCoodirExterno = document.getElementById('contenedor_coodirector_externo');

    function actualizarVisibilidadCoodirector() {
        if (!radCoodirNinguno || !radCoodirRegistrado || !radCoodirExterno) return;
        if (radCoodirNinguno.checked) {
            if (contenedorCoodirRegistrado) contenedorCoodirRegistrado.style.display = 'none';
            if (contenedorCoodirExterno) contenedorCoodirExterno.style.display = 'none';
        } else if (radCoodirRegistrado.checked) {
            if (contenedorCoodirRegistrado) contenedorCoodirRegistrado.style.display = 'block';
            if (contenedorCoodirExterno) contenedorCoodirExterno.style.display = 'none';
        } else {
            if (contenedorCoodirRegistrado) contenedorCoodirRegistrado.style.display = 'none';
            if (contenedorCoodirExterno) contenedorCoodirExterno.style.display = 'block';
        }
    }

    if (radCoodirNinguno && radCoodirRegistrado && radCoodirExterno) {
        radCoodirNinguno.addEventListener('change', actualizarVisibilidadCoodirector);
        radCoodirRegistrado.addEventListener('change', actualizarVisibilidadCoodirector);
        radCoodirExterno.addEventListener('change', actualizarVisibilidadCoodirector);
        actualizarVisibilidadCoodirector();
    }

    const fileInput = document.getElementById('archivo_pdf');
    if (fileInput) {
        fileInput.addEventListener('change', function () {
            if (this.files && this.files[0]) {
                const file = this.files[0];
                const filename = file.name.toLowerCase();
                if (!filename.endsWith('.pdf')) {
                    alert('Por favor selecciona un archivo en formato PDF válido.');
                    this.value = '';
                    return;
                }
                const maxBytes = 25 * 1024 * 1024;
                if (file.size > maxBytes) {
                    alert('El archivo supera el tamaño máximo permitido de 25 MB.');
                    this.value = '';
                }
            }
        });
    }

    const form = document.querySelector('form');
    if (form) {
        form.addEventListener('submit', () => {
            const btn = form.querySelector('button[type="submit"]');
            if (btn) {
                btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>Guardando...';
                btn.style.pointerEvents = 'none';
            }
        });
    }
});

