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
}

export interface SideToken {
  ticker: string;
  name: string;
  icon: React.ReactElement;
}

export interface SearchOfferWidgetProps {
  isPlaceholder: boolean;
}
