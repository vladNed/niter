import { Check, Clipboard } from 'components/Icons'
import { useState } from 'react'

type UserGuideSteps = {
  completed: boolean
  text: string
}

type StepProps = {
  offerId: string
}


export const OfferCreatedStep = (props: StepProps) => {
  const [steps, setSteps] = useState<UserGuideSteps[]>([
    { completed: false, text: 'Copy and share the offer id with a friend*' },
    { completed: false, text: 'Wait for your friend start the swap' },
  ]);
  const [copyPopup, setCopyPopup] = useState<boolean>(false);

  const handleCopyClipboard = () => {
    navigator.clipboard.writeText(props.offerId).then(() => {
      setSteps(steps.map((step, index) => index === 0 ? { ...step, completed: true } : step))
      setCopyPopup(true)
      setTimeout(() => {
        setCopyPopup(false)
      }, 1000)
    })
  }

  return (
    <div className='h-5/6'>
      {/* Body */}
      <div className='flex flex-col gap-4 justify-center place-items-center h-2/3 px-8'>
        <div className='flex flex-row gap-2 items-center relative'>
          <span className='text-3xl font-medium'>Offer ID: {props.offerId}</span>
          <button onClick={handleCopyClipboard} className='bg-slate-100 p-2 rounded-lg text-black hover:bg-blue-400 hover:text-white transition duration-200 ease-in-out active:bg-blue-700'>
            <Clipboard />
          </button>
          {copyPopup && <span className='absolute bg-slate-200 p-1 text-slate-500 rounded-lg -end-[7rem]'>Offer copied !!</span>}
        </div>
        <div>What to do next:</div>
        <ul className='flex flex-col gap-4'>
          {steps.map((step, index) => (
            <li key={index} className='flex flex-row gap-2 items-center'>
              <div className={`rounded-full h-8 w-8 flex items-center justify-center ${step.completed ? 'bg-green-400' : 'bg-slate-300 '} text-white`}>{step.completed ? <Check /> : index + 1}</div>
              <div className={`text-lg ${step.completed ? 'text-black' : 'text-slate-400'}`}>{step.text}</div>
            </li>
          ))}
        </ul>
      </div>

      {/* Footer */}
      <div className='flex flex-col h-1/3 items-center gap-2'>
        <button
          className='
              bg-primary-500 text-white text-2xl font-medium w-1/2 p-4
              rounded-lg shadow-md shadow-slate-300 disabled:bg-primary-400 disabled:text-white'
          disabled
        >
          Wait for peer to start
        </button>
        <span className='text-slate-400'>*Offer will disappear if you exit this window</span>
      </div>

    </div>
  )
}