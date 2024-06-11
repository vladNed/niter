import { SwapEvents, TransactionRequestTypes } from 'localConstants';


// Hook used to interact with the core wasm module
export const useWasm = () => {

  const getSwapEvents = async (): Promise<SwapEvents[]> => {
    let fetchedSwapEvents: string[];
    try{
      fetchedSwapEvents = await wasmGetSwapEvents();
    } catch (e) {
      return [];
    };

    let data: SwapEvents[] = [];
    for (const event of fetchedSwapEvents) {
      data.push(SwapEvents[event as keyof typeof SwapEvents]);
    }

    return data;
  };

  const resetPeer = async (): Promise<void> => {
    try {
      wasmResetPeer();
    } catch (e) {
      throw new Error('Failed to reset peer:' + e);
    }
  }

  const emitSwapEvent = async (event: SwapEvents, data: object): Promise<void> => {
    const encodedData = Buffer.from(JSON.stringify(data)).toString('base64');
    try {
      await wasmEmitSwapEvent(event.toString(), encodedData);
    } catch (e) {
      throw new Error('Failed to emit swap event:' + e);
    }
  };

  const getTransactionRequest = async (type: TransactionRequestTypes): Promise<any> => {
    try {
      return await wasmTransactionRequest(type.toString());
    } catch (e) {
      throw new Error('Failed to get transaction request:' + e);
    }
  }

  return {
    getSwapEvents,
    emitSwapEvent,
    resetPeer,
    getTransactionRequest
  };
};