import NProgress from 'nprogress';
import Cosmos, { toMicro } from './cosmos';
import { API_HOST } from '../config';

export const coinPriceBetBlockchain = new (class extends Cosmos {
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
  host: API_HOST + '/vb-rest',
  chainId: 'band-consumer',
  // gasInfo: { maxFee: toMicro(1), denom: 'stake' },
});

export const gaiaBlockchain = new (class extends Cosmos {
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
  host: API_HOST + '/gaia-rest',
  chainId: 'band-cosmoshub',
  gasInfo: { maxFee: toMicro(1), denom: 'uatom' },
});

export const bandBlockchain = new (class extends Cosmos {
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
  host: API_HOST + '/band-rest',
  chainId: 'ibc-bandchain',
  gasInfo: { maxFee: toMicro(1), denom: 'uband' },
});
