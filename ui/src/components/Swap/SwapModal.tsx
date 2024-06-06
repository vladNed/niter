import { CloseIcon } from 'components/Icons'
import { SwapSide } from 'localConstants';
import { useEffect, useState } from 'react'
import { InitiatorStepOne, OfferConnecting, OfferCreatedStep, ParticipantStepOne } from './Steps';
import { type OfferDetails } from 'types';
import { InitiatorStepTwo } from './Steps/InitiatorStepTwo';


type SwapModalProps = {
  onClose: () => void;
  offerId: string;
  offerData?: OfferDetails;
  swapSide?: SwapSide;
};

export const SwapModal = (props: SwapModalProps) => {
  const [peerState, setPeerState] = useState<string>('');
  const [swapStarted, setSwapStarted] = useState<boolean>(false);
  const [currentSwapStep, setCurrentSwapStep] = useState<number>(0);

  const handleProgress = () => {
    setCurrentSwapStep(currentSwapStep + 1);
  }

  const InitiatorStates = [
    <InitiatorStepOne offerData={props.offerData} progressHandler={handleProgress}/>,
    <InitiatorStepTwo offerData={props.offerData} progressHandler={handleProgress}/>
  ];

  const ReceiverStates = [
    <ParticipantStepOne />
  ];

  const getCurrentStep = () => {
    switch (props.swapSide) {
      case SwapSide.INITIATOR:
        return InitiatorStates[currentSwapStep];
      case SwapSide.PARTICIPANT:
        return ReceiverStates[currentSwapStep];
      default:
        return undefined;
    };
  };

  const getOpeningStep = () => {
    const peerState = wasmGetPeerState();
    if (props.offerData?.swapCreator === peerState.id) {
      return <OfferCreatedStep offerId={props.offerId} />;
    }
    return <OfferConnecting offerId={props.offerId} />;
  };

  useEffect(() => {
    const fetchPeerState = setInterval(() => {
      const peerState = wasmGetPeerState();
      setPeerState(peerState.state);
    }, 100);

    return () => clearInterval(fetchPeerState);
  }, []);

  useEffect(() => {
    const startSwapListener = setInterval(() => {
      if (peerState === 'PeerCommunicating') {
        setSwapStarted(true);
        clearInterval(startSwapListener);
      }
    }, 100);

    return () => clearInterval(startSwapListener);
  }, [peerState, setSwapStarted]);

  return (
    <div className='bg-slate-200/20 w-screen backdrop-blur-sm fixed inset-0 flex items-center justify-center z-20'>
      <div className='bg-white text-black h-3/4 md:w-3/4 lg:w-2/4 xs:w-full xs:mx-2 rounded-lg shadow-lg z-20'>
        {/* Header */}
        <div className='flex flex-row w-full justify-between place-items-center p-4 border-b border-zinc-200 h-1/6 px-8'>
          <div className='flex flex-row gap-8 items-center'>
            <div className='text-3xl font-bold'>Swap Room</div>
            <div className='flex flex-col'>
              <div className='grid grid-cols-6'>
                <span className='text-slate-600 col-span-2'>Offer ID:</span>
                <span className='col-span-4'>1245789</span>
              </div>
              <div className='grid grid-cols-6'>
                <span className='text-slate-600 col-span-2'>Status:</span>
                <span className='col-span-4'>{peerState}</span>
              </div>
            </div>
          </div>
          <button onClick={props.onClose} className='bg-slate-100 p-1 rounded-lg hover:bg-red-300 hover:text-red-900 active:bg-red-400 duration-300 transition ease-in-out'>
            <CloseIcon />
          </button>
        </div>

        {!swapStarted && getOpeningStep()}
        {swapStarted && getCurrentStep()}
      </div>
    </div>
  );
};
