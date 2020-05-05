import { writable, get, derived } from 'svelte/store';
import {
  coinPriceBetBlockchain,
  gaiaBlockchain,
  bandBlockchain,
} from '../utils/blockchains';
import { toMicro, generateMnemonic } from '../utils/cosmos';
import { sleep } from '../utils';
import sl from '../utils/sl';
import xhr from '../utils/xhr';

export const address = writable(null);
export const info = writable(null);
export const myInfo = writable(null);
export const balances = writable({});

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
  coinPriceBetBlockchain.disconnect();
  address.set(null);
}

async function loadInfo() {
  info.set(await coinPriceBetBlockchain.query('/coinpricebet/info'));
}

export async function loadAccount(mnemonic) {
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
  const [coinPriceBet, gaia, band] = await Promise.all([
    coinPriceBetBlockchain.query(`/bank/balances/${get(address)}`),
    gaiaBlockchain.query(`/bank/balances/${get(address)}`),
    // bandBlockchain.query(`/bank/balances/${bandBlockchain.address}`)
  ]);
  balances.set({ coinPriceBet, gaia, band });
}

async function loadMyInfo() {
  myInfo.set(
    await coinPriceBetBlockchain.query(`/coinpricebet/info/${get(address)}`)
  );
}

export const parsedBalances = derived([balances, info], function ([
  $balances,
  $info,
]) {
  const b = {
    coinPriceBet: {},
    gaia: {},
  };

  const { coinPriceBet, gaia } = $balances;
  if (coinPriceBet) {
    for (let i = 0; i < coinPriceBet.length; i++) {
      const coin = coinPriceBet[i];
      if (~coin.denom.search($info.betchainTransferChannel)) {
        b.coinPriceBet.atom = coin.amount;
      }
      if (coin.denom === 'stake') {
        b.coinPriceBet.stake = coin.amount;
      }
    }
  }

  if (gaia) {
    for (let i = 0; i < gaia.length; i++) {
      const coin = gaia[i];
      if (coin.denom === 'uatom') {
        b.gaia.atom = coin.amount;
        break;
      }
    }
  }

  return b;
});

export async function rechargeAtomFromFaucet() {
  await xhr('post', '/gaia-faucet/', {
    address: get(address),
    'chain-id': 'band-cosmoshub',
  });
  sl('success', 'Waiting for confirmation...');
  setTimeout(() => loadBalance(), 2000);
}

export async function rechargeStakeFromFaucet() {
  await xhr('get', `/vb-rest/coinpricebet/faucet/${get(address)}`);
  sl('success', 'Waiting for confirmation...');
  await sleep(3000);
  await loadBalance();
}

export async function rechargeAtomFromGaia() {
  const amount = parseInt(prompt('Amount?'));
  if (!amount) {
    return sl('error', 'Invalid amount');
  }

  // const $address = get(address);
  const $info = get(info);
  const channel = $info.gaiaTransferChannel;
  // const channel = $info.betchainTransferChannel;
  const port = 'transfer';

  try {
    await gaiaBlockchain.tx(
      'post',
      `/ibc/ports/${port}/channels/${channel}/transfer`,
      {
        amount: [
          {
            denom: `transfer/${$info.betchainTransferChannel}/uatom`,
            amount: toMicro(amount).toString(),
          },
        ],
        receiver: get(address),
        source: true,
      }
    );
    sl('success', 'WAITING FOR CONFIRMATION...');
    await sleep(3000);
    await loadBalance();
  } catch (e) {
    sl('error', e);
  }
}

export function generateAccount() {
  const mnemonic = generateMnemonic();
  console.log(mnemonic);
  loadAccount(mnemonic);
}
