import { CloseIcon } from 'components/Icons'


type SwapModalProps = {
  onClose: () => void
  offerId: string
  offerData?: string
  bodyElement?: JSX.Element
}

export const SwapModal = (props: SwapModalProps) => {

  return (
    <div className='bg-slate-200/20 w-screen backdrop-blur-sm fixed inset-0 flex items-center justify-center'>
      <div className='bg-white text-black h-3/4 md:w-3/4 lg:w-3/4 xs:w-full xs:mx-2 rounded-lg shadow-lg'>
        {/* Header */}
        <div className='flex flex-row w-full justify-between place-items-center p-4 border-b border-zinc-200 h-1/6 px-8'>
          <div className='text-3xl font-bold'>Swap Room</div>
          <button onClick={props.onClose} className='bg-slate-100 p-1 rounded-lg hover:bg-red-300 hover:text-red-900 active:bg-red-400 duration-300 transition ease-in-out'>
            <CloseIcon />
          </button>
        </div>

        {props.bodyElement}
      </div>
    </div>
  )
}
