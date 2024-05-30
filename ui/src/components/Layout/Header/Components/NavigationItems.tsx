import { MxLink } from 'components/MxLink';
import { environment } from 'config';

interface NavItem {
  title: string;
  path: string;
}

interface NavProps {
  items: NavItem[];
}

export const NavigationItems = (props: NavProps) => {
  return (
    <nav className='flex gap-4 text-zinc-500 place-items-center content-center'>
      {props.items.map((item, index) => (
        <MxLink to={item.path} key={index} className='hover:bg-zinc-800 hover:text-zinc-400 rounded-xl ease-in transition duration-100 py-2 px-4'>
          {item.title}
        </MxLink>
      ))}

    </nav>
  )
}