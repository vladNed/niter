import { DB_NAME, DB_VERSION } from 'config';
import Dexie, { Table } from 'dexie';

/* Wallet schema */
interface Wallet {
  label: string;
  wif: string;
}

export class NiterDB extends Dexie {
  wallets!: Table<Wallet, string>;

  constructor() {
    super(DB_NAME);
    this.version(DB_VERSION).stores({
      wallets: 'label, wif',
    });
  }

  async addWallet(wallet: Wallet): Promise<void> {
    try {
      await this.wallets.add(wallet);
    } catch (e) {
      switch (true) {
        case e instanceof Dexie.ConstraintError:
          throw new Error('Wallet already exists.');
        default:
          throw new Error('Failed to add wallet.');
      }
    }
  }

  async getWallet(label: string): Promise<Wallet | undefined> {
    return await this.wallets.get(label);
  }
}

export const db = new NiterDB();
export type { Wallet };