import { useGetAccountInfo } from "@multiversx/sdk-dapp/hooks/account/useGetAccountInfo";
import { ChevronIcon } from "components/Icons";
import { useEffect, useState } from "react";
import { DrawerProps } from "types";

export const DrawerField = ({ label, value }: { label: string, value: string}) => {
  return (
    <div className='flex justify-between font-roboto items-center py-1 grid grid-cols-12 gap-2'>
      <span className='text-slate-600 text-base col-span-4'>{label}:</span>
      <span className='text-base col-span-8 overflow-hidden truncate'>{value}</span>
    </div>
  )
}

type NetworkData = {
  id?: string,
  status?: string,
  version?: string,
  state?: string,
}

export const NetworkDrawer = (props: DrawerProps) => {
  const [networkData, setNetworkData] = useState<NetworkData>({})
  const { address, account } = useGetAccountInfo();
  const networkFields = [
    { label: 'ID', value: networkData.id || 'N/A'},
    { label: 'Status', value: 'active' },
    { label: 'Version', value: networkData.version || 'N/A'},
    { label: 'State', value: networkData.state || 'N/A' },
  ]
  const accountFields = [
    { label: 'Address', value: address },
    { label: 'Balance', value: `${account.balance} EGLD` },
  ]

  useEffect(() => {
    const fetchNetworkData = setInterval(() => {
      try {
        const nodeVersion: string = wasmVersion;
        const peerState: PeerInfo = wasmGetPeerState();
        setNetworkData({ ...networkData, version: nodeVersion, id: peerState.id, state: peerState.state });
      } catch (error) {
        console.error(error);
      }
    }, 50);

    return () => clearInterval(fetchNetworkData);
  }, []);

  return (
    <div className={`text-black fixed h-screen w-[25rem] bg-slate-300/30 backdrop-blur-md rounded-l-3xl top-0 right-0 flex gap-2 transition-transform duration-300 ` + (props.isOpen ? 'translate-x-0' : 'translate-x-full') }>
      <div className=' hover:text-primary-600 p-5 hover:bg-slate-400/30 rounded-l-3xl transition duration-300 ease-in-out' onClick={props.toggleDrawer}><ChevronIcon /></div>
      <div className='h-full w-full p-4'>
        <div className='flex flex-col border-b-[1px] border-neutral-600 pb-4'>
          <span className='text-2xl font-bold'>Network</span>
          <span className='text-slate-600'>Information about the Niter network node</span>
        </div>
        <div className='flex flex-col py-4 '>
          {networkFields.map((field, index) => (
            <DrawerField key={index} label={field.label} value={field.value} />
          ))}
        </div>
        <div className='flex flex-col border-b-[1px] border-neutral-600 pb-4'>
          <span className='text-2xl font-bold '>Account</span>
          <span className='text-slate-600'>Information about the MVX account</span>
        </div>
        <div className='flex flex-col py-4 '>
          {accountFields.map((field, index) => (
            <DrawerField key={index} label={field.label} value={field.value} />
          ))}
        </div>
      </div>
    </div>
  );
}