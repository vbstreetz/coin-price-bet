import NProgress from "nprogress";
import {writable, get} from 'svelte/store';
import {sleep} from '../utils';
import xhr from '../utils/xhr';
import Cosmos, {toMicro, generateMnemonic} from '../utils/cosmos';
import {API_HOST} from "../config";

export const blockchain = new (class extends Cosmos {
  async xhr(...args) {
    NProgress.start();
    NProgress.set(0.4);
    try {
      return await super.xhr(...args);
    } finally {
      NProgress.done();
    }
  }
})({
  host: API_HOST + '/bet-rest',
  chainId: 'band-consumer',
  gasInfo: {maxFee: toMicro(1), denom: 'stake'},
});

export const address = writable(null);
export const info = writable(null);
export const myInfo = writable(null);
export const balance = writable({});

export async function load() {
  await Promise.all([loadInfo(), loadAccount()]);
}

export async function connectAccount() {
  const mnemonic = prompt('Enter mnemonic:');
  if (!mnemonic) {
    return;
  }
  loadAccount(mnemonic);
}

export function disconnectAccount() {
  blockchain.disconnect();
  address.set(null);
}

async function loadInfo() {
  info.set(await blockchain.query('/coinpricebet/info'));
}

export async function loadAccount(mnemonic) {
  if (mnemonic) {
    if (!(await blockchain.loadPrivateKeyFromMnemonic(mnemonic))) {
      return;
    }
  } else {
    if (!(await blockchain.loadPrivateKeyFromCache())) {
      return;
    }
  }

  address.set(blockchain.deriveAddress());

  await blockchain.loadAccount();

  Promise.all([loadBalance(), loadMyInfo()]);
}

export async function loadBalance() {
  const balances = await blockchain.query(`/bank/balances/${get(address)}`);
  for (let i = 0; i < balances.length; i++) {
    const coin = balances[i];
    if (~coin.denom.search('stake')) {
      return balance.set(parseInt(coin.amount));
    }
  }
  balance.set(0);
}

export async function tryReloadBalance() {
  let b = get(balance);
  for (let i = 0; i < 3; i++) {
    if (b !== get(balance)) break;
    b = get(balance);
    await sleep(2000);
    await loadBalance();
  }
}

async function loadMyInfo() {
  myInfo.set(
    await blockchain.query(`/coinpricebet/info/${get(address)}`)
  );
}

export async function rechargeFromFaucet() {
  await xhr('get', `/bet-rest/coinpricebet/faucet/${get(address)}`);
}

export async function generateAccount() {
  const mnemonic = generateMnemonic();
  console.log(mnemonic);
  await loadAccount(mnemonic);
  await rechargeFromFaucet();
  await tryReloadBalance();
}
