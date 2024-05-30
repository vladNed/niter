import { RouteNamesEnum } from 'localConstants';
import { Swap, Explore, Home, Connect } from 'pages';
import { RouteType } from 'types';

interface RouteWithTitleType extends RouteType {
  title: string;
}

export const routes: RouteWithTitleType[] = [
  {
    path: RouteNamesEnum.home,
    title: 'Home',
    component: Home
  },
  {
    path: RouteNamesEnum.swap,
    title: 'Swap',
    component: Swap
  },
  {
    path: RouteNamesEnum.explore,
    title: 'Explore',
    component: Explore
  },
  {
    path: RouteNamesEnum.connect,
    title: 'Connect',
    component: Connect
  }
];
