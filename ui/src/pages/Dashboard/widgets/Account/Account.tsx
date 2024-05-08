import { Label } from 'components/Label';
import { OutputContainer } from 'components/OutputContainer';
import { FormatAmount } from 'components/sdkDappComponents';
import { useGetAccountInfo, useGetNetworkConfig } from 'hooks';
import { useEffect, useState } from 'react';

export const Account = () => {
  const { network } = useGetNetworkConfig();
  const { address, account } = useGetAccountInfo();
  const [peerStatus, setPeerStatus] = useState<string>('');
  const [peerId, setPeerId] = useState<string>('');

  useEffect(() => {
    const stateInterval = setInterval(() => {
      const peerInfo = wasmGetPeerState()
      setPeerStatus(peerInfo.state)
      setPeerId(peerInfo.id)
    }, 100);

    return () => clearInterval(stateInterval);
  }, [setPeerStatus, peerStatus]);

  return (
    <OutputContainer>
      <div className='flex flex-col text-black' data-testid='topInfo'>
        <p className='truncate'>
          <Label>Address: </Label>
          <span data-testid='accountAddress'> {address}</span>
        </p>
        <p>
          <Label>Peer Status:</Label>
          <span data-testid='peerStatus'> {peerStatus}</span>
        </p>
        <p>
          <Label>Peer ID:</Label>
          <span data-testid='peerId'> {peerId}</span>
        </p>
        <p>
          <Label>Balance: </Label>
          <FormatAmount
            value={account.balance}
            egldLabel={network.egldLabel}
            data-testid='balance'
          />
        </p>
      </div>
    </OutputContainer>
  );
};
