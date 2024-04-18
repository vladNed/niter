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

export type OfferData = {
  receivingCoin: Coin;
  givingCoin: Coin;
  receivingAmount: string;
  givingAmount: string;
}