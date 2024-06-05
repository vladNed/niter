import { CreateOfferButton } from "components/Button"


export const FindOffer = () => {
  return (
    <div className='flex flex-col gap-3 relative mb-4'>
      <input
        type="text"
        placeholder='Enter offer ID'
        className='w-full bg-slate-100 outline-none text-2xl p-2 rounded-xl text-center'
      />
     <CreateOfferButton text='Search' onClick={() => console.log("ceva")} />
    </div>
  )
}