import { useEffect, useState, createContext } from 'react';
import './wasm_exec.js';
import _ from 'wasm/wasmTypes.js';
import { environment, SIGNALLING_SERVER_URL, WASM_LOG_LEVEL } from 'config';
import { LoadingModal } from 'components/index.js';
import { db } from 'db';

export const LoadWasm = (props: any) => {
  const [isWasmLoaded, setWasmLoaded] = useState(false);

  useEffect(() => {

    const wasm_config = {
      logLevel: WASM_LOG_LEVEL,
      signallingServerUrl: SIGNALLING_SERVER_URL,
      network: environment
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

      const mainWallet = await db.getWallet('default');
      if (mainWallet) {
        await wasmInitWallet(mainWallet.wif)
      } else {
        const wif = await wasmInitWallet();
        await db.addWallet({ label: 'default', wif: wif });
      }

      setWasmLoaded(true);
    };

    loadWasm();
  }, [setWasmLoaded]);

  return isWasmLoaded ? props.children : <LoadingModal />;
}
