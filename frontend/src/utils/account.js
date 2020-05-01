import NProgress from 'nprogress';
import Cosmos from './cosmos';

const API_HOST =
  window.location.hostname === 'localhost'
    ? 'http://localhost:1317'
    : 'http://144.202.100.245:1317';
const CHAIN_ID = 'vbstreetz';

export default new (class extends Cosmos {
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
  host: API_HOST,
  chainId: CHAIN_ID,
  gasInfo: { minFee: '525000', denom: 'stake' },
});
