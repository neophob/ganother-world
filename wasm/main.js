(function() {
  'use strict';

  const go = new Go();

  console.info('Loading assets...');
  const assetsPromise = Promise.all([
    fetch('lib.wasm'),
    loadFileAsBytes('assets/memlist.bin'),
    // TODO load all bank assets
  ]);

  assetsPromise
    .then(([wasmLib, memList, ...banks]) => {
      console.info('Assets loaded:', {memList});

      return WebAssembly.instantiateStreaming(wasmLib, go.importObject)
    })
    .then((gotherworld) => {
      console.info('Running gotherworld...');
      go.run(gotherworld.instance);
      // TODO copy memlist to go
      // TODO copy banks to go
    });

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
    return chucks.reduce((bytes, chunk) => {
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
