(function() {
  'use strict';

  const go = new Go();

  // TODO consider loading all assets (wasm, memlist, banks)
  // and only then initializing the webassembly code.

  const memlist = loadFileAsBytes("assets/memlist.bin");
  memlist
    .then((bytes) => {
      console.log('Success:', bytes);
    })
    .catch((error) => {
      console.error('Failed to load assets:', error);
    });

  WebAssembly
    .instantiateStreaming(fetch("lib.wasm"), go.importObject)
    .then((gotherworld) => {
      go.run(gotherworld.instance);
      // TODO load memlist and use CopyBytesToGo to pass it in
      // TODO then load all bank assets and pass them in to build bankFilesMap
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
