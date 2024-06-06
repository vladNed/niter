import { useState } from 'react';
import { CreateOfferButton } from 'components/Button';
import { BTCLogo, MultiversxLogo } from 'components/Icons';
import {
  type CreateOfferProps,
  type OfferDetails,
  type SideToken,
  type SwapFieldProps,
} from 'types';
import { CoinCurrency, OfferSide } from 'localConstants';

const SwapField = (props: SwapFieldProps) => {
  return (
    <div className='px-6 py-4 rounded-xl text-black bg-slate-100'>
      <div className='text-xl font-medium'>{props.side}</div>
      <div className='w-full flex py-4 rounded-md justify-between place-items-left grid grid-cols-12'>
        <div className='flex gap-4 col-span-3'>
          {props.icon}
          <div className='flex flex-col place-content-center'>
            <div className='text-xl font-bold leading-none'>{props.ticker}</div>
            <div className='text-sm leading-tight text-zinc-500'>{props.name}</div>
          </div>
        </div>
        <input
          onChange={props.onChange}
          value={props.value}
          data-side={props.dataSide}
          type='text'
          placeholder='0'
          className='col-span-9 w-full text-right outline-none text-4xl text-center bg-inherit'
        />
      </div>
      <div className='h-[1rem]'></div>
    </div>
  )
}

const tokens: SideToken[] = [
  { ticker: 'BTC', name: CoinCurrency.BTC, icon: <BTCLogo /> },
  { ticker: 'EGLD', name: CoinCurrency.EGLD, icon: <MultiversxLogo /> }
]

export const CreateOffer = (props: CreateOfferProps) => {
  const [errMsg, setErrMsg] = useState<string>('')
  const [sendingToken, setSendingToken] = useState<SideToken>(tokens[0]);
  const [receivingToken, setReceivingToken] = useState<SideToken>(tokens[1]);
  const [sendingAmount, setSendingAmount] = useState<string>('');
  const [receivingAmount, setReceivingAmount] = useState<string>('');

  const handleSideSwap = () => {
    const temp = sendingToken
    setSendingToken(receivingToken)
    setReceivingToken(temp)
    const tempAmount = sendingAmount
    setSendingAmount(receivingAmount)
    setReceivingAmount(tempAmount)
  }

  const handleSwapAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    var value = e.target.value
    const side = e.target.getAttribute('data-side')

    if (value === '00') {
      value = '0.0'
    }

    const isDecimal = /^\d*\.?\d{0,4}$/.test(value)
    if (!isDecimal && value !== '') return

    switch (side) {
      case OfferSide.SENDING:
        setSendingAmount(value)
        break
      case OfferSide.RECEIVING:
        setReceivingAmount(value)
        break
      default:
        break
    }
  }

  const handleCreateOffer = () => {
    const sendingAmountNum = parseFloat(sendingAmount)
    const receivingAmountNum = parseFloat(receivingAmount)
    if (isNaN(sendingAmountNum) || isNaN(receivingAmountNum) || sendingAmountNum <= 0 || receivingAmountNum <= 0){
      setErrMsg('Not all amounts are provided')
      return
    }

    const peerState = wasmGetPeerState()
    const data: OfferDetails = {
      swapCreator: peerState.id,
      sendingAmount: sendingAmount,
      sendingCurrency: sendingToken.ticker,
      receivingAmount: receivingAmount,
      receivingCurrency: receivingToken.ticker,
      createdAt: new Date().toISOString(),
      expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24).toISOString()
    }
    props.handleReceiptShow(data)
  }


  return (
    <div>
      <div className='flex flex-col gap-1 relative'>
        <SwapField
          icon={sendingToken.icon}
          ticker={sendingToken.ticker}
          name={sendingToken.name}
          side='Swap'
          value={sendingAmount}
          onChange={handleSwapAmountChange}
          dataSide={OfferSide.SENDING}
        />
        <div
          className={`
              rounded-full bg-slate-100 ring ring-[4px] ring-white absolute inset-0
              p-2 w-10 h-10 flex place-items-center justify-center m-auto z-10
              hover:bg-primary-500 transition duration-300 ease-in-out hover:text-primary-100 font-bold
            `}
          onClick={handleSideSwap}
        >
          <svg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' strokeWidth='1.5' stroke='currentColor' className='h-5 w-5'>
            <path strokeLinecap='round' strokeLinejoin='round' d='M19.5 13.5 12 21m0 0-7.5-7.5M12 21V3' />
          </svg>
        </div>
        <SwapField
          icon={receivingToken.icon}
          ticker={receivingToken.ticker}
          name={receivingToken.name}
          side='For'
          value={receivingAmount}
          onChange={handleSwapAmountChange}
          dataSide={OfferSide.RECEIVING}
        />
      </div>
      <div className='mt-1'>
        <CreateOfferButton
          text={'Create offer'}
          onClick={handleCreateOffer}
        />
      </div>
      {errMsg && <div className='text-red-500 my-4 text-lg w-full text-center'>{errMsg}</div>}
    </div>
  )
}