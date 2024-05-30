import { OFFERS_POLLING_INTERVAL } from 'config'
import { useEffect, useState } from 'react'

type OfferType = {
  id: string,
  sendAmount?: string,
  sendToken?: string,
  getAmount?: string,
  getToken?: string,
  timeLeft?: string
}

export const ExploreWidget = () => {
  const columns = [
    'ID',
    'You send',
    'You get',
    'Time left',
    'Action'
  ]
  const [sendToken, setSendToken] = useState<string>('BTC')
  const [getToken, setGetToken] = useState<string>('EGLD')
  const [offers, setOffers] = useState<OfferType[]>([])
  const [offerIds, setOfferIds] = useState<string[]>([])

  const handleSendTokenChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    switch (e.target.value) {
      case 'BTC':
        if (e.target.id === 'sendOption') {
          setSendToken('BTC')
          setGetToken('EGLD')
        } else {
          setSendToken('EGLD')
          setGetToken('BTC')
        }
        break
      case 'EGLD':
        if (e.target.id === 'sendOption') {
          setSendToken('EGLD')
          setGetToken('BTC')
        } else {
          setSendToken('BTC')
          setGetToken('EGLD')
        }
        break
    }
  }



  useEffect(() => {
    const onOffersPool = async () => {
      try {
        const newOffers = await wasmPollOffers();
        for(let i = 0; i < newOffers.length; i++) {
          if (!offerIds.includes(newOffers[i].id)) {
            let offer = newOffers[i];
            setOffers([...offers, offer]);
            setOfferIds([...offerIds, offer.id]);
          }
        }
      } catch (e) {
        console.error('Error polling offers', e);
      }
    }

    onOffersPool();

    const offersPool = setInterval(onOffersPool, OFFERS_POLLING_INTERVAL);

    return () => clearInterval(offersPool);


  }, [offers, offerIds, setOffers, setOfferIds]);

  return (
    <div className='h-full container flex flex-col text-neutral-200 min-w-[1000px]'>
      <div className='flex flex-row mb-6 justify-between'>
        <div className='flex place-items-center gap-4'>
          <div className='text-xl font-bold'>Swap</div>
          <select id='sendOption' value={sendToken} className='text-neutral-200 bg-neutral-800 px-2 py-1 rounded-xl' onChange={handleSendTokenChange}>
            <option key='BTC' value='BTC'>BTC</option>
            <option key='EGLD' value='EGLD'>EGLD</option>
          </select>
          <div className='text-xl font-bold'>for</div>
          <select id='getOption' value={getToken} className='text-neutral-200 bg-neutral-800 px-2 py-1 rounded-xl' onChange={handleSendTokenChange}>
            <option key='BTC' value='BTC'>BTC</option>
            <option key='EGLD' value='EGLD'>EGLD</option>
          </select>
        </div>
        <input type='text' className='px-4 bg-neutral-800 outline-none border-[1px] border-neutral-700 rounded-xl' placeholder='Enter amount' />
        <button className='bg-neutral-800 text-neutral-400 py-2 px-4 rounded-xl'>Clear all</button>
        <select className='text-neutral-200 bg-neutral-800 px-2 py-1 rounded-xl' defaultValue='Sort by'>
          <option value='Sort by' disabled>Sort by</option>
          {columns.map((column) => (
            <option key={column} value={column}>{column}</option>
          ))}
        </select>
      </div>
      <div className='h-full bg-neutral-900 rounded-xl flex flex-col border-[1px] border-neutral-800 pb-6'>
        <table className='table-auto'>
          <thead>
            <tr>
              {columns.map((column) => (
                <th key={column} className='py-4'>{column}</th>
              ))}
            </tr>
          </thead>
          <tbody className='text-center bg-neutral-800 '>
            {offers.length > 0 && offers.map((offer) => (
              <tr key={offer.id} className='border-b-2 border-neutral-700'>
                <td className='py-4 '><span className='text-neutral-400'>s-</span>{offer.id}</td>
                <td className='py-4'>{offer.getAmount || 'N/A'}</td>
                <td className='py-4'>{offer.sendAmount || 'N/A'}</td>
                <td className='py-4'>{offer.timeLeft || 'N/A'}</td>
                <td className='py-4'><button className='bg-secondary-500 px-2 py-1 rounded-md text-neutral-800 hover:bg-secondary-400 hover:text-neutral-900'>Connect</button></td>
              </tr>
            ))}
          </tbody>
        </table>
        {offers.length == 0 && <div className='text-center text-2xl mt-5'>No offers available</div>}
      </div>
    </div>
  )
}