import { SwapEvents } from 'localConstants';


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
  }

  return {
    getSwapEvents,
  }
};