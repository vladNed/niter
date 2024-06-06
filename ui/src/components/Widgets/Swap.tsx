import { useEffect, useState } from 'react';
import { FindOffer, CreateOffer, ReceiptOffer } from 'components/Widgets';
import { type OfferDetails } from 'types';
import {
  InitiatorStepOne,
  OfferConnecting,
  OfferCreatedStep,
  ParticipantStepOne,
  SwapModal,
} from 'components/Swap';
import { SwapSide, SwapWidgetType } from 'localConstants';


export const SwapWidget = () => {
  const [currentOffer, setCurrentOffer] = useState<OfferDetails | null>(null);
  const [currentOfferId, setCurrentOfferId] = useState<string>('');
  const [swapWidgetType, setSwapWidgetType] = useState<SwapWidgetType>(SwapWidgetType.CREATE);
  const [swapActive, setSwapActive] = useState<boolean>(false);
  const [swapSide, setSwapSide] = useState<SwapSide | undefined>(undefined);
  const [currentSwapStep, setCurrentSwapStep] = useState<number>(0);
  const swapModeText = swapWidgetType === SwapWidgetType.FIND ? 'Find an existing offer' : 'Create a new offer';

  const InitiatorStates = [
    <OfferCreatedStep offerId={currentOfferId} />,
    <InitiatorStepOne />
  ]

  const ReceiverStates = [
    <OfferConnecting offerId={currentOfferId}/>,
    <ParticipantStepOne />
  ]

  const getCurrentStep = () => {
    switch(swapSide) {
      case SwapSide.INITIATOR:
        return InitiatorStates[currentSwapStep]
      case SwapSide.PARTICIPANT:
        return ReceiverStates[currentSwapStep]
      default:
        return undefined
    }
  }

  const handleSwapClose = () => {
    setSwapActive(false)
    setCurrentOfferId('')
    setCurrentOffer(null)
    setSwapWidgetType(SwapWidgetType.CREATE)
  }

  const handleSwapMode = (mode: SwapWidgetType) => {
    setSwapWidgetType(mode)
  }

  const handleReceiptOffer = (offerData: OfferDetails, offerId?: string) => {
    if(swapWidgetType === 'Create') {
      setSwapSide(SwapSide.INITIATOR)
    } else {
      setSwapSide(SwapSide.PARTICIPANT)
    }
    setSwapWidgetType(SwapWidgetType.RECEIPT)
    setCurrentOffer(offerData)

    if(offerId) {
      setCurrentOfferId(offerId)
    }

  }

  const handleCreateConfirmation = async () => {
    const offerId = await wasmCreateOffer(JSON.stringify(currentOffer))
    setCurrentOfferId(offerId)
    setSwapActive(true)
  }

  const handleSearchConfirmation = async () => {
    setSwapActive(true)
  }

  useEffect(() => {
    const fetchPeerState = setInterval(() => {
      const peerState = wasmGetPeerState();
      switch (peerState.state) {
        case 'PeerCommunicating':
          setCurrentSwapStep(1)
          break;
        default:
          break;
      }
    }, 1000);

    return () => clearInterval(fetchPeerState)
  }, [])

  return (
    <div className='h-full text-black font-outfit rounded-md p-4 w-full min-w-[500px] max-w-[600px] rounded-xl bg-white'>
      <div className='mb-4 flex flex-col place-items-left gap-5'>
        <div className='p-2 flex gap-4 rounded-2xl bg-slate-100'>
          <button
            className={
              `w-1/2 px-4 py-2 rounded-2xl hover:bg-white hover:text-blue-600 transition duration-500 ease-in-out`
              + (swapWidgetType !== SwapWidgetType.FIND ? ' bg-white text-blue-600 font-medium' : ' text-neutral-500')
            }
            onClick={() => handleSwapMode(SwapWidgetType.CREATE)}
          >Create Offer
          </button>
          <button
            className={
              `w-1/2 px-4 py-2 rounded-2xl hover:bg-white hover:text-blue-600 transition duration-500 ease-in-out`
              + (swapWidgetType === SwapWidgetType.FIND ? ' bg-white text-blue-600 font-medium' : ' text-neutral-500')
            }
            onClick={() => handleSwapMode(SwapWidgetType.FIND)}
          >Find Offer
          </button>
        </div>
        <span className='text-left text-2xl'>{swapModeText}</span>
      </div>
      {swapWidgetType === SwapWidgetType.RECEIPT &&
      <ReceiptOffer
        offerData={currentOffer}
        handleConfirmation={swapSide === SwapSide.INITIATOR ? handleCreateConfirmation : handleSearchConfirmation}
        swapSide={swapSide}
      />}
      {swapWidgetType === SwapWidgetType.CREATE && <CreateOffer handleReceiptShow={handleReceiptOffer} />}
      {swapWidgetType === SwapWidgetType.FIND && <FindOffer handleReceiptShow={handleReceiptOffer}/>}
      {swapActive && <SwapModal onClose={handleSwapClose} offerId={currentOfferId} bodyElement={getCurrentStep()}/>}
    </div>
  )
}