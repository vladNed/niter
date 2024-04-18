import { AuthRedirectWrapper } from 'wrappers';
import {
  Account,
  Offers,
} from './widgets';
import { useScrollToElement } from 'hooks';
import { Widget } from './components';
import { WidgetType } from 'types/widget.types';

const WIDGETS: WidgetType[] = [
  {
    title: 'Account',
    widget: Account,
    description: 'Connected account details',
  },
  {
    title: 'Offers',
    widget: Offers,
    description: 'List of available offers',
  },
];

export const Dashboard = () => {
  useScrollToElement();

  return (
    <AuthRedirectWrapper>
      <div className='flex flex-col gap-6 max-w-3xl w-full'>
        {WIDGETS.map((element) => (
          <Widget key={element.title} {...element} />
        ))}
      </div>
    </AuthRedirectWrapper>
  );
};
