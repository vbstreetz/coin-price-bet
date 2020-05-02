export const formatFiat = (val, currency) =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(val);
