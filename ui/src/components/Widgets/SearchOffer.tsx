import { CreateOfferButton } from "components/Button";
import { ConnectWalletPlaceholderBtn } from "components/Button/ConnectWalletButton";
import { MoneroLogo, MultiversxLogo } from "components/Icons";
import { useState } from "react";
import { SearchOfferWidgetProps, SideToken, SwapFieldProps } from "types";


const SwapField = (props: SwapFieldProps) => {
  return (
    <div className='bg-neutral-800 px-6 py-4 rounded-xl border-[1px] border-neutral-800 hover:border-[1px] hover:border-neutral-700 active:border-neutral-700'>
      <div className='text-xl font-medium text-neutral-500'>{props.side}</div>
      <div className='w-full flex py-4 rounded-md justify-between place-items-left grid grid-cols-12'>
        <div className='flex gap-4 col-span-3'>
          {props.icon}
          <div className='flex flex-col place-content-center text-neutral-200'>
            <div className='text-xl font-bold leading-none'>{props.ticker}</div>
            <div className='text-sm leading-tight'>{props.name}</div>
          </div>
        </div>
        <input type="text" placeholder='0' className='col-span-9 w-full text-right bg-neutral-800 outline-none text-4xl text-center' />
      </div>
      <div className='h-[1rem]'></div>
    </div>
  )
}

const tokens: SideToken[] = [
  { ticker: 'XMR', name: 'Monero', icon: <MoneroLogo /> },
  { ticker: 'EGLD', name: 'Multiversx', icon: <MultiversxLogo /> }
]

export const SearchOfferWidget = (props: SearchOfferWidgetProps) => {
  const [swapMode, setSwapMode] = useState<'Create' | 'Find'>('Find')
  const swapModeText = swapMode === 'Find' ? 'Find an existing offer' : 'Create a new offer'

  const [sendingToken, setSendingToken] = useState<SideToken>(tokens[0])
  const [receivingToken, setReceivingToken] = useState<SideToken>(tokens[1])

  const handleSwapMode = (mode: 'Create' | 'Find') => {
    setSwapMode(mode)
  }

  const handleSideSwap = () => {
    const temp = sendingToken
    setSendingToken(receivingToken)
    setReceivingToken(temp)
  }

  return (
    <div className='h-full text-white font-outfit rounded-md p-10 w-full min-w-[500px] max-w-[600px]'>
      <div className='mb-4 flex flex-col place-items-left gap-5'>
        <div className='flex gap-4'>
          <div
            className={
              `px-4 py-2 rounded-2xl hover:bg-neutral-700 hover:text-neutral-400 transition duration-500 ease-in-out`
              + (swapMode === 'Create' ? ' bg-neutral-700 text-neutral-400' : ' bg-neutral-800 text-neutral-500')
            }
            onClick={() => handleSwapMode('Create')}
          >Create
          </div>
          <div
            className={
              `px-4 py-2 rounded-2xl hover:bg-neutral-700 hover:text-neutral-400 transition duration-500 ease-in-out`
              + (swapMode === 'Find' ? ' text-neutral-400 bg-neutral-700' : ' bg-neutral-800 text-neutral-500')
            }
            onClick={() => handleSwapMode('Find')}
          >Find
          </div>
          <div className='px-4 py-2 rounded-2xl hover:bg-neutral-700 hover:text-neutral-400 bg-neutral-800 text-neutral-500 transition duration-500 ease-in-out' onClick={handleSideSwap}>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="w-6 h-6 hover:rotate-180 transition duration-300 ease-in-out">
              <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 0 0-3.7-3.7 48.678 48.678 0 0 0-7.324 0 4.006 4.006 0 0 0-3.7 3.7c-.017.22-.032.441-.046.662M19.5 12l3-3m-3 3-3-3m-12 3c0 1.232.046 2.453.138 3.662a4.006 4.006 0 0 0 3.7 3.7 48.656 48.656 0 0 0 7.324 0 4.006 4.006 0 0 0 3.7-3.7c.017-.22.032-.441.046-.662M4.5 12l3 3m-3-3-3 3" />
            </svg>
          </div>
        </div>
        <span className='text-left text-2xl text-neutral-500'>{swapModeText}</span>
      </div>
      <div className='flex flex-col gap-1'>
        <SwapField icon={sendingToken.icon} ticker={sendingToken.ticker} name={sendingToken.name} side='Swap' />
        <SwapField icon={receivingToken.icon} ticker={receivingToken.ticker} name={receivingToken.name} side='For' />
        {props.isPlaceholder ?
          <ConnectWalletPlaceholderBtn /> :
          <CreateOfferButton text={swapMode === 'Find' ? 'See offers' : 'Create offer'} />
        }
        {swapMode === 'Find' && <span className='text-center text-neutral-500 px-10'>Tip: If there are not active offers, don't worry, you can be the one to create the offer.</span>}
      </div>
    </div>
  )
}