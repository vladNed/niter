

/**
 * The build version of the wasm module.
 */
declare const wasmVersion;

/**
 * The function that initializes the wasm module.
 */
declare function wasmInit(config: string): Promise<void>;