import { useState } from 'react'

const Step = ({
  step,
  stepIndex,
  currentStep
}: {
  step: string,
  stepIndex: number,
  currentStep: number
}) => {
  return (
    <li className={`
      flex w-full items-center
      after:content-[''] after:w-full after:h-1 after:border-b after:border-4 after:inline-block
      ${stepIndex < currentStep ? ' after:border-green-500 ' : ' after:border-gray-700 '}
      last:w-fit last:after:hidden
      `}
    >
      <span className={
        `flex items-center justify-center w-10 h-10 rounded-full
        lg:h-12 lg:w-12 shrink-0 ${currentStep === stepIndex ? 'bg-green-500 text-green-700' : 'bg-neutral-700'}
        ${stepIndex < currentStep ? 'text-green-400 bg-green-700' : 'text-neutral-300'}
        `}
      >
        {stepIndex + 1}
      </span>
    </li>
  )
}


export const SwapStepper = ({ steps, currentStep }: {steps: string[], currentStep: number }) => {
  return (
    <div className='w-full my-2'>
      <ol className='flex items-center w-full'>
        {steps.map((step, index) => (
          <Step key={index} step={step} stepIndex={index} currentStep={currentStep} />
        ))}
      </ol>
    </div>
  )
}