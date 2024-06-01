export type WidgetProps = {
  callbackRoute: string;
};

export type WidgetType<T = any> = {
  title: string;
  widget: (props: T) => JSX.Element;
  description?: string;
  props?: { receiver?: string };
  anchor?: string;
};

export interface SwapFieldProps {
  icon: React.ReactElement;
  ticker: string;
  name: string;
  side: 'Swap' | 'For';
  value: string;
  dataSide: 'sending' | 'receiving';
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export interface SideToken {
  ticker: string;
  name: string;
  icon: React.ReactElement;
}

interface SearchOfferWidgetProps {
  isPlaceholder: boolean;
}

export interface DrawerProps {
  isOpen: boolean;
  toggleDrawer: () => void;
}

export type {
  SearchOfferWidgetProps,
};
