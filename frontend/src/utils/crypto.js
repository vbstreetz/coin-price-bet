// Most funcs here have been adapted from https://github.com/bluzelle/blzjs/blob/devel/src/swarmClient/cosmos.js
// import { Buffer } from 'buffer/';
import { generateWalletFromMnemonic } from '@everett-protocol/cosmosjs/utils/key';
// import { sortJSON } from '@everett-protocol/cosmosjs/utils/sortJson';
import { CHAIN_ID, HD_PATH } from '../config';
import secp256k1 from 'secp256k1';
import elliptic from 'elliptic';

export function getECPrivateKey(mnemonic) {
  return generateWalletFromMnemonic(mnemonic, HD_PATH);
}

export function getAddress(privateKey, prefix) {
  return privateKey.toPubKey().toAddress().toBech32(prefix);
}

export function signTransaction({
  accountNumber,
  accountSequence,
  memo,
  fee,
  msg,
  privateKey,
}) {
  if (!memo) {
    const s = {
      account_number: accountNumber.toString(),
      chain_id: CHAIN_ID,
      fee: { amount: [], gas: '200000' },
      memo: '',
      msgs: [
        {
          type: 'coin_price_bet/BuyGold',
          value: {
            amount: [{ denom: 'transfer//uatom', amount: '1000000000' }],
            buyer: 'cosmos15d4apf20449ajvwycq8ruaypt7v6d34522frnd',
          },
        },
      ],
      sequence: accountSequence.toString(),
    };

    const sig = privateKey.sign(JSON.stringify(sortJSON(s)));
    console.log(sig.toString('base64'));
    console.log(privateKey);
    console.log(
      secp256k1
        .publicKeyConvert(privateKey.toPubKey().pubKey, false)
        .toString('hex')
    );

    const sd = new elliptic.ec('secp256k1');

    console.log(
      Buffer.from(
        sd.keyFromPrivate(privateKey.privKey).getPublic(true, 'hex'),
        'hex'
      ).toString('hex')
    );
    console.log(
      'eb5ae9872102db19bba04424d63f05a1765f5efb39d2b3f550e3b4152535b47c1a12e5a886f6'
    );
    return;
  }
  const payload = {
    account_number: accountNumber.toString(),
    chain_id: CHAIN_ID,
    fee: sortJSON(fee),
    memo,
    msgs: sortJSON(msg),
    sequence: accountSequence.toString(),
  };

  const signature = privateKey.sign(JSON.stringify(payload));

  return {
    pub_key: {
      type: 'tendermint/PubKeySecp256k1',
      value: privateKey.toPubKey().toString('base64'),
    },
    signature: signature.toString('base64'),
    account_number: accountNumber.toString(),
    sequence: accountSequence.toString(),
  };
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
