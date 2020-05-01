<script>
  import { onMount } from 'svelte';
  import account from '../utils/account';
  import sl from '../utils/sl';
  import {BETCHAIN_TRANSFER_CHANNEL} from '../config';

  let address;
  let balances = [];

  onMount(async function () {
    if (await account.loadPrivateKeyFromCache()) {
      await loadAccount();
    }
  });

  async function connectAccount() {
    const mnemonic = prompt('Enter mnemonic:');
    if (!mnemonic) {return;}
    if (await account.loadPrivateKeyFromMnemonic(mnemonic)) {
      await loadAccount();
    }
  }

  function disconnectAccount() {
    account.disconnect();
    address = null;
  }

  async function loadAccount() {
    address = account.deriveAddress();
    await account.loadAccount();
    await onLoadBalance();
  }

  async function onLoadBalance() {
    balances = await account.query(`/bank/balances/${address}`);
  }

  async function onBuyGold() {
    try {
      await account.tx('post', '/coinpricebet/buy', {
        amount: `1000000000transfer/${BETCHAIN_TRANSFER_CHANNEL}/uatom`
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
    <h1 class="main-heading flex-grow">GOLDCDP</h1>
    {#if address}
      <button class="button is-light is-small" on:click={onBuyGold}>
        BUY GOLD
      </button>
      <button class="button is-light is-small ml-2" on:click={onLoadBalance}>
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

  {#if address}
    <div class="mb-5">
      <div>Account: {address}</div>
      <div>Balance:</div>
      <div class="ml-5">
        {#each balances as c}
          <div>{c.amount} {c.denom}</div>
        {/each}
      </div>
    </div>
  {/if}
</div>
