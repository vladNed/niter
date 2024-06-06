import { useSwapRouterTransactions } from "hooks";
import { CoinCurrency, TransactionRequestTypes } from "localConstants";
import { useEffect, useState } from "react";
import { type SwapStepProps } from "types";

export const InitiatorStepOne = (props: SwapStepProps) => {
  const [errMsg, setErrMsg] = useState<string>('');

  const { sendCreateSwapTransaction, transactionStatus } = useSwapRouterTransactions();

  const onCreateSwap =  async() => {
    let txData: any;
    try{
      txData = await wasmTransactionRequest(TransactionRequestTypes.CreateSwap.toString());
    } catch (e) {
      setErrMsg('Cannot create transaction');
      return;
    }
    await sendCreateSwapTransaction({
      amount: txData.amount,
      claimProof: txData.claimProof,
      refundProof: txData.refundProof,
      callbackRoute: 'initiator-create-swap'
    });
  };

  const getAmount = () => {
    const amount = props.offerData?.sendingCurrency === CoinCurrency.EGLD ? props.offerData?.sendingAmount : props.offerData?.receivingAmount;
    const currency = props.offerData?.sendingCurrency === CoinCurrency.EGLD ? CoinCurrency.EGLD : props.offerData?.receivingCurrency;
    return `${amount} ${currency}`;
  }

  useEffect(() => {
    if (transactionStatus.isSuccessful) {
      props.progressHandler();
    }
  }, [transactionStatus]);

  return (
    <div className='h-5/6'>
      {/* Body */}
      <div className='flex flex-col gap-4 justify-center place-items-center h-2/3 px-8'>
        <div className='text-2xl font-medium'>Lock funds:</div>
        <ul className='text-xl justify-center'>
          <li>
            <div className='flex flex-row gap-4 grid grid-cols-2'>
              <span className='col-span-1'>Amount:</span>
              <span className='col-span-1 text-right'>{getAmount()}</span>
            </div>
          </li>
          <li>
            <div className='flex flex-row gap-4 grid grid-cols-2'>
              <span className='col-span-1 '>Lock for:</span>
              <span className='col-span-1 text-right'>Not Implemented</span>
            </div>
          </li>
          <li>
            <div className='flex flex-row gap-4 grid grid-cols-2'>
              <span className='col-span-1'>Refund after:</span>
              <span className='col-span-1 text-right'>Not Implemented</span>
            </div>
          </li>
        </ul>
      </div>

      <div className='flex flex-col h-1/3 items-center gap-2'>
        <button
          className='bg-primary-500 text-white text-2xl font-medium w-1/2 p-4 rounded-lg shadow-md shadow-slate-300 hover:bg-primary-800 hover:text-white transition duration-200 ease-in-out'
          onClick={onCreateSwap}
        >Lock funds</button>
        <span className='text-xl font-medium text-red-500'>{errMsg}</span>
      </div>
    </div>
  )
}
