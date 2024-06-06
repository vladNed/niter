import { CreateOfferButton } from 'components/Button'
import { useState } from 'react'
import { type OfferDetails } from 'types'

type FindOfferProps = {
  handleReceiptShow: (offerData: OfferDetails, offerId?: string) => void
}

export const FindOffer = (props: FindOfferProps) => {
  const [offerId, setOfferId] = useState<string>('')
  const [errMsg, setErrMsg] = useState<string | null>(null)

  const handleOnValueChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setOfferId(e.target.value)
  }

  const handleSearch = async () => {
    if (!offerId) {
      setErrMsg('Offer ID is required')
      return
    }
    const offers = await wasmPollOffers();
    const offer = offers.find(offer => offer.id === offerId)
    if (!offer) {
      setErrMsg('Offer not found')
      return
    }
    setErrMsg(null)
    const offerData: OfferDetails = {
      swapCreator: offer.swapCreator,
      sendingAmount: offer.sendingAmount,
      sendingCurrency: offer.sendingCurrency,
      receivingAmount: offer.receivingAmount,
      receivingCurrency: offer.receivingCurrency,
    }
    props.handleReceiptShow(offerData, offerId)
  }

  return (
    <div className='flex flex-col gap-3 relative mb-2'>
      <input
        type='text'
        placeholder='Enter offer ID'
        className='w-full bg-slate-100 outline-none text-2xl p-2 rounded-xl text-center'
        onChange={handleOnValueChange}
        value={offerId}
      />
     <CreateOfferButton text='Search' onClick={handleSearch} />
     {errMsg && <div className='text-red-500 my-2 text-lg w-full text-center'>{errMsg}</div>}
    </div>
  )
}