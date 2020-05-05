<script>
  import moment from 'moment';
  import {onMount} from 'svelte';
  import {get} from 'svelte/store';
  import Day from '../components/day.svelte';
  import {coinPriceBetBlockchain, gaiaBlockchain, bandBlockchain} from '../utils/blockchains';
  import {
    address,
    info,
    myInfo,
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
  import {fromMicro, toMicro, generateMnemonic} from '../utils/cosmos';
  import sl from '../utils/sl';
  import xhr from '../utils/xhr';
  import {sleep} from '../utils';
  import {COINS} from '../config';

  const tomorrow = moment.utc().add(1, 'days');
  const today = moment.utc();
  const yesterday = moment.utc().subtract(1, 'days');
  const days = [
    {label: 'Tomorrow', date: tomorrow},
    {label: 'Today', date: today},
    {label: 'Yesterday', date: yesterday}
  ];

  onMount(async function () {
    await load();
    // coinPriceBetBlockchain.setGasInfo({
    //   denom: `transfer/${$info.betchainTransferChannel}/uatom`,
    //   maxFee: toMicro(1)
    // });
  });
</script>

<div class="dark"> <!-- todo: find out why purgecss doesn't look up html.dark -->
  <div class="flex">
    <h1 class="main-heading flex-grow">BET TODAY,<br/>THE BEST CRYPTO OF TOMORROW, AND WIN!</h1>

    {#if $address}
      <div class="flex flex-col text-sm">
        <div class='mr-3'>Account: {$address}</div>
        <table class='balances'>
          {#if $parsedBalances.coinPriceBet}
          <tr><td>Balance:</td><td>{fromMicro($parsedBalances.coinPriceBet.atom || 0)}atom</td></tr>
          {/if}
          {#if $myInfo}
          <tr><td>Total Bets:</td><td>{fromMicro($myInfo.totalBetsAmount)}atom</td></tr>
          <tr><td>Total Wins:</td><td>{fromMicro($myInfo.totalWinsAmount)}atom</td></tr>
          {/if}
        </table>
      </div>
      <button class="button is-light is-small ml-2" on:click={loadBalance}>
        REFRESH BALANCE
      </button>
      <button class="button is-primary is-small ml-2" on:click={disconnectAccount}>
        DISCONNECT
      </button>
    {:else}
      <button class="button is-primary is-small mr-3" on:click={generateAccount}>
        GENERATE
      </button>
      <button class="button is-primary is-small" on:click={connectAccount}>
        CONNECT
      </button>
    {/if}
  </div>

  {#if info}
    <div class="main-container">
      {#each days as day}
        {#if $info}
        <Day day={day.date} label={day.label} />
        {/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .main-heading {
    font-weight: bolder;
    margin-left: 50px;
  }

  .main-container {
    margin-top: 50px;
    padding: 30px 16px 16px 0;
    position: relative;
  }

  .main-container:before {
    right: auto;
    left: 47px;
    background: var(--border-color);
    bottom: 0;
    content: "";
    height: 100%;
    position: absolute;
    top: 0;
    width: 1px;
  }

  .balances {
    width: 140px;
  }
</style>
