type CreateOfferProps = {
  handleReceiptShow: (offerData: OfferDetails) => void
}

type OfferDetails = {
  swapCreator: string;
  sendingAmount: string;
  sendingCurrency: string;
  receivingAmount: string;
  receivingCurrency: string;
  createdAt?: string;
  expiresAt?: string;
  isSwapCreator?: boolean;
}

export type { OfferDetails, CreateOfferProps }