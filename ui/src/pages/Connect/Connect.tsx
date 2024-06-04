import { useEffect, useState } from "react";

type NetworkData = {
  id?: string,
  status?: string,
  version?: string,
  state?: string,
  remotePeer?: string,
}


export const Connect = () => {
  const [networkData, setNetworkData] = useState<NetworkData>({});
  const [swapEvents, setSwapEvents] = useState<string[]>([]);
  useEffect(() => {
    const fetchNetworkData = setInterval(() => {
      const nodeVersion: string = wasmVersion;
      const peerState: PeerInfo = wasmGetPeerState();

      setNetworkData({ ...networkData, version: nodeVersion, id: peerState.id, state: peerState.state, remotePeer: peerState.remotePeer });
    }, 50);

    const fetchSwapEvents = setInterval(() => {
      try{
        const swapEventsFetched: string[] = wasmGetSwapEvents();
        setSwapEvents([...swapEventsFetched]);
      } catch(error) {
        console.error(error);
      }
    }, 10);

    return () => {
      clearInterval(fetchNetworkData)
      clearInterval(fetchSwapEvents)
    };

  }, [networkData, swapEvents]);

  return (
    <div className='h-full text-white font-outfit rounded-md p-10 w-full min-w-[500px] max-w-[600px] flex-col flex gap-10'>
      <div className='w-full bg-red-500 text-red-800 font-bold text-center'>Developer Tool</div>
      <div className='relative h-[40rem] bg-neutral-800 rounded-lg border-[1px] border-neutral-700 p-4 flex flex-col gap-2'>
        <div className='flex place-items-center gap-2'>
          Status: {networkData.state}
          {networkData.state === 'PeerCommunicating' ?
            <div className='h-3 w-3 bg-green-500 rounded-full'></div> :
            <div className='h-3 w-3 bg-yellow-500 rounded-full'></div>}
        </div>
        <div>Version: {networkData.version}</div>
        <div>Local ID: {networkData.id}</div>
        <div>Remote ID: {networkData.remotePeer ? networkData.remotePeer : 'N/A'}</div>

        <ol className="relative text-gray-500 border-s border-gray-200 dark:border-zinc-700 dark:text-gray-400 m-auto">
          {swapEvents.map((event, index) => (
            <li className="mb-10 ms-6">
              <span className="absolute flex items-center justify-center w-8 h-8 bg-green-200 ring-4 ring-neutral-800 rounded-full -start-4 dark:bg-green-900">
                <svg className="w-3.5 h-3.5 text-green-500 dark:text-green-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
                  <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M1 5.917 5.724 10.5 15 1.5" />
                </svg>
              </span>
              <h3 className="font-medium leading-tight">Initialization</h3>
              <p className="text-sm">Step details here</p>
            </li>
          ))}
        </ol>


      </div>
    </div>
  );
}