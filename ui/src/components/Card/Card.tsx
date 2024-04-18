import type { PropsWithChildren } from 'react';
import { WithClassnameType } from 'types';

interface CardType extends PropsWithChildren, WithClassnameType {
  title: string;
  description?: string;
  anchor?: string;
}

export const Card = (props: CardType) => {
  const { title, children, description, anchor } = props;

  return (
    <div
      className='flex flex-col flex-1 rounded-xl bg-white p-6 justify-center'
      data-testid={props['data-testid']}
      id={anchor}
    >
      <h2 className='flex text-xl font-medium group'>
        {title}
      </h2>
      {description && <p className='text-gray-400 mb-6'>{description}</p>}
      {children}
    </div>
  );
};
