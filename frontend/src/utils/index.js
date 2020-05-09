export * from './sl';
export * from './xhr';
export * from './cache';
export * from './cosmos';

export function formatFiat(val, currency) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(
    val
  );
}

export function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

export function toMicro(n) {
  return n * Math.pow(10, 6);
}

export function fromMicro(n) {
  return n / Math.pow(10, 6);
}
