import { EnvironmentsEnum } from 'types';

export * from './sharedConfig';

export const API_URL = 'https://testnet-gateway.multiversx.com';
export const sampleAuthenticatedDomains = [API_URL];
export const environment = EnvironmentsEnum.testnet;
export const WASM_LOG_LEVEL = 0;
export const SIGNALLING_SERVER_URL = 'ws://192.168.1.129:8080/ws/v1/';
export const SWAP_ROUTER_CONTRACT_ADDRESS = 'erd1qqqqqqqqqqqqqpgqmknzxse8fktzxz6eyg54h3mvl7wug85p7pyqm2hn7y';
export const CREATE_SWAP_GAS_LIMIT = 10000000;
