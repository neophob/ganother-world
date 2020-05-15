(function() {
  'use strict';

  const ASSETS_PATH = 'assets';
  const go = new Go();
  const params = parseParameters();

  const SHIFT_TO_RED_BITS = 16;
  const SHIFT_TO_GREEN_BITS = 8;
  const FULL_ALPHA = 0xFF;

  const canvas = document.getElementById('gotherworld');
  const touchControls = document.getElementById('touch-controls');
  const ctx = canvas.getContext('2d');
  const tempImage = ctx.createImageData(320, 200);

  let showKeyboard = false;

  ctx.blit = function(buffer) {
    const pixel = tempImage.data;
    let ofs = 0;
    buffer.forEach((p) => {
      pixel[ofs++] = (p >> SHIFT_TO_RED_BITS) & 0xFF;
      pixel[ofs++] = (p >> SHIFT_TO_GREEN_BITS) & 0xFF;
      pixel[ofs++] = p & 0xFF;
      pixel[ofs++] = FULL_ALPHA;
    });
    ctx.putImageData(tempImage, 0, 0);
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
      initializeTouchEventListners();
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

  function initializeTouchEventListners() {
    const keyMappings = [
      { id: 'key-up', key: 'ArrowUp', keyCode: 38 },
      { id: 'key-left', key: 'ArrowLeft', keyCode: 37 },
      { id: 'key-down', key: 'ArrowDown', keyCode: 40 },
      { id: 'key-right', key: 'ArrowRight', keyCode: 39 },
      { id: 'key-sp', key: ' ', keyCode: 32 },
      { id: 'key-esc', key: 'Escape', keyCode: 27 },
    ]
    keyMappings.forEach(({id, key, keyCode}) => {
      const keyButton = document.getElementById(id);
      keyButton.addEventListener('mousedown', (e) => {
        handleKeyEvent(key, keyCode, 'keydown');
      });
      // Note we have to handle both up and leave to release button
      keyButton.addEventListener('mouseup', (e) => {
        handleKeyEvent(key, keyCode, 'keyup');
      });
      keyButton.addEventListener('mouseleave', (e) => {
        handleKeyEvent(key, keyCode, 'keyup');
      });
    });

    const toggleButton = document.getElementById('toggle-keyboard');
    toggleButton.addEventListener('click', (e) => {
      toggleKeyboard(toggleButton);
    });
  }

  function toggleKeyboard(button) {
    showKeyboard = !showKeyboard;
    const changeToLabel = showKeyboard ? '⌨︎ ON' : '⌨︎ OFF';
    console.log('toggle keyboard', changeToLabel);
    button.text = changeToLabel;
    document.getElementById('touch-controls').className = showKeyboard ?
      '' :
      'hide';
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
