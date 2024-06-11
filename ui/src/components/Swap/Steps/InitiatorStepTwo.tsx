import { faSpinner } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { type SwapStepProps } from 'types';


export const InitiatorStepTwo = (props: SwapStepProps) => {
  return (
    <div className='h-5/6'>
      {/* Body */}
      <div className='flex flex-col gap-4 justify-center place-items-center h-2/3 px-8'>
        <p className='text-2xl font-bold text-center'>Waiting for peer to lock BTC.</p>
        <FontAwesomeIcon icon={faSpinner} spin size='2x' />
      </div>
    </div>
  );
};