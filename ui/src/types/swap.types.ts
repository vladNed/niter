import { type OfferDetails } from 'types';

type CreateSwapProps = {
  amount: string;
  claimProof: string;
  refundProof: string;
  callbackRoute: string;
};

type SwapStepProps = {
  offerData?: OfferDetails;
  progressHandler: () => void;
};

export type {
  CreateSwapProps,
  SwapStepProps,
};
