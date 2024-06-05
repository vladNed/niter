
interface SwapButtonProps {
  onClick?: () => void;
  text: string;
}

export const CreateOfferButton = (props: SwapButtonProps) => {
  return (
    <button onClick={props.onClick} className='
      bg-primary-500 text-2xl font-kanit
      text-white rounded-xl px-4 py-4 w-full
      hover:bg-primary-700 active:bg-primary-800 transition-colors duration-300 ease-in-out
    '>
      {props.text}
    </button>
  );
}