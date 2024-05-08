
interface SwapButtonProps {
  onClick?: () => void;
  text: string;
}

export const CreateOfferButton = (props: SwapButtonProps) => {
  return (
    <button onClick={props.onClick} className='
      bg-primary-600/40 text-2xl font-kanit
      text-primary-300 rounded-xl px-4 py-4
      hover:bg-primary-700 active:bg-primary-800 transition-colors hover:text-white hover:shadow-md hover:shadow-primary-900
    '>
      {props.text}
    </button>
  );
}