
export const PeerInfo = () => {
  return (
    <div className='w-full rounded-lg p-6 text-white grid gap-6'>
      <div className=''>
        <h1 className='text-3xl text-white font-kanit'>Peer Info</h1>
        <span className='text-[12px] text-zinc-400'>Data about the peer connection</span>
      </div>
      <div>
        <div className='grid grid-cols-12 flex flex-row'>
          <div className='col-span-4 font-bold text-zinc-400'>Network ID:</div>
          <div className='col-span-8'>#as3as6sq2</div>
        </div>
        <div className='grid grid-cols-12 flex flex-row'>
          <div className='col-span-4 font-bold text-zinc-400'>Status:</div>
          <div className='col-span-8'>
            <span className='flex gap-2 place-items-center'>
              Waiting Peer
              <span className="relative flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-500 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
              </span>
            </span>
          </div>
        </div>
        <div className='grid grid-cols-12 flex flex-row'>
          <div className='col-span-4 font-bold text-zinc-400'>Active Offers:</div>
          <div className='col-span-8'>1</div>
        </div>
      </div>
    </div>
  )
}
