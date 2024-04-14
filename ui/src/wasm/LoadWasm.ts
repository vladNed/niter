import { useEffect, useState } from 'react';
import './wasm_exec.js';
import './wasmTypes.d.ts';


export const LoadWasm = (props: any) => {
  const [isWasmLoaded, setWasmLoaded] = useState(false);

  useEffect(() => {
    const loadWasm = async () => {
      const goWasm = new window.Go();
      const res = await WebAssembly.instantiateStreaming(fetch('wasm/niter.wasm'), goWasm.importObject);
      goWasm.run(res.instance);
      wasmGenerateWallet();
      setWasmLoaded(true);
    };

    loadWasm();
  }, [setWasmLoaded]);

  return isWasmLoaded ? props.children : null;
}

