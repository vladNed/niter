import { CloseIcon } from "components/Icons"

export const NetworkField = ({ label, value }: { label: string, value: string}) => {
  return (
    <div className='flex justify-between font-kanit items-center py-1 grid grid-cols-12 gap-2'>
      <span className='text-lg col-span-6 text-neutral-300'>{label}:</span>
      <span className='text-lg col-span-6 text-neutral-100'>{value}</span>
    </div>
  )
}

export const NetworkModal = () => {
  const fields = [
    { label: 'ID', value: '#as9l2h8q0mr8j' },
    { label: 'Type', value: 'mainnet' },
    { label: 'Status', value: 'active' },
    { label: 'Version', value: '1.0.0' },
    { label: 'State', value: 'PeerIdle' },
  ]

  return (
    <div className='fixed inset-0 z-50 overflow-y-auto font-outfit overflow-x-hidden bg-zinc-900/30 backdrop-blur-sm items-center justify-center flex'>
      <div className='bg-neutral-800 border-[1px] border-neutral-700 shadow-lg shadow-zinc-900 rounded-xl p-6 min-w-[25rem] max-w-[25rem] min-h-[25rem]'>
        {/* Modal Header */}
        <div className='flex justify-between items-center border-b-[1px] border-b-neutral-500 pb-6'>
          <div className='flex flex-col'>
            <span className='text-2xl font-bold text-neutral-100'>Network</span>
            <span className='text-neutral-500'>Information about the Niter network node</span>
          </div>
          <button className='text-neutral-100 hover:text-red-400 hover:bg-neutral-700 transition duration-300 ease-in-out rounded-lg'>
            <CloseIcon />
          </button>
        </div>

        {/* Modal Body */}
        <div className='flex flex-col py-4'>
          {fields.map((field, index) => (
            <NetworkField key={index} label={field.label} value={field.value} />
          ))}
        </div>

      </div>
    </div>
  )
}