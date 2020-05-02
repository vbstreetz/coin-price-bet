<script>
  import moment from 'moment';
  import {onMount} from 'svelte';
  import Day from '../components/day.svelte';
  import account from '../utils/account';
  import {fromMicro} from '../utils/cosmos';
  import {BETCHAIN_TRANSFER_CHANNEL, COINS} from '../config';

  let address;
  let balance = 0;
  let info;
  let myInfo;

  const tomorrow = moment.utc().add(1);
  const today = moment.utc();
  const yesterday = moment.utc().subtract(1);
  const days = [
    {label: 'Tomorrow', date: tomorrow},
    {label: 'Today', date: today},
    {label: 'Yesterday', date: yesterday}
  ];

  onMount(async function () {
    await Promise.all([loadInfo(), loadAccount()]);
  });

  async function connectAccount() {
    const mnemonic = prompt('Enter mnemonic:');
    if (!mnemonic) {
      return;
    }
    loadAccount(mnemonic);
  }

  function disconnectAccount() {
    account.disconnect();
    address = null;
  }

  async function loadInfo() {
    info = await account.query(`/coinpricebet/info`);
  }

  async function loadAccount(mnemonic) {
    if (mnemonic) {
      if (!await account.loadPrivateKeyFromMnemonic(mnemonic)) {
        return;
      }
    } else {
      if (!await account.loadPrivateKeyFromCache()) {
        return;
      }
    }
    address = account.deriveAddress();
    await Promise.all([
      account.loadAccount(),
      loadBalance(),
      loadMyInfo(),
    ]);
  }

  async function loadBalance() {
    const coins = await account.query(`/bank/balances/${address}`);
    for (let i = 0; i < coins.length; i++) {
      const coin = coins[i];
      if (~coin.denom.search(`transfer/${BETCHAIN_TRANSFER_CHANNEL}/uatom`)) {
        balance = coin.amount;
        break;
      }
    }
  }

  async function loadMyInfo() {
    myInfo = await account.query(`/coinpricebet/info/${address}`);
  }
</script>

<div>
  <div class="flex">
    <h1 class="main-heading flex-grow">BET TODAY,<br/>THE BEST CRYPTO OF TOMORROW, AND WIN!</h1>

    {#if address}
      <div class="flex flex-col text-sm">
        <div>Account: {address}</div>
        <div>Balance: {fromMicro(balance)}atom</div>
        {#if myInfo}
        <div>Total Bets: {fromMicro(myInfo.totalBetsAmount)}atom</div>
        <div>Total Wins: {fromMicro(myInfo.totalWinsAmount)}atom</div>
        {/if}
      </div>
      <button class="button is-light is-small ml-2" on:click={disconnectAccount}>
        DISCONNECT
      </button>
    {:else}
      <button class="button is-light is-small" on:click={connectAccount}>
        CONNECT
      </button>
    {/if}
  </div>

  {#if info}
    <div class="main-container">
      {#each days as day}
        <Day firstDay={info.firstDay} day={day.date} label={day.label} address={address} />
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
    background: rgba(0, 0, 0, .12);
    bottom: 0;
    content: "";
    height: 100%;
    position: absolute;
    top: 0;
    width: 1px;
  }
</style>
