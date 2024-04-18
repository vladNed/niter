import { useEffect, useState } from 'react';
import './wasm_exec.js';
import './wasmTypes.d.ts';
import { SIGNALLING_SERVER_URL, WASM_LOG_LEVEL } from 'config';


export const LoadWasm = (props: any) => {
  const [isWasmLoaded, setWasmLoaded] = useState(false);

  useEffect(() => {

    const wasm_config = {
      logLevel: WASM_LOG_LEVEL,
      signallingServerUrl: SIGNALLING_SERVER_URL
    }

    const loadWasm = async () => {
      const go = new window.Go();
      const wasm = await WebAssembly.instantiateStreaming(fetch('wasm/niter.wasm'), go.importObject);
      go.run(wasm.instance);

      try {
        await wasmInit(JSON.stringify(wasm_config))
      } catch (e) {
        console.error('Error initializing wasm', e);
        return;
      }

      setWasmLoaded(true);
    };

    loadWasm();
  }, [setWasmLoaded]);

  return isWasmLoaded ? props.children : null;
}

