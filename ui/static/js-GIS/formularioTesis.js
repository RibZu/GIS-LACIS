/**
 * formularioTesis.js - Comportamiento interactivo para los formularios de Crear y Editar Tesis.
 * Manejo de alternancia entre integrantes registrados y autores/directores externos,
 * y validación de extensiones de archivos PDF.
 */

document.addEventListener('DOMContentLoaded', function () {
    // --- 1. Control de Autor (Registrado vs Externo) ---
    const radAutorRegistrado = document.getElementById('autor_tipo_registrado');
    const radAutorExterno = document.getElementById('autor_tipo_externo');
    const contenedorAutorRegistrado = document.getElementById('contenedor_autor_registrado');
    const contenedorAutorExterno = document.getElementById('contenedor_autor_externo');
    const selectAutorId = document.getElementById('autor_id');
    const inputAutorHistorico = document.getElementById('autor_historico');

    function actualizarVisibilidadAutor() {
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
        radAutorRegistrado.addEventListener('change', actualizarVisibilidadAutor);
        radAutorExterno.addEventListener('change', actualizarVisibilidadAutor);
        actualizarVisibilidadAutor();
    }

    // --- 2. Control de Director (Registrado vs Externo) ---
    const radDirRegistrado = document.getElementById('director_tipo_registrado');
    const radDirExterno = document.getElementById('director_tipo_externo');
    const contenedorDirRegistrado = document.getElementById('contenedor_director_registrado');
    const contenedorDirExterno = document.getElementById('contenedor_director_externo');
    const selectDirId = document.getElementById('director_id');
    const inputDirHistorico = document.getElementById('director_historico');

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

    // --- 3. Control de Codirector (Ninguno vs Registrado vs Externo) ---
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

    // --- 4. Validación de archivo PDF ---
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
                const maxBytes = 25 * 1024 * 1024; // 25 MB
                if (file.size > maxBytes) {
                    alert('El archivo supera el tamaño máximo permitido de 25 MB.');
                    this.value = '';
                }
            }
        });
    }

    // Sincronizar nivel académico con carreras sugeridas si es necesario
    const selectNivel = document.getElementById('nivel');
    const selectCarrera = document.getElementById('carrera_origen');
    if (selectNivel && selectCarrera) {
        selectNivel.addEventListener('change', function () {
            const nivelVal = this.value;
            if (nivelVal === 'Doctorado' && selectCarrera.value === '') {
                selectCarrera.value = 'Doctorado en Ingeniería en Informática';
            } else if (nivelVal === 'Especialización' && selectCarrera.value === '') {
                selectCarrera.value = 'Especialización en Ingeniería de Software';
            }
        });
    }
});
