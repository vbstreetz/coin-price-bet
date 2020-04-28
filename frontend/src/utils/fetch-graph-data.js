import {writable, get} from "svelte/store";
import moment from 'moment';

const API_HOST = window.location.hostname; // window.location.hostname === 'localhost' ? 'localhost': '144.202.100.245';

export const data = writable({
  times: [],
  prices: []
});

export async function fetchData(coin = 0) {
  const {result: {times, prices}} = await (await fetch(`http://${API_HOST}:1317/coinpricebet/latest-coin-prices/${coin}`)).json();
  data.set({times: times.map(s => moment.unix(parseInt(s)).format('HH:mm:ss')), prices: prices.map(s => parseInt(s) / 1000000)});
}
