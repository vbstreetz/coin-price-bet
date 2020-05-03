<script>
  import moment from 'moment';
  import {onMount} from 'svelte';
  import {derived, get} from 'svelte/store';
  import Day from '../components/day.svelte';
  import {coinPriceBetBlockchain, gaiaBlockchain, bandBlockchain} from '../utils/blockchains';
  import {address, info, myInfo, balances as allBalances, load, disconnectAccount, connectAccount} from '../stores/blockchains';
  import {fromMicro, toMicro} from '../utils/cosmos';
  import sl from '../utils/sl';
  import xhr from '../utils/xhr';
  import {COINS} from '../config';

  const tomorrow = moment.utc().add(1, 'days');
  const today = moment.utc();
  const yesterday = moment.utc().subtract(1, 'days');
  const days = [
    {label: 'Tomorrow', date: tomorrow},
    {label: 'Today', date: today},
    {label: 'Yesterday', date: yesterday}
  ];
  const balances = derived(allBalances, function({coinPriceBet, gaia}) {
    const $info = get(info);
    const b = {}

    if (coinPriceBet) {
      for (let i = 0; i < coinPriceBet.length; i++) {
        const coin = coinPriceBet[i];
        if (~coin.denom.search($info.betchainTransferChannel)) {
          b.coinPriceBet = coin.amount;
          break;
        }
      }
    }

    if (gaia) {
      for (let i = 0; i < gaia.length; i++) {
        const coin = gaia[i];
        if (coin.denom === 'uatom') {
          b.gaia = coin.amount;
          break;
        }
      }
    }

    return b;
  });

  onMount(load);

  async function rechargeFromFaucet() {
    await xhr('post', '/gaia-faucet/', {
      "address": get(address),
      "chain-id": "band-cosmoshub"
    });
    sl('success', 'Waiting for confirmation...');
  }

  async function rechargeFromGaia() {
    // const amount = parseInt(prompt('Amount?'));
    // if (!amount) return sl('error', 'Invalid amount');

    const amount = 10;

    const $address = get(address);
    const $info = get(info);
    // const channel = $info.gaiaTransferChannel;
    const channel = $info.betchainTransferChannel;
    const port = 'transfer';

    try {
      await gaiaBlockchain.tx('post', `/ibc/ports/${port}/channels/${channel}/transfer`, {
        "amount": [
          {
            "denom": "uatom",
            // "amount": `${toMicro(amount).toString()}transfer/${$info.betchainTransferChannel}/uatom`
            "amount": `${toMicro(amount).toString()}`
          }
        ],
        "receiver": get(address),
        // "source": true
      });
      sl('success', 'WAITING FOR CONFIRMATION...');
    } catch (e) {
      sl('error', e);
    }
  }
</script>

<div>
  <div class="flex">
    <h1 class="main-heading flex-grow">BET TODAY,<br/>THE BEST CRYPTO OF TOMORROW, AND WIN!</h1>

    {#if $address}
      <div class="flex flex-col text-sm">
        <div>Account: {$address}</div>
        <div>Balance: {fromMicro($balances.gaia || 0)}atom (gaia) <span class="cursor-pointer underline" on:click={rechargeFromFaucet}>recharge from faucet</span></div>
        <div>Balance: {fromMicro($balances.coinPriceBet || 0)}atom (coinpricebet) <span class="cursor-pointer underline" on:click={rechargeFromGaia}>recharge from gaia</span></div>
        {#if $myInfo}
        <div>Total Bets: {fromMicro($myInfo.totalBetsAmount)}atom</div>
        <div>Total Wins: {fromMicro($myInfo.totalWinsAmount)}atom</div>
        {/if}
      </div>
      <button class="button is-primary is-small ml-2" on:click={disconnectAccount}>
        DISCONNECT
      </button>
    {:else}
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
</style>
