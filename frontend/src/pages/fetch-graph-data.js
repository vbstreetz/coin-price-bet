import {writable, get} from "svelte/store";
import moment from 'moment';

export const data = writable({
  times: [],
  prices: []
});

export async function fetchData(coin = 0) {
  const {result: {times, prices}} = await (await fetch(`http://localhost:1317/coinpricebet/latest-coin-prices/${coin}`)).json();
  data.set({times: times.map(s => moment.unix(parseInt(s)).format('hh:mm:ss')), prices: prices.map(s => parseInt(s) / 1000000)});
}
