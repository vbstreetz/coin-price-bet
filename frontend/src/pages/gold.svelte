<script>
  import { onMount } from 'svelte';
  import {get} from 'svelte/store';
  import {coinPriceBetBlockchain} from '../utils/blockchains';
  import {address, info, balances, load, loadBalance, disconnectAccount, connectAccount} from '../stores/blockchains';
  import sl from '../utils/sl';

  onMount(load);

  async function onBuyGold() {
    const $info = get(info);
    try {
      await coinPriceBetBlockchain.tx('post', '/coinpricebet/buy', {
        amount: `1000000000transfer/${$info.betchainTransferChannel}/uatom`
      });
      sl('success', 'WAITING FOR CONFIRMATION...');
    } catch (e) {
      sl('error', e);
    }
  }
</script>

<style>

</style>

<div class="flex flex-grow flex-col">
  <div class="flex mb-5 items-center">
    <h1 class="main-heading flex-grow">GOLD CDP</h1>
    {#if $address}
      <button class="button is-light is-small" on:click={onBuyGold}>
        BUY GOLD
      </button>
      <button class="button is-light is-small ml-2" on:click={loadBalance}>
        REFRESH BALANCE
      </button>
      <button class="button is-light is-small ml-2" on:click={disconnectAccount}>
        DISCONNECT
      </button>
    {:else}
      <button class="button is-light is-small" on:click={connectAccount}>
        CONNECT
      </button>
    {/if}
  </div>

  {#if $address}
    <div class="mb-5">
      <div>Account: {$address}</div>
      <div>Balance:</div>
      <div class="ml-5">
        {#if $balances.coinPriceBet}
          {#each $balances.coinPriceBet as c}
            <div>{c.amount} {c.denom}</div>
          {/each}
        {/if}
      </div>
    </div>
  {/if}
</div>
