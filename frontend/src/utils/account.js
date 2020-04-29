import { getECPrivateKey, getAddress, signTransaction } from './crypto';
import { query, mutate } from './xhr';

export default new (class {
  constructor() {}

  async setMnemonic(mnemonic) {
    this.privateKey = getECPrivateKey(mnemonic);
    this.address = getAddress(this.privateKey);
    this.account = (await query(`/auth/accounts/${this.address}`)).value;
  }

  async query(endpoint, data) {
    return query(endpoint, data);
  }

  async tx(endpoint, data) {
    return signTransaction({
      accountNumber: this.account.account_number,
      accountSequence: this.account.sequence,
      privateKey: this.privateKey,
    });
    // const signedTransaction = signTransaction(data);
    // return await this.broadcastSignedTransaction(endpoint, signedTransaction);
  }

  async broadcastSignedTransaction(endpoint, txn) {
    const res = await mutate(endpoint, txn);
    const { data, code, raw_log } = await res.json();
    if (code) {
      throw new Error(raw_log);
    }
    if (!data) {
      return;
    }
    const hex = data.toString('hex');
    return hex ? JSON.parse(hex) : {};
  }
})();
