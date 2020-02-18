(function() {
  'use strict';

  const ASSETS_PATH = 'assets';
  const go = new Go();

  console.info('Loading assets...');
  const assetsPromise = Promise.all([
    loadGoWasm('lib.wasm'),
    loadFileAsBytes(`${ASSETS_PATH}/memlist.bin`),
    loadBankAssets()
  ]);

  assetsPromise
    .then(([gotherworld, memList, banks]) => {
      console.info('Assets loaded:', {memList, banks});
      go.run(gotherworld.instance);
      initGameFromURI(memList, ...banks);
      // TODO read offset from get parameters
      startGameFromPart();
    });

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
})();
