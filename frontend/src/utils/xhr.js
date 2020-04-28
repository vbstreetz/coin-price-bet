import {signTransaction} from './crypto';

const API_HOST = window.location.hostname; // window.location.hostname === 'localhost' ? 'localhost': '144.202.100.245';

export async function query(endpoint) {
  const res = await fetch(buildUrl(endpoint));
  const {result} = await res.json();
  return result;
}

export async function tx(endpoint, txn) {
  const signedTransaction = signTransaction(txn);
  return await broadcastSignedTransaction(endpoint, signedTransaction);
}

async function broadcastSignedTransaction(endpoint, txn) {
  const res = await fetch(buildUrl(endpoint), {method: 'POST', body: data});
  const {data, code, raw_log} = await res.json();
  if (code) {
    throw new Error(raw_log);
  }
  const hex = data?.toString('hex');
  return hex ? JSON.parse(hex) : {};
}

function buildUrl(path) {
  return `http://${API_HOST}:1317/coinpricebet/${path}`;
}