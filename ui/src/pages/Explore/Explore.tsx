import { useState } from "react"
import { AuthRedirectWrapper } from "wrappers"


export const Explore = () => {
  const columns = [
    'ID',
    'You send',
    'You get',
    'Time left',
    'Action'
  ]
  const [sendToken, setSendToken] = useState<string>('XMR')
  const [getToken, setGetToken] = useState<string>('EGLD')

  const handleSendTokenChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    switch (e.target.value) {
      case 'XMR':
        if(e.target.id === 'sendOption') {
          setSendToken('XMR')
          setGetToken('EGLD')
        } else {
          setSendToken('EGLD')
          setGetToken('XMR')
        }
        break
      case 'EGLD':
        if(e.target.id === 'sendOption') {
          setSendToken('EGLD')
          setGetToken('XMR')
        } else {
          setSendToken('XMR')
          setGetToken('EGLD')
        }
        break
    }
  }

  return (
    <div className='h-full container flex flex-col text-neutral-200'>
      <div className='flex flex-row mb-6 justify-between'>
        <div className='flex place-items-center gap-4'>
          <div className='text-xl font-bold'>Swap</div>
          <select id='sendOption' value={sendToken} className='text-neutral-200 bg-neutral-800 px-2 py-1 rounded-xl' onChange={handleSendTokenChange}>
            <option value="XMR" selected>XMR</option>
            <option value="EGLD">EGLD</option>
          </select>
          <div className='text-xl font-bold'>for</div>
          <select id='getOption' value={getToken} className='text-neutral-200 bg-neutral-800 px-2 py-1 rounded-xl' onChange={handleSendTokenChange}>
            <option value="XMR">XMR</option>
            <option value="EGLD" selected>EGLD</option>
          </select>
        </div>
        <input type='text' className='px-4 bg-neutral-800 outline-none border-[1px] border-neutral-700 rounded-xl' placeholder='Enter amount'/>
        <button className='bg-neutral-800 text-neutral-400 py-2 px-4 rounded-xl'>Clear all</button>
        <select className='text-neutral-200 bg-neutral-800 px-2 py-1 rounded-xl'>
          <option value='' selected disabled>Sort by</option>
          {columns.map((column) => (
            <option key={column} value={column}>{column}</option>
          ))}
        </select>
      </div>
      <div className='h-full bg-neutral-900 rounded-xl flex flex-col border-[1px] border-neutral-800 pb-6'>
        <table className="table-auto">
          <thead>
            <tr>
              {columns.map((column) => (
                <th key={column} className='py-4'>{column}</th>
              ))}
            </tr>
          </thead>
          <tbody className="text-center bg-neutral-800 ">
            <tr className='border-b-2 border-neutral-700'>
              <td className='py-4 '><span className='text-neutral-400'>s-</span>12i9lp3r</td>
              <td className='py-4'>1000 XMR</td>
              <td className='py-4'>1 EGLD</td>
              <td className='py-4'>10 minutes</td>
              <td className='py-4'><button className='bg-secondary-500 px-2 py-1 rounded-md text-neutral-800 hover:bg-secondary-400 hover:text-neutral-900'>Connect</button></td>
            </tr>
            <tr className='border-b-2 border-neutral-700'>
              <td className='py-4 '><span className='text-neutral-400'>s-</span>1af9ld3r</td>
              <td className='py-4'>1000 EGLD</td>
              <td className='py-4'>1 XMR</td>
              <td className='py-4'>10 minutes</td>
              <td className='py-4'><button className='bg-secondary-500 px-2 py-1 rounded-md text-neutral-800  hover:bg-secondary-400 hover:text-neutral-900'>Connect</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  )
}