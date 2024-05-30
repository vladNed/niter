

interface PeerInfo {
  id: string;
  state: string;
  remotePeer: string;
}

interface OfferData {
  id: string;
}


declare type PeerData = {
  side: string,
  data: string,
  timestamp: string,
}


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
declare function wasmCreateOffer(): Promise<string>;

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
 * The function that polls exchange data.
 */
declare function wasmPollExchangeData(): Promise<PeerData[]>;

/**
 * Sends data to the peer.
 */
declare function wasmSendData(data: string): Promise<void>;