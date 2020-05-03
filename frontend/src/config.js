export const COINS = [
  "BTC",
  "ETH",
  "LTC",
  "BAND",
  "ATOM",
  "LINK",
  "XTZ",
];

export const DAY_STATES = {
  BET: 0,
  DRAWING: 1,
  PAYOUT: 2
};

export const API_HOST =
  window.location.hostname === 'localhost'
    ? 'http://localhost:1317'
    : 'https://witnet.tools';
