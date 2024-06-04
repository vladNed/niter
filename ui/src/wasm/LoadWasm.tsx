import { useEffect, useState } from 'react';
import './wasm_exec.js';
import _ from 'wasm/wasmTypes.js';
import { API_URL, environment, SIGNALLING_SERVER_URL, WASM_LOG_LEVEL } from 'config';
import { LoadingModal } from 'components/index.js';
import { db } from 'db';
import { useGetAccountInfo } from '@multiversx/sdk-dapp/hooks/account/useGetAccountInfo.js';

export const LoadWasm = (props: any) => {
  const [isWasmLoaded, setWasmLoaded] = useState(false);
  const { account } = useGetAccountInfo();

  useEffect(() => {

    const wasm_config = {
      logLevel: WASM_LOG_LEVEL,
      signallingServerUrl: SIGNALLING_SERVER_URL,
      network: environment,
      mvxGatewayUrl: API_URL
    }

    const loadWasm = async () => {
      const go = new window.Go();
      const wasm = await WebAssembly.instantiateStreaming(fetch('wasm/niter.wasm'), go.importObject);
      go.run(wasm.instance);
      const mvxWallet = account.address;
      const mainWallet = await db.getWallet('default');
      if (mainWallet) {
        await wasmInitWallet(mainWallet.wif, mvxWallet)
      } else {
        const wif = await wasmInitWallet(mvxWallet);
        await db.addWallet({ label: 'default', wif: wif });
      }
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

  return isWasmLoaded ? props.children : <LoadingModal />;
}
