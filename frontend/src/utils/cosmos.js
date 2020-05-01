// Most funcs here have been adapted from https://github.com/bluzelle/blzjs/blob/devel/src/swarmClient/cosmos.js
import shajs from 'sha.js';
import ripemd160 from 'ripemd160';
import ec from 'elliptic';
import { Buffer } from 'buffer/';
import BN from 'bn.js';
import bech32 from 'bech32';
import bitcoinjs from 'bitcoinjs-lib';
import bip32 from 'bip32';
import bip39 from 'bip39';

const BROADCAST_MAX_RETRIES = 10;
const BROADCAST_RETRY_INTERVAL_SECONDS = 1;
const PRIVATE_KEY_STORE_KEY = 'priv-key';
const HD_PATH = "m/44'/118'/0'/0/0"; // eslint-disable-line quotes

const secp256k1 = new ec.ec('secp256k1');

export default class {
  constructor({ host, chainId, gasInfo }) {
    this.host = host;
    this.chainId = chainId;
    this.gasInfo = gasInfo;
  }

  async loadPrivateKeyFromMnemonic(mnemonic) {
    const seed = await bip39.mnemonicToSeed(mnemonic);
    const node = await bip32.fromSeed(seed);
    const child = node.derivePath(HD_PATH);
    const ecpair = bitcoinjs.ECPair.fromPrivateKey(child.privateKey, {
      compressed: false,
    });
    this.privateKey = ecpair.privateKey.toString('hex');
    this.storePrivateKeyInCache(this.privateKey);
    return true;
  }

  loadPrivateKeyFromCache() {
    this.privateKey = cache(PRIVATE_KEY_STORE_KEY);
    return !!this.privateKey;
  }

  storePrivateKeyInCache(key) {
    cache(PRIVATE_KEY_STORE_KEY, key);
  }

  getPublicKey() {
    return secp256k1
      .keyFromPrivate(this.privateKey, 'hex')
      .getPublic(true, 'hex');
  }

  deriveAddress(prefix = 'cosmos') {
    const s = shajs('sha256')
      .update(Buffer.from(this.getPublicKey(), 'hex'))
      .digest();
    const r = new ripemd160().update(s).digest();
    const bytes = r;
    this.address = bech32.encode(prefix, bech32.toWords(bytes));
    return this.address;
  }

  disconnect() {
    this.storePrivateKeyInCache(null);
  }

  async loadAccount() {
    this.account = (await this.query(`/auth/accounts/${this.address}`)).value;
  }

  async query(endpoint) {
    return await this.xhr('get', endpoint);
  }

  async tx(method, endpoint, data) {
    this.broadcastRetries = 0;
    const tx = await this.validateTransaction(method, endpoint, data);
    return await this.broadcastSignedTransaction(tx);
  }

  async validateTransaction(method, endpoint, data) {
    const { value: tx } = await this.xhr(
      method,
      endpoint,
      sortJSON({
        ...this.generateBaseRequestPayload(),
        ...data,
      })
    );
    return tx;
  }

  generateBaseRequestPayload() {
    return {
      base_req: {
        chain_id: this.chainId,
        from: this.address,
      },
    };
  }

  async broadcastSignedTransaction(tx) {
    tx.memo = makeRandomString(32);

    tx.fee = {
      amount: [{ amount: this.gasInfo.minFee, denom: this.gasInfo.denom }],
      gas: tx.fee.gas,
    };

    // signature pub key
    const pubKeyValue = Buffer.from(this.getPublicKey(), 'hex').toString(
      'base64'
    );

    // signature value
    const payload = JSON.stringify({
      account_number: this.account.account_number.toString(),
      chain_id: this.chainId,
      fee: sortJSON(tx['fee']),
      memo: tx['memo'],
      msgs: sortJSON(tx['msg']),
      sequence: this.account.sequence.toString(),
    });
    const hash = shajs('sha256').update(Buffer.from(payload)).digest('hex');
    const sig = convertSignature(
      secp256k1.sign(hash, this.privateKey, 'hex', { canonical: true })
    ).toString('base64');

    tx.signatures = [
      {
        pub_key: {
          type: 'tendermint/PubKeySecp256k1',
          value: pubKeyValue,
        },
        signature: sig,
        account_number: this.account.account_number.toString(),
        sequence: this.account.sequence.toString(),
      },
    ];

    const response = await this.xhr('post', '/txs', {
      tx,
      mode: 'block',
    });
    const { data, raw_log } = response;
    if (!('code' in response)) {
      this.account.sequence += 1;
      if (!data) {
        return;
      }
      const hex = data.toString('hex');
      return hex ? JSON.parse(hex) : {};
    }

    if (~raw_log.indexOf('signature verification failed')) {
      this.broadcastRetries += 1;
      console.warn(
        'transaction failed ... retrying(%d) ...',
        this.broadcastRetries
      );
      if (this.broadcastRetries >= BROADCAST_MAX_RETRIES) {
        throw new Error('transaction failed after max retry attempts');
      }
      await sleep(BROADCAST_RETRY_INTERVAL_SECONDS * 1000);
      // lookup changed sequence
      await this.loadAccount();
      return await this.broadcastSignedTransaction(tx);
    }

    throw new Error(raw_log);
  }

  async xhr(method, endpoint, data) {
    const opts = {};
    if (data) {
      opts.method = method.toUpperCase();
      opts.body = JSON.stringify(data);
      opts.headers = {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      };
    }
    const res = await fetch(`${this.host}${endpoint}`, opts);
    const response = await res.json();
    return response.result || response; // todo
  }
}

function sortJSON(obj) {
  if (
    obj === null ||
    ~['undefined', 'string', 'number', 'boolean', 'function'].indexOf(
      typeof obj
    )
  ) {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.sort().map((i) => sortJSON(i));
  } else {
    const sortedObj = {};

    Object.keys(obj)
      .sort()
      .forEach((key) => {
        sortedObj[key] = sortJSON(obj[key]);
      });

    return sortedObj;
  }
}

function convertSignature(sig) {
  let r = new BN(sig.r);

  if (r.cmp(secp256k1.curve.n) >= 0) {
    r = new BN(0);
  }

  let s = new BN(sig.s);
  if (s.cmp(secp256k1.curve.n) >= 0) {
    s = new BN(0);
  }

  return Buffer.concat([
    r.toArrayLike(Buffer, 'be', 32),
    s.toArrayLike(Buffer, 'be', 32),
  ]);
}

function makeRandomString(length) {
  let result = '';
  const characters =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

async function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

function cache(k, v) {
  switch (arguments.length) {
    case 2:
      if (v === null) {
        return window.localStorage.removeItem(k);
      }
      return window.localStorage.setItem(k, JSON.stringify(v));

    case 1:
      try {
        return JSON.parse(window.localStorage.getItem(k));
      } catch (e) {
        return null;
      }

    default:
      return;
  }
}

export function toMicro(n) {
  return n * Math.pow(10, 6);
}

export function fromMicro(n) {
  return n / Math.pow(10, 6);
}
