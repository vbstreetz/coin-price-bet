import { writable, get } from 'svelte/store';
import moment from 'moment';
import { query } from './xhr';

export const data = writable({
  times: [],
  prices: [],
});

export async function fetchData(coin = 0) {
  const { times, prices } = await query(
    `/coinpricebet/latest-coin-prices/${coin}`
  );
  data.set({
    times: times.map((s) => moment.unix(parseInt(s)).format('HH:mm:ss')),
    prices: prices.map((s) => parseInt(s) / 1000000),
  });
}
