import {
  AbiRegistry,
  Address,
  SmartContractTransactionsFactory,
  TransactionsFactoryConfig,
} from '@multiversx/sdk-core/out';
import { useCallback, useState } from 'react';
import swapRouterJson from 'contracts/swap-router.abi.json';
import {
  CREATE_SWAP_GAS_LIMIT,
  SWAP_ROUTER_CONTRACT_ADDRESS
} from 'config';
import {
  deleteTransactionToast,
  removeAllSignedTransactions,
  removeAllTransactionsToSign,
} from '@multiversx/sdk-dapp/services/transactions/clearTransactions';
import { useGetAccountInfo } from '@multiversx/sdk-dapp/hooks/account/useGetAccountInfo';
import { getChainId } from 'utils/getChainId';
import { signAndSendTransactions } from 'helpers';
import { type CreateSwapProps } from 'types';
import { useTrackTransactionStatus } from '@multiversx/sdk-dapp/hooks/transactions/useTrackTransactionStatus';

const SWAP_CREATE_INFO = {
  processingMessage: 'Locking funds in contract',
  errorMessage: 'An error occurred while creating the swap',
  successMessage: 'Fund locked successfully!',
}

export const useSwapRouterTransactions = () => {
  const { address } = useGetAccountInfo();
  const swapRouterAbi = AbiRegistry.create(swapRouterJson);
  const factoryConfig = new TransactionsFactoryConfig({ chainID: getChainId() })
  const factory = new SmartContractTransactionsFactory({ config: factoryConfig, abi: swapRouterAbi });

  const [swapRouterSessionId, setSwapRouterSessionId] = useState(sessionStorage.getItem('swapRouterSessionId') || '');
  const transactionStatus = useTrackTransactionStatus({
    transactionId: swapRouterSessionId ?? '0'
  });

  const clearAllTransactions = () => {
    removeAllSignedTransactions();
    removeAllTransactionsToSign();
    deleteTransactionToast(swapRouterSessionId ?? '0');
  }

  const sendCreateSwapTransaction = useCallback(async (props: CreateSwapProps) => {
    clearAllTransactions();

    const deploySwapTransaction = factory.createTransactionForExecute({
      sender: new Address(address),
      contract: Address.fromBech32(SWAP_ROUTER_CONTRACT_ADDRESS),
      function: 'createSwap',
      gasLimit: BigInt(CREATE_SWAP_GAS_LIMIT),
      arguments: [props.claimProof, props.refundProof],
      nativeTransferAmount: BigInt(props.amount),
    });

    const sessionId = await signAndSendTransactions({
      transactions: [deploySwapTransaction],
      callbackRoute: props.callbackRoute,
      transactionsDisplayInfo: SWAP_CREATE_INFO
    });

    sessionStorage.setItem('swapRouterSessionId', sessionId);
    setSwapRouterSessionId(sessionId);
  }, []);

  return {
    sendCreateSwapTransaction,
    transactionStatus
  };
};