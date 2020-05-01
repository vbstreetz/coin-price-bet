<script>
  import {onMount} from 'svelte';
  import account from '../utils/account';
  import Chart from 'svelte-frappe-charts';
  import { writable } from 'svelte/store';
  import moment from 'moment';

  const tokens = [
    "BTC",
    "ETH",
    "LTC",
    "BAND",
    "ATOM",
    "LINK",
    "XTZ"
  ];
  let currentToken = tokens[0];
  export const data = writable({
    times: [],
    prices: [],
  });

  $: graph = {
    labels: $data.times,
    datasets: [{values: $data.prices}]
  };

  onMount(fetchData);

  function switchChart(event) {
    const token = event.target.dataset.token;
    currentToken = tokens[token];
    fetchData(token);
  }

  async function fetchData(coin = 0) {
    const { times, prices } = await account.query(`/coinpricebet/latest-coin-prices/${coin}`);
    data.set({
      times: times.map((s) => moment.unix(parseInt(s)).format('HH:mm:ss')),
      prices: prices.map((s) => parseInt(s) / 1000000),
    });
  }
</script>

<style>
  .heading {
    margin-bottom: 30px;
  }
</style>

<div class="flex flex-col">
  <div class="heading">
    <h3>Price Performance</h3>
  </div>

  <div>
    <div class="dropdown is-hoverable">
      <div class="dropdown-trigger cursor-pointer">
        <div
          aria-haspopup="true"
          aria-controls="dropdown-menu"
          class="flex"
        >
          <span>{currentToken}</span>
          <span class="icon is-small dropdown-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="12"
              height="12"
              viewBox="0 0 512 512"
            >
              <title>ionicons-v5-a</title>
              <polyline
                points="112 184 256 328 400 184"
                style="fill:none;stroke:#000;stroke-linecap:round;stroke-linejoin:round;stroke-width:48px"
              />
            </svg>
          </span>
        </div>

      </div>
      <div class="dropdown-menu" id="draw-dropdown-menu" role="menu">
        <div class="dropdown-content">
          {#each tokens as token, id}
            <a class="dropdown-item" href="javascript:" data-token={id} on:click={switchChart}>
              {token}
            </a>
          {/each}
        </div>
      </div>
    </div>
  </div>

  {#if $data.times.length}
    <Chart data={graph} type="line"/>
  {/if}
</div>
