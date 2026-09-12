

document.addEventListener('DOMContentLoaded', function () {
  const cards = Array.from(document.querySelectorAll('.tesis-card-item'));
  const inputPalabra = document.getElementById('filtroPalabraClave');
  const selectAnio = document.getElementById('filtroAnio');
  const selectCarrera = document.getElementById('filtroCarrera');
  const tipoBtns = document.querySelectorAll('.filtro-tipo-btn');
  const contadorBadge = document.getElementById('contadorTesisBadge');
  const noResultsMsg = document.getElementById('noResultsTesisMsg');
  const btnLimpiarWrap = document.getElementById('btnLimpiarWrap');
  const btnLimpiar = document.getElementById('btnLimpiarFiltros');
  const btnResetDesdeVacio = document.getElementById('btnResetDesdeVacio');

  if (!cards.length && !contadorBadge) return;

  let activeTipo = 'all'; 
  let activeAnio = 'all';
  let activeCarrera = 'all';
  let activePalabra = '';

 
  if (selectAnio) {
    const aniosSet = new Set();
    cards.forEach(card => {
      const anio = (card.dataset.anio || '').trim();
      if (anio && anio !== '0' && anio !== '-') {
        aniosSet.add(anio);
      }
    });

    const aniosOrdenados = Array.from(aniosSet).sort((a, b) => parseInt(b) - parseInt(a));
    aniosOrdenados.forEach(anio => {
      const opt = document.createElement('option');
      opt.value = anio;
      opt.textContent = anio;
      selectAnio.appendChild(opt);
    });
  }


  if (selectCarrera) {
    const carrerasSet = new Set();
    cards.forEach(card => {
      const carrera = (card.dataset.carrera || '').trim();
      if (carrera) {
        carrerasSet.add(carrera);
      }
    });

    const carrerasOrdenadas = Array.from(carrerasSet).sort((a, b) => a.localeCompare(b));
    carrerasOrdenadas.forEach(carrera => {
      const opt = document.createElement('option');
      opt.value = carrera;
      opt.textContent = carrera;
      selectCarrera.appendChild(opt);
    });
  }


  function aplicarFiltros() {
    let visibles = 0;
    const palabraNorm = activePalabra.toLowerCase().trim();

    cards.forEach(card => {
      const nivel = (card.dataset.nivel || '').toLowerCase();
      const tipo = (card.dataset.tipo || '').toLowerCase();
      const anio = (card.dataset.anio || '').trim();
      const carrera = (card.dataset.carrera || '').trim();
      const titulo = (card.dataset.titulo || '').toLowerCase();
      const autor = (card.dataset.autor || '').toLowerCase();
      const director = (card.dataset.director || '').toLowerCase();
      const palabrasClave = (card.dataset.palabras || '').toLowerCase();
      const resumen = (card.dataset.resumen || '').toLowerCase();

   
      let cumpleTipo = true;
      if (activeTipo === 'posgrado') {
        cumpleTipo = (tipo === 'posgrado') || (nivel !== 'grado' && !nivel.includes('grado'));
      } else if (activeTipo === 'grado') {
        cumpleTipo = (tipo === 'grado') || (nivel === 'grado' || nivel.includes('grado'));
      }

    
      let cumpleAnio = true;
      if (activeAnio !== 'all') {
        cumpleAnio = (anio === activeAnio);
      }

    
      let cumpleCarrera = true;
      if (activeCarrera !== 'all') {
        cumpleCarrera = (carrera.toLowerCase() === activeCarrera.toLowerCase());
      }

   
      let cumplePalabra = true;
      if (palabraNorm) {
        const bolsaTexto = `${titulo} ${palabrasClave} ${autor} ${director} ${carrera} ${resumen}`;
        const terminos = palabraNorm.split(/\s+/);
        cumplePalabra = terminos.every(term => bolsaTexto.includes(term));
      }

      if (cumpleTipo && cumpleAnio && cumpleCarrera && cumplePalabra) {
        card.style.display = '';
        visibles++;
      } else {
        card.style.display = 'none';
      }
    });

   
    if (contadorBadge) {
      if (cards.length === 0) {
        contadorBadge.textContent = '0 tesis disponibles';
      } else if (visibles === cards.length) {
        contadorBadge.textContent = `Mostrando ${visibles} ${visibles === 1 ? 'tesis' : 'tesis'}`;
      } else {
        contadorBadge.textContent = `Mostrando ${visibles} de ${cards.length} tesis`;
      }
    }

  
    if (noResultsMsg) {
      noResultsMsg.style.display = (visibles === 0 && cards.length > 0) ? '' : 'none';
    }

   
    const hayFiltrosActivos = (activeTipo !== 'all' || activeAnio !== 'all' || activeCarrera !== 'all' || activePalabra !== '');
    if (btnLimpiarWrap) {
      btnLimpiarWrap.style.display = hayFiltrosActivos ? '' : 'none';
    }


    if (window.AOS && typeof window.AOS.refresh === 'function') {
      window.AOS.refresh();
    }
  }


  function restablecerFiltros() {
    activeTipo = 'all';
    activeAnio = 'all';
    activeCarrera = 'all';
    activePalabra = '';

    if (inputPalabra) inputPalabra.value = '';
    if (selectAnio) selectAnio.value = 'all';
    if (selectCarrera) selectCarrera.value = 'all';

    tipoBtns.forEach(btn => {
      btn.classList.toggle('active', btn.dataset.tipo === 'all');
    });

    aplicarFiltros();
  }


  tipoBtns.forEach(btn => {
    btn.addEventListener('click', function () {
      tipoBtns.forEach(b => b.classList.remove('active'));
      this.classList.add('active');
      activeTipo = this.dataset.tipo;
      aplicarFiltros();
    });
  });

  if (inputPalabra) {
    inputPalabra.addEventListener('input', function () {
      activePalabra = this.value;
      aplicarFiltros();
    });
  }

  if (selectAnio) {
    selectAnio.addEventListener('change', function () {
      activeAnio = this.value;
      aplicarFiltros();
    });
  }

  if (selectCarrera) {
    selectCarrera.addEventListener('change', function () {
      activeCarrera = this.value;
      aplicarFiltros();
    });
  }

  if (btnLimpiar) {
    btnLimpiar.addEventListener('click', restablecerFiltros);
  }

  if (btnResetDesdeVacio) {
    btnResetDesdeVacio.addEventListener('click', restablecerFiltros);
  }

  
  const urlParams = new URLSearchParams(window.location.search);
  const paramTipo = urlParams.get('tipo');
  const paramAnio = urlParams.get('anio');
  const paramCarrera = urlParams.get('carrera');
  const paramPalabra = urlParams.get('palabra');

  if (paramTipo && ['posgrado', 'grado', 'all'].includes(paramTipo.toLowerCase())) {
    activeTipo = paramTipo.toLowerCase();
    tipoBtns.forEach(b => b.classList.toggle('active', b.dataset.tipo === activeTipo));
  }
  if (paramAnio && selectAnio) {
    activeAnio = paramAnio;
    selectAnio.value = paramAnio;
  }
  if (paramCarrera && selectCarrera) {
    activeCarrera = paramCarrera;
    selectCarrera.value = paramCarrera;
  }
  if (paramPalabra && inputPalabra) {
    activePalabra = paramPalabra;
    inputPalabra.value = paramPalabra;
  }

  
  aplicarFiltros();
});
