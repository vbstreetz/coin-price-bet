// Most funcs here have been adapted from https://github.com/bluzelle/blzjs/blob/devel/src/swarmClient/cosmos.js

export async function signTransaction() {
  let payload = {
    account_number: account_info.account_number || '0',
    chain_id: chain_id,
    fee: util.sortJson(data.value.fee),
    memo: data.value.memo,
    msgs: util.sortJson(data.value.msg),
    sequence: (account_info.sequence || '0')
  };

  // Calculate the SHA256 of the payload object
  let jstr = JSON.stringify(payload);
  let sstr = sanitize_string(jstr);
  let jsonHash = util.hash('sha256', Buffer.from(sstr));

  return {
    pub_key: {
      type: 'tendermint/PubKeySecp256k1',
      value: Buffer.from(
          secp256k1
              .keyFromPrivate(key, 'hex')
              .getPublic(true, 'hex'),
          'hex'
      ).toString('base64'),
    },

    // We have to convert the signature to the format that Tendermint uses
    signature: util.convertSignature(
        secp256k1.sign(jsonHash, key, 'hex', {
          canonical: true,
        }),
    ).toString('base64'),

    account_number: account_info.account_number,
    sequence: account_info.sequence
  }
}

async function getECPrivateKey(mnemonic) {
  const seed = await bip39.mnemonicToSeed(mnemonic);
  const node = await bip32.fromSeed(seed);
  const child = node.derivePath(path);
  const ecpair = bitcoinjs.ECPair.fromPrivateKey(child.privateKey, {compressed: false});
  return ecpair.privateKey.toString('hex');
}
