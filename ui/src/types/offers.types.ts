export enum Coin {
  EGLD = 'EGLD',
  MNR = 'MNR'
}

export enum TradeType {
  Buy = 'buy',
  Sell = 'sell'
}

export type TradeTypeOptions = {
  label: string;
  value: string;
}

interface OfferDetails {
  swapCreator: string;
  sendingAmount: string;
  sendingCurrency: string;
  receivingAmount: string;
  receivingCurrency: string;
  createdAt?: string;
  expiresAt?: string;
}

export type { OfferDetails }