import { Check } from "components/Icons";
import { useState } from "react";

type UserGuideSteps = {
  completed: boolean
  text: string
}

type StepProps = {
  offerId: string
}

export const OfferConnecting = (props: StepProps) => {
  const [steps, setSteps] = useState<UserGuideSteps[]>([
    { completed: false, text: 'Start the swap' },
  ]);

  const handleStartSwap = async () => {
    try{
      await wasmCreateAnswer(props.offerId)
    } catch (error) {
      console.log(error)
    }
  }

  return (
    <div className='h-5/6'>
      {/* Body */}
      <div className='flex flex-col gap-4 justify-center place-items-center h-2/3 px-8'>
        <div className='flex flex-row gap-2 items-center relative'>
          <span className='text-3xl font-medium'>Offer ID: {props.offerId}</span>
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

      <div className='flex flex-col h-1/3 items-center gap-2'>
        <button
          className='
              bg-primary-500 text-white text-2xl font-medium w-1/2 p-4
              rounded-lg shadow-md shadow-slate-300'
          onClick={handleStartSwap}
        >Start the swap </button>
      </div>
    </div>
  )
}