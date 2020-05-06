(function() {
  'use strict';

  const ASSETS_PATH = 'assets';
  const go = new Go();
  const params = parseParameters();

  const canvas = document.getElementById('gotherworld');
  const ctx = canvas.getContext('2d');

  ctx.blit = function(buffer) {
    //console.log('this is foo', buffer)
    let offset = 0;
    let lastPixel = '';
    buffer.forEach((pixel) => {
      if (pixel !== lastPixel) {
        lastPixel = pixel;
        ctx.fillStyle = '#' + pixel.toString(16).padStart(6, '0') + 'FF';
      }
      const x = offset % 320;
      const y = parseInt(offset / 320, 10);
      ctx.fillRect(x, y, 1, 1)
      offset++;
    });
  }

  loadAllAssets()
    .then(([gotherworld, memList, banks]) => {
      console.info('Assets loaded:', {memList, banks});
      go.run(gotherworld.instance);
      console.info('Initializing with:', params);
      if (isFinite(params.logLevel)) {
        console.info('Updating log level to:', params.logLevel)
        setLogLevel(params.logLevel);
      }

      initializeKeyEventListner();
      initGameFromURI(memList, ...banks);
      startGameFromPart(params.gamePart);
    });

  function loadAllAssets() {
    return Promise.all([
      loadGoWasm('lib.wasm'),
      loadFileAsBytes(`${ASSETS_PATH}/memlist.bin`),
      loadBankAssets()
    ]);
  }

  function loadGoWasm(filename) {
    return fetch(filename)
      .then((wasmLib) => {
        return WebAssembly.instantiateStreaming(wasmLib, go.importObject)
      });
  }

  function loadBankAssets() {
    const filePromises = [
      'bank01',
      'bank02',
      'bank03',
      'bank04',
      'bank05',
      'bank06',
      'bank07',
      'bank08',
      'bank09',
      'bank0a',
      'bank0b',
      'bank0c',
      'bank0d',
    ].map((filename) => {
      return loadFileAsBytes(`${ASSETS_PATH}/${filename}`);
    });
    return Promise.all(filePromises);
  }

  async function loadFileAsBytes(filename) {
    return fetch(filename)
      .then((response) => response.arrayBuffer())
      .then((buffer) => new Uint8Array(buffer));
  }

  function initializeKeyEventListner() {
    document.addEventListener('keydown', forwardKeyEvent);
    document.addEventListener('keyup', forwardKeyEvent);
  }

  function forwardKeyEvent(event) {
    if (event.repeat) {
      return; //Ignore repeat, only up and down is important.
    }
    handleKeyEvent(event.key, event.keyCode, event.type);
  }

  function parseParameters() {
    const rawQuery = location.search.substr(1);
    return rawQuery.split("&")
      .filter((pair) => Boolean(pair))
      .reduce((map, pair) => {
        console.log("parsing pair", pair)
        const [key, value] = pair.split("=");
        const intValue = parseInt(value, 10);
        map[key] = isFinite(intValue) ? intValue :value;
        return map;
      }, {});
  }
})();
