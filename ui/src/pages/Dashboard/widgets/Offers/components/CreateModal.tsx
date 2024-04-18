import { faXmark, faArrowRight } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { NumericInput } from '.';
import { useState } from 'react';
import { Button } from 'components';
import { Coin, OfferData } from 'types';

export const CreateModal = (props: {
  onExit: () => void;
  onSubmit: (data: OfferData) => void;
}) => {
  const [amount1, setAmount1] = useState<string>('');
  const [amount2, setAmount2] = useState<string>('');

  const [receivingCoin, setReceivingCoin] = useState<Coin>(Coin.EGLD);
  const [givingCoin, setGivingCoin] = useState<Coin>(Coin.MNR);

  const handleSubmit = () => {
    const data: OfferData = {
      receivingAmount: amount1,
      givingAmount: amount2,
      receivingCoin: receivingCoin,
      givingCoin: givingCoin
    };

    props.onSubmit(data);
  }

  const handleSwitch = () => {
    const temp = givingCoin;
    setGivingCoin(receivingCoin);
    setReceivingCoin(temp);
  }

  return (
    <div className='absolute top-0 left-0 w-full h-full bg-black bg-opacity-50 flex justify-center items-center'>
      <div className='bg-white p-6 rounded-xl justify-center flex flex-col gap-6 w-1/2 shadow-xl shadow-gray-500'>

        {/* This is the modal header */}
        <div className='flex justify-between border-b-[1px] border-gray-200 py-2'>
          <div className='font-bold text-xl'>Create Offer</div>
          <button onClick={props.onExit}>
            <FontAwesomeIcon
              icon={faXmark}
              size='lg'
              className='p-2 rounded-full bg-slate-100 hover:bg-red-200 hover:text-red-500 transition
                         ease-in-out duration-200 text-slate-400'
            />
          </button>
        </div>

        {/* This is the modal body */}
        <div className='flex flex-row gap-6 justify-between py-10'>
          <div className='flex flex-col gap-4 w-full'>
            <NumericInput
              value={amount1}
              onChange={(e) => setAmount1(e)}
              placeholder='0.0'
            />
            <div className='text-lg flex gap-2'>
              Pay:
              <span className='font-extrabold text-blue-500 bg-blue-200 rounded-lg px-2'>{givingCoin}</span>
            </div>
          </div>
          <FontAwesomeIcon
            icon={faArrowRight}
            size='lg'
            className='p-2 bg-slate-200 rounded-full text-blue-600'
          />
          <div className='flex flex-col gap-4 w-full'>
            <NumericInput
              value={amount2}
              onChange={setAmount2}
              placeholder='0.0'
            />
            <div className='text-lg flex gap-2'>
              Receive:
              <span className='font-extrabold text-blue-500 bg-blue-200 rounded-lg px-2'>{receivingCoin}</span>
            </div>
          </div>
        </div>

        {/* This is the modal footer */}
        <div className='border-t-[1px] py-2 flex flex-row gap-4'>
          <Button onClick={handleSubmit}>Create</Button>
          <Button onClick={handleSwitch}>Switch</Button>
        </div>
      </div>
    </div>
  )
}