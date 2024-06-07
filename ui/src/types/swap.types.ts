import { type OfferDetails } from 'types';

export type CreateSwapProps = {
  amount: string;
  claimProof: string;
  refundProof: string;
  callbackRoute: string;
};

export type SwapStepProps = {
  offerData?: OfferDetails;
};
