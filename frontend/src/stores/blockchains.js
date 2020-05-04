import { writable, get } from 'svelte/store';
import {
  coinPriceBetBlockchain,
  gaiaBlockchain,
  bandBlockchain,
} from '../utils/blockchains';

export const address = writable(null);
export const info = writable(null);
export const myInfo = writable(null);
export const balances = writable({});

export async function load() {
  await Promise.all([loadInfo(), loadAccount()]);
}

export async function connectAccount(mnemonic) {
  if (!mnemonic) {
    mnemonic = prompt('Enter mnemonic:');
  }
  if (!mnemonic) {
    return;
  }
  loadAccount(mnemonic);
}

export function disconnectAccount() {
  coinPriceBetBlockchain.disconnect();
  address.set(null);
}

async function loadInfo() {
  info.set(await coinPriceBetBlockchain.query('/coinpricebet/info'));
}

async function loadAccount(mnemonic) {
  if (mnemonic) {
    if (!(await coinPriceBetBlockchain.loadPrivateKeyFromMnemonic(mnemonic))) {
      return;
    }
  } else {
    if (!(await coinPriceBetBlockchain.loadPrivateKeyFromCache())) {
      return;
    }
  }

  address.set(coinPriceBetBlockchain.deriveAddress());

  gaiaBlockchain.loadPrivateKeyFromCache();
  gaiaBlockchain.deriveAddress();

  bandBlockchain.loadPrivateKeyFromCache();
  bandBlockchain.deriveAddress('band');

  await Promise.all([
    coinPriceBetBlockchain.loadAccount(),
    gaiaBlockchain.loadAccount(),
    // bandBlockchain.loadAccount()
  ]);

  Promise.all([loadBalance(), loadMyInfo()]);
}

export async function loadBalance() {
  balances.set({
    coinPriceBet: await coinPriceBetBlockchain.query(
      `/bank/balances/${get(address)}`
    ),
    gaia: await gaiaBlockchain.query(`/bank/balances/${get(address)}`),
    // band: await bandBlockchain.query(`/bank/balances/${bandBlockchain.address}`),
  });
}

async function loadMyInfo() {
  myInfo.set(
    await coinPriceBetBlockchain.query(`/coinpricebet/info/${get(address)}`)
  );
}
