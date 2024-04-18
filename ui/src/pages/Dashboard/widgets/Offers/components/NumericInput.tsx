type TextInputProps = {
  onChange: (v: string) => void;
  value: string;
  placeholder: string;
}

export const NumericInput = ({ onChange, value, placeholder }: TextInputProps) => {

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let newValue = e.target.value.replace(/[^\d.]/g, '');
    newValue = newValue.replace(/\.+/g, '.');
    newValue = newValue.replace(/^0+(?=\d)|(?<=\.)0+(?=\d)|^\.+/g, '');
    onChange(newValue);
  }

  return (
    <div className='flex flex-col'>
      <input
        type='text'
        value={value}
        onChange={handleChange}
        placeholder={placeholder}
        className='px-2 border-b-[2px] border-gray-300 outline-none duration-100 ease-in-out transition
        active:border-blue-500 focus:border-blue-500 focus:border-b-[2px] w-full'
      />
    </div>
  );
}