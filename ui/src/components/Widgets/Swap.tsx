import { useEffect, useState } from 'react';
import { FindOffer, CreateOffer, ReceiptOffer } from 'components/Widgets';
import { type OfferDetails } from 'types';
import { SwapModal } from 'components/Swap';
import { CoinCurrency, SwapSide, SwapWidgetType } from 'localConstants';


export const SwapWidget = () => {
  const [currentOffer, setCurrentOffer] = useState<OfferDetails | undefined>(undefined);
  const [currentOfferId, setCurrentOfferId] = useState<string>('');
  const [swapWidgetType, setSwapWidgetType] = useState<SwapWidgetType>(SwapWidgetType.CREATE);
  const [swapActive, setSwapActive] = useState<boolean>(false);
  const [swapFlowStarted, setSwapFlowStarted] = useState<boolean>(false);
  const [swapSide, setSwapSide] = useState<SwapSide | undefined>(undefined);
  const [isSwapCreator, setIsSwapCreator] = useState<boolean>(false);
  const swapModeText = swapWidgetType === SwapWidgetType.FIND ? 'Find an existing offer' : 'Create a new offer';

  const handleSwapStart = () => {
    setSwapFlowStarted(true);
  };

  const handleSwapClose = () => {
    setSwapActive(false);
    setCurrentOfferId('');
    setCurrentOffer(undefined);
    setSwapWidgetType(SwapWidgetType.CREATE);
    setSwapFlowStarted(false);
  }

  const handleSwapMode = (mode: SwapWidgetType) => {
    setSwapWidgetType(mode)
  }

  const handleReceiptOffer = (offerData: OfferDetails, offerId?: string) => {
    const peerId = wasmGetPeerState().id;
    const isSwapCreator = offerData.swapCreator === peerId;
    const isInitiator = (isSwapCreator && offerData.sendingCurrency === CoinCurrency.EGLD) || (!isSwapCreator && offerData.receivingCurrency === CoinCurrency.EGLD);
    offerData.isSwapCreator = isSwapCreator;
    setSwapSide(isInitiator ? SwapSide.INITIATOR : SwapSide.PARTICIPANT);
    setSwapWidgetType(SwapWidgetType.RECEIPT);
    setCurrentOffer(offerData);
    setIsSwapCreator(isSwapCreator);

    if(offerId) {
      setCurrentOfferId(offerId);
    };
  };

  const handleCreateConfirmation = async () => {
    const offerId = await wasmCreateOffer(JSON.stringify(currentOffer));
    setCurrentOfferId(offerId);
    setSwapActive(true);
  };

  const handleSearchConfirmation = async () => {
    setSwapActive(true);
  };

  const getConfirmationHandler = () => {
    return isSwapCreator ? handleCreateConfirmation : handleSearchConfirmation;
  };

  useEffect(() => {
    const fetchPeerState = setInterval(() => {
      const peerState = wasmGetPeerState();
      if(peerState.state === 'PeerIdle' && swapFlowStarted) {
        handleSwapClose();
      }
    }, 500);

    return () => clearInterval(fetchPeerState);
  }, [handleSwapClose, swapActive, currentOfferId, currentOffer, swapWidgetType, isSwapCreator, swapFlowStarted]);

  return (
    <div className='h-full text-black font-outfit p-4 w-full min-w-[500px] max-w-[600px] rounded-xl bg-white'>
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
      {swapWidgetType === SwapWidgetType.RECEIPT && <ReceiptOffer offerData={currentOffer} handleConfirmation={getConfirmationHandler()} isSwapCreator={isSwapCreator}/>}
      {swapWidgetType === SwapWidgetType.CREATE && <CreateOffer handleReceiptShow={handleReceiptOffer} />}
      {swapWidgetType === SwapWidgetType.FIND && <FindOffer handleReceiptShow={handleReceiptOffer}/>}
      {swapActive && <SwapModal
        offerId={currentOfferId}
        offerData={currentOffer}
        swapSide={swapSide}
        handleStartFlow={handleSwapStart}
        handleClose={handleSwapClose}
        />}
    </div>
  )
}