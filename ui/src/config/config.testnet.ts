import { EnvironmentsEnum } from 'types';

export * from './sharedConfig';

export const contractAddress =
  'erd1qqqqqqqqqqqqqpgqq3qsdxf55rlz5ka8mw3jdnacm8dlkuy09l5ql0wrlm';
export const API_URL = 'https://testnet-gateway.multiversx.com';
export const sampleAuthenticatedDomains = [API_URL];
export const environment = EnvironmentsEnum.testnet;
export const WASM_LOG_LEVEL = 0;
export const SIGNALLING_SERVER_URL = 'ws://192.168.1.129:8080/ws/v1/';
