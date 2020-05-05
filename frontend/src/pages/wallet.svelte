<script>
  import { onMount } from 'svelte';
  import {get} from 'svelte/store';
  import {coinPriceBetBlockchain} from '../utils/blockchains';
  import {
    address,
    info,
    balances,
    parsedBalances,
    load,
    loadBalance,
    disconnectAccount,
    connectAccount,
    generateAccount,
    rechargeAtomFromFaucet,
    rechargeAtomFromGaia,
    rechargeStakeFromFaucet,
  } from '../stores/blockchains';
  import sl from '../utils/sl';
  import {fromMicro} from '../utils/cosmos';

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
    <h1 class="main-heading flex-grow">WALLET</h1>
    {#if $address}
      <button class="button is-light is-small ml-2" on:click={loadBalance}>
        REFRESH BALANCE
      </button>
      <button class="button is-light is-small ml-2" on:click={disconnectAccount}>
        DISCONNECT
      </button>
    {:else}
      <button class="button is-light is-small mr-3" on:click={generateAccount}>
        GENERATE
      </button>
      <button class="button is-light is-small" on:click={connectAccount}>
        CONNECT
      </button>
    {/if}
  </div>

  {#if $address}
    <div class="mb-5">
      <div>Account: {$address}</div>
      <div>Balances:</div>
      <div class="ml-5">
        {#if $balances.coinPriceBet}
          ---
          {#each $balances.coinPriceBet as c}
            <div>{c.amount} {c.denom}</div>
          {/each}
          {#each $balances.gaia as c}
            <div>{c.amount} {c.denom}</div>
          {/each}
        {/if}

        {#if $address}
          ---
          <div class="flex flex-col text-sm">
            <div>{fromMicro($parsedBalances.gaia.atom || 0)}atom (gaia) <span class="cursor-pointer underline" on:click={() => rechargeAtomFromFaucet()}>recharge from faucet</span></div>
            <div>{fromMicro($parsedBalances.coinPriceBet.atom || 0)}atom (coinpricebet) <span class="cursor-pointer underline" on:click={() => rechargeAtomFromGaia()}>recharge from your gaia account</span></div>
            <div>{fromMicro($parsedBalances.coinPriceBet.stake || 0)}stake (coinpricebet) <span class="cursor-pointer underline" on:click={() => rechargeStakeFromFaucet()}>request from faucet (used for transactions fee)</span></div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>
