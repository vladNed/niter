type CreateOfferProps = {
  handleReceiptShow: (offerData: OfferDetails) => void
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

export type { OfferDetails, CreateOfferProps }