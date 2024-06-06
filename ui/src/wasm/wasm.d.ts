

interface PeerInfo {
  id: string;
  state: string;
  remotePeer: string;
};

interface OfferData {
  id: string;
  receivingAmount: string;
  receivingCurrency: string;
  sendingAmount: string;
  sendingCurrency: string;
  swapCreator: string;
};


declare type PeerData = {
  side: string,
  data: string,
  timestamp: string,
};


/**
 * The build version of the wasm module.
 */
declare const wasmVersion;

/**
 * The function that initializes the wasm module.
 */
declare function wasmInit(config: string): Promise<void>;

/**
 * The function that creates an offer and SDP for the peer.
 */
declare function wasmCreateOffer(offer: string): Promise<string>;

/**
 * The function that creates an answer and SDP for the peer.
 */
declare function wasmCreateAnswer(offerId: string): Promise<void>;

/**
 * The function that polls the broadcasted offers.
 */
declare function wasmPollOffers(): Promise<OfferData[]>;

/**
 * Get the peer state
 */
declare function wasmGetPeerState(): PeerInfo;

/**
 * The function that connects to the peer.
 */
declare function wasmMakeConnect(): Promise<string>;

/**
 * The function that sends the SDP to the peer.
 */
declare function wasmConnect(sdp: string): Promise<string>;


/**
 * Sends data to the peer.
 */
declare function wasmSendData(data: string): Promise<void>;


declare function wasmInitWallet(wif?: string, mvxAddress?: string): Promise<string>;

/*
  * Gets all the events that happened in a swap context
  */
declare function wasmGetSwapEvents(): string[];


/**
 * The function that retrieves the data from the swap state to create a transaction.
 */
declare function wasmTransactionRequest(transactionType: string): Promise<any>;