
type SelectInputProps = {
  label: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  options: {
    label: string;
    value: string;
  }[];
}

export const SelectInput = (props: SelectInputProps) => {
  return (
    <div className='flex flex-col'>
      <label>{props.label}</label>
      <select
        value={props.value}
        onChange={props.onChange}
        className='px-2 border-black border-[1px] border-slate-300 rounded-lg
        active:border-blue-500 focus-visible:border-blue-500 bg-white'
      >
        {props.options.map((option, index) => (
          <option key={index} value={option.value}>{option.label}</option>
        ))}
      </select>
    </div>
  );
}