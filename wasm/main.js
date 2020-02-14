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
    const response = await fetch(filename);
    const chunks = [];
    let responseSize = 0;
    for await (const chunk of streamAsyncIterator(response.body)) {
      chunks.push(chunk);
      responseSize += chunk.length;
    }

    if (chunks.length === 1) {
      return chunks[0];
    }

    let index = 0;
    return chunks.reduce((bytes, chunk) => {
      bytes.set(chunk, index);
      index += chunk;
      return bytes;
    }, new Uint8Array(responseSize));
  }

  async function* streamAsyncIterator(stream) {
    const reader = stream.getReader();
    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) {
          return;
        }
        yield value;
      }
    } finally {
      reader.releaseLock();
    }
  }
})();
