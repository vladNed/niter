import { CloseIcon } from 'components/Icons'
import {
  SwapEvents,
  SwapSide
} from 'localConstants';
import { useEffect, useState } from 'react'
import {
  InitiatorStepOne,
  InitiatorStepTwo,
  OfferConnecting,
  OfferCreatedStep,
  ParticipantStepOne,
  ParticipantStepTwo
} from 'components/Swap/Steps';
import {
  type OfferDetails,
} from 'types';
import { useWasm } from 'hooks';


type SwapModalProps = {
  offerId: string;
  offerData?: OfferDetails;
  swapSide?: SwapSide;
  handleStartFlow: () => void;
  handleClose: () => void;
};

type InitiatorStateMap = {
  [key in SwapEvents]: JSX.Element | undefined;
};

type ParticipantStateMap = {
  [key in SwapEvents]: JSX.Element | undefined;
};

export const SwapModal = (props: SwapModalProps) => {
  const [peerState, setPeerState] = useState<string>('');
  const [swapStarted, setSwapStarted] = useState<boolean>(false);
  const [swapEvents, setSwapEvents] = useState<SwapEvents[]>([]);
  const { getSwapEvents, resetPeer } = useWasm();
  const initiatorStatesMap: InitiatorStateMap = {
    [SwapEvents.SInit]: undefined,
    [SwapEvents.SInitDone]: <InitiatorStepOne offerData={props.offerData} />,
    [SwapEvents.SLockedEGLD]: <InitiatorStepTwo offerData={props.offerData} />,
    [SwapEvents.SLockeedBTC]: undefined,
    [SwapEvents.SRefund]: undefined,
    [SwapEvents.SClaimed]: undefined,
    [SwapEvents.SOk]: undefined,
    [SwapEvents.SFailed]: undefined
  };
  const participantStatesMap: ParticipantStateMap = {
    [SwapEvents.SInit]: undefined,
    [SwapEvents.SInitDone]:  <ParticipantStepOne />,
    [SwapEvents.SLockedEGLD]: <ParticipantStepTwo offerData={props.offerData} />,
    [SwapEvents.SLockeedBTC]: undefined,
    [SwapEvents.SClaimed]: undefined,
    [SwapEvents.SRefund]: undefined,
    [SwapEvents.SOk]: undefined,
    [SwapEvents.SFailed]: undefined
  };

  const onClose = () => {
    resetPeer();
    props.handleClose();
  };

  const getCurrentStep = () => {
    const lastEvent = swapEvents.slice(-1)[0];
    switch (props.swapSide) {
      case SwapSide.INITIATOR:
        return initiatorStatesMap[lastEvent];
      case SwapSide.PARTICIPANT:
        return participantStatesMap[lastEvent];
      default:
        return;
    }
  };

  const getOpeningStep = () => {
    const peerState = wasmGetPeerState();
    if (props.offerData?.swapCreator === peerState.id) {
      return <OfferCreatedStep offerId={props.offerId} />;
    }
    return <OfferConnecting offerId={props.offerId} />;
  };

  useEffect(() => {
    const fetchSwapEvents = setInterval(async () => {
      const swapEvents = await getSwapEvents();
      setSwapEvents(swapEvents);
    }, 500);
    return () => clearInterval(fetchSwapEvents);
  }, [setSwapEvents, swapEvents]);

  useEffect(() => {
    const fetchPeerState = setInterval(() => {
      const peerState = wasmGetPeerState();
      setPeerState(peerState.state);
    }, 500);

    return () => clearInterval(fetchPeerState);
  }, []);

  useEffect(() => {
    const startSwapListener = setInterval(() => {
      if (peerState === 'PeerCommunicating') {
        setSwapStarted(true);
        props.handleStartFlow();
        clearInterval(startSwapListener);
      }
    }, 1000);

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
          <button onClick={onClose} className='bg-slate-100 p-1 rounded-lg hover:bg-red-300 hover:text-red-900 active:bg-red-400 duration-300 transition ease-in-out'>
            <CloseIcon />
          </button>
        </div>

        {!swapStarted && getOpeningStep()}
        {swapStarted && getCurrentStep()}
      </div>
    </div>
  );
};
