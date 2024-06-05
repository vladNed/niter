import { CreateOfferButton } from 'components/Button'
import { OfferDetails } from 'types'

type ReceiptOfferProps = {
  offerData: OfferDetails | null
  handleConfirmation: (offerId: string) => void
}

export const ReceiptOffer = (props: ReceiptOfferProps) => {

  const calculateProtocolFee = () => {
    if (!props.offerData) return 0
    const sendingAmount = parseFloat(props.offerData.sendingAmount)
    return sendingAmount * 0.005
  }

  const handleConfirmSwap = async () => {
    const offerId = await wasmCreateOffer(JSON.stringify(props.offerData))
    props.handleConfirmation(offerId)
  }

  return (
    <div className='text-xl flex flex-col gap-4'>
      <div className='text-2xl font-medium py-2 border-b-[1px] border-zinc-200'>Swap Details</div>
      <div className='flex flex-col'>
        <span className='text-zinc-500'>You send:</span>
        <span>{props.offerData?.sendingAmount} {props.offerData?.sendingCurrency}</span>
      </div>
      <div className='flex flex-col'>
        <span className='text-zinc-500'>Exchange Rate:</span>
        <span>Not Implemented</span>
      </div>
      <div className='flex flex-col'>
        <span className='text-zinc-500'>Protocol fee (0.5%):</span>
        <span>{calculateProtocolFee()} {props.offerData?.sendingCurrency}</span>
      </div>
      <div className='flex flex-row items-center gap-4 border-b-[1px] border-zinc-200 pb-4'>
        <span className='text-zinc-500'>Network fee estimate:</span>
        <svg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' strokeWidth='1.5' stroke='currentColor' className='w-6 h-6 text-zinc-500'>
          <path strokeLinecap='round' strokeLinejoin='round' d='m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z' />
        </svg>
      </div>
      <div className='flex flex-col'>
        <span className='text-2xl font-medium'>You receive:</span>
        <span>{props.offerData?.receivingAmount} {props.offerData?.receivingCurrency}</span>
      </div>
      <CreateOfferButton text='Confirm Swap' onClick={handleConfirmSwap} />
    </div>
  )
}
