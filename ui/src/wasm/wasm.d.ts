

interface PeerInfo {
  id: string;
  state: string;
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
declare function wasmPollOffers(): Promise<string[]>;

/**
 * Get the peer state
 */
declare function wasmGetPeerState(): PeerInfo;
