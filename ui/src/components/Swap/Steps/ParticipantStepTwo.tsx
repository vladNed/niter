import { type SwapStepProps } from 'types';
import { Clipboard } from 'components/Icons';
import { useWasm } from 'hooks';
import { TransactionRequestTypes } from 'localConstants';
import { useEffect, useState } from 'react';
import QRCode from 'react-qr-code';


export const ParticipantStepTwo = (props: SwapStepProps) => {
  const { getTransactionRequest } = useWasm();
  const [lockingAddress, setLockingAddress] = useState<string>('Not Available');

  const getAmount = () => {
    const amount = props.offerData?.isSwapCreator ? props.offerData?.sendingAmount : props.offerData?.receivingAmount;
    return `${amount} ${props.offerData?.sendingCurrency}`
  }

  useEffect(() => {
    const fetchLockingAddress = async () => {
      try {
        const txData = await getTransactionRequest(TransactionRequestTypes.CreateSwap);
        setLockingAddress(txData.address);
      } catch (e) {
        console.log(e);
      }
    }

    fetchLockingAddress();
  }, [lockingAddress, setLockingAddress]);

  return (
    <div className='h-5/6'>
      {/* Body */}
      <div className='flex flex-col gap-4 justify-center place-items-center h-2/3 px-8'>
        <p className='text-2xl font-bold text-center'>Lock BTC by sending to the following address</p>
        <QRCode value={lockingAddress} size={256} />
        <ul className='text-center'>
          <li>
            <div className='p-1 flex flex-row w-1/2 rounded-lg  place-items-center'>
              <span className='m-auto text-md font-semibold'>{lockingAddress}</span>
              <button className='p-1 text-slate-500 rounded-md hover:bg-slate-500 hover:text-slate-200 transition duration-200 ease-in-out'>
                <Clipboard />
              </button>
            </div>
          </li>
          <li className='text-md font-semibold'>Amount: {getAmount()}</li>
          <li>The address will be monitored to validate when the funds arrive.</li>
        </ul>
      </div>
    </div>
  )
};

