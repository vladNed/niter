import { useEffect, useState } from "react";
import { Clipboard } from "components/Icons";

type NetworkData = {
  id?: string,
  status?: string,
  version?: string,
  state?: string,
  remotePeer?: string,
}


export const Connect = () => {
  const [msgVisible, setMsgVisible] = useState<boolean>(false);
  const [networkData, setNetworkData] = useState<NetworkData>({})
  const [offerId, setOfferId] = useState<string>('')
  const [exchangeData, setExchangeData] = useState<PeerData[]>([]);
  const [sendData, setSendData] = useState<string>('');

  const handleConnect = async () => {
    await wasmCreateAnswer(offerId);
  }

  const handleSendData = async () => {
    await wasmSendData(sendData);
    setSendData('');
  }

  const handleCreateOffer = async () => {
    const newOffer = await wasmCreateOffer();
    setOfferId(newOffer);
  }

  const handleClipboardCopy = () => {
    navigator.clipboard.writeText(offerId).then(() => {
      setMsgVisible(true);
      setTimeout(() => setMsgVisible(false), 1000);
    });
  }

  useEffect(() => {
    const fetchNetworkData = setInterval(() => {
      const nodeVersion: string = wasmVersion;
      const peerState: PeerInfo = wasmGetPeerState();

      setNetworkData({ ...networkData, version: nodeVersion, id: peerState.id, state: peerState.state, remotePeer: peerState.remotePeer});
    }, 50);

    const fetchExchangeData = setInterval(async () => {
      const pollData: PeerData[] = await wasmPollExchangeData();
      setExchangeData(pollData);
    }, 50);

    return () => {
      clearInterval(fetchNetworkData)
      clearInterval(fetchExchangeData)
    };

  }, [networkData, exchangeData]);

  return (
    <div className='h-full text-white font-outfit rounded-md p-10 w-full min-w-[500px] max-w-[600px] flex-col flex gap-10'>
      <div className='w-full bg-red-500 text-red-800 font-bold text-center'>Developer Tool</div>
      <div className='w-full flex-col flex gap-2'>
        <div className='text-2xl font-bold'>Friendly Connect</div>
        <span className='text-sm text-neutral-500 mb-4'>*Connect directly to a friend through connection string.</span>
        <input
          type='text'
          className='w-full bg-neutral-800 px-2 py-2 outline-none text-md rounded-lg border-[1px] border-neutral-700'
          onChange={(e) => setOfferId(e.target.value)}
        />
        <button
          className='w-full bg-primary-700 px-2 py-2 rounded-lg text-md font-bold hover:bg-primary-600 transition duration-300 active:bg-secondary-600'
          onClick={handleConnect}
        >
          Connect
        </button>
      </div>
      <div className='flex flex-col gap-2'>
        <div className='relative flex flex-row gap-2'>
          <input type='text' value={offerId} disabled={true} className='w-full bg-neutral-800 px-2 py-2 outline-none text-md rounded-lg border-[1px] border-neutral-700' />
          <button
            className='p-2 text-neutral-500 rounded-lg hover:bg-neutral-700 hover:text-neutral-400 transition duration-300 ease-in-out'
            onClick={handleClipboardCopy}
          >
            <Clipboard />
          </button>
          {msgVisible && <div className='absolute end-1 top-10 bg-neutral-800 p-2 rounded-lg text-neutral-400 border-[1px] border-neutral-600'>Copied!!</div>}
        </div>
        <button
          className='w-full bg-primary-700 px-2 py-2 rounded-lg text-md font-bold hover:bg-primary-600 transition duration-300 active:bg-secondary-600'
          onClick={handleCreateOffer}
        >
          Generate offer
        </button>
      </div>
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

        <div className='h-full w-full border-[1px] border-neutral-600 rounded-xl p-4'>
          {exchangeData.map((data, index) => (
            <div key={index} className='flex flex-row gap-2'>
              <div className={`w-1/6` + (data.side === 'local' ? ' text-green-400' : ' text-orange-500')}>{data.side}</div>
              <div className='w-5/6'>{data.data}</div>
            </div>
          ))}
        </div>
        <input
          type='text'
          className='w-full bg-neutral-800 px-2 py-2 outline-none text-md rounded-lg border-[1px] border-neutral-700'
          placeholder="Type message..."
          value={sendData}
          onChange={(e) => setSendData(e.target.value)}
        />
        <button
          className='bg-primary-700 w-1/3 rounded-lg hover:bg-primary-600 transition duration-300 active:bg-secondary-600 font-bold'
          onClick={handleSendData}
        >
          Send
        </button>
      </div>

    </div>
  );
}