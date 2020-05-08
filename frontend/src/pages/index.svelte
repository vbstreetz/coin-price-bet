<script>
  import moment from 'moment';
  import {onMount} from 'svelte';
  import Day from '../components/day.svelte';
  import {
    address,
    info,
    myInfo,
    balance,
    load,
    loadBalance,
    tryReloadBalance,
    disconnectAccount,
    connectAccount,
    generateAccount,
    rechargeFromFaucet,
  } from '../stores/blockchain';
  import {fromMicro, toMicro, generateMnemonic} from '../utils/cosmos';
  import sl from '../utils/sl';
  import {sleep} from '../utils';

  const tomorrow = moment.utc().add(1, 'days');
  const today = moment.utc();
  const yesterday = moment.utc().subtract(1, 'days');
  const days = [
    {label: 'Tomorrow', date: tomorrow},
    {label: 'Today', date: today},
    {label: 'Yesterday', date: yesterday}
  ];
  let loaded = false;

  onMount(async function () {
    await load();
    loaded = true;
  });

  async function onRechargeFromFaucet() {
    await rechargeFromFaucet();
    sl('success', 'Waiting for confirmation...');
    await tryReloadBalance();
  }
</script>

{#if loaded}
  <div class="dark"> <!-- todo: find out why purgecss doesn't look up html.dark -->
    <div class="flex">
      <h1 class="main-heading flex-grow">BET TODAY,<br/>THE BEST CRYPTO OF TOMORROW, AND WIN!</h1>

      {#if $address}
        <div class="flex flex-col text-sm">
          <div class='mr-3'>Account: {$address}</div>
          <table class='balances'>
            <tr>
              <td>Balances:</td>
              <td>
                {fromMicro($balance || 0)}BET <span class="cursor-pointer underline" on:click={onRechargeFromFaucet}>recharge from faucet</span>
              </td>
            </tr>
            <!--
            {#if $myInfo}
              <tr>
                <td>Total Bets:</td>
                <td>{fromMicro($myInfo.totalBetsAmount)}BET</td>
              </tr>
              <tr>
                <td>Total Wins:</td>
                <td>{fromMicro($myInfo.totalWinsAmount)}BET</td>
              </tr>
            {/if}
            -->
          </table>
        </div>
        <button class="button is-light is-small ml-2" on:click={loadBalance}>
          REFRESH BALANCE
        </button>
        <button class="button is-light is-small ml-2" on:click={disconnectAccount}>
          DISCONNECT
        </button>
      {:else}
        <button class="button is-light is-small mr-3" on:click={generateAccount}>
          GENERATE ACCOUNT
        </button>
        <button class="button is-light is-small" on:click={connectAccount}>
          CONNECT ACCOUNT
        </button>
      {/if}
    </div>

    {#if info}
      <div class="main-container">
        {#each days as day}
          {#if $info}
            <Day day={day.date} label={day.label}/>
          {/if}
        {/each}
      </div>
    {/if}
  </div>
{/if}

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

  }
</style>
