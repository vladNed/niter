import { RouteNamesEnum } from 'localConstants';
import { Swap, Home } from 'pages';
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
];
