(function() {
  'use strict';

  const go = new Go();

  WebAssembly
    .instantiateStreaming(fetch("lib.wasm"), go.importObject)
    .then((gotherworld) => {
      go.run(gotherworld.instance);
    });
})();
