export const formatFiat = (val, currency) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(val);

export const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
