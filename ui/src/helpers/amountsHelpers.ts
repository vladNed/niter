import { BTC_DECIMALS, CoinCurrency, DECIMALS } from "localConstants";

export const formatCoinAmount = (amount: string, currency: CoinCurrency): number => {
  const value = parseFloat(amount)
  switch (currency) {
    case CoinCurrency.BTC:
      return value / Math.pow(10, BTC_DECIMALS)
    case CoinCurrency.EGLD:
      return value / Math.pow(10, DECIMALS)
    default:
      return value
  }
}
