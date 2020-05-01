<script>
  import moment from 'moment';
  import {onMount} from 'svelte';
  import Clock from '../components/clock.svelte';
  import account from '../utils/account';
  import {fromMicro, toMicro} from '../utils/cosmos';
  import sl from '../utils/sl';
  import {BETCHAIN_TRANSFER_CHANNEL} from '../config';

  let address;
  let balance = 0;

  const tomorrow = moment.utc().add(1);
  const today = moment.utc();
  const yesterday = moment.utc().subtract(1);
  let days = [
    {label: 'Tomorrow', dateLabel: tomorrow.format('YYYY-MM-DD'), active: true},
    {label: 'Today', dateLabel: today.format('YYYY-MM-DD')},
    {label: 'Yesterday', dateLabel: yesterday.format('YYYY-MM-DD')}
  ];
  let coins = [
    "BTC",
    "ETH",
    "LTC",
    "BAND",
    "ATOM",
    "LINK",
    "XTZ",
  ];

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
    const coins = await account.query(`/bank/balances/${address}`);
    for (let i = 0; i < coins.length; i++) {
      const coin = coins[i];
      if (~coin.denom.search(`transfer/${BETCHAIN_TRANSFER_CHANNEL}/uatom`)) {
        balance = coin.amount;
        break;
      }
    }
  }

  async function onSendPrediction(event) {
    event.preventDefault();
    const coin = event.target.coin.value;
    const atom = event.target.atom.value;

    try {
      await account.tx('post', '/coinpricebet/place-bet', {
        amount: `${toMicro(atom)} transfer/${BETCHAIN_TRANSFER_CHANNEL}/uatom`,
        coin: coins.indexOf(coin)
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

    {#if address}
      <div class="flex flex-col text-sm">
        <div>Account: {address}</div>
        <div>Balance: {fromMicro(balance)}atom</div>
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

  <div class="main-container">
    {#each days as day}
      <div class="day-container flex">
        <div class="day-divider flex justify-center">
          <div class="day-divider-dot">
            <div
              class="day-divider-dot-inner {day.active ? 'day-divider-dot-inner--active' : ''} flex items-center justify-center"></div>
          </div>
        </div>
        <div class="day-body">
          <div class="day-card">
            <!-- day-state -->
            <div class="column day-state">
              <div class="row field">
                <header>
                  <h1 class="day">{day.label}</h1>
                  <h3 class="date">({day.dateLabel})</h3>
                </header>
                <main class="status">
                  {#if day.active}
                    Predictions are open
                  {:else}
                    Finalized
                  {/if}
                </main>
              </div>
              <div class="row field">
                {#if day.active}
                  <header><h1>Closing in</h1></header>
                  <Clock/>
                {/if}
              </div>
              <div class="row field">
                <header><h1>Current grand prize</h1></header>
                <div class="grand-prize">
                  <div class="atom">0.0000<span class="currency">ATOM</span></div>
                  <ul class="fiat">
                    <li>~0.00<span class="currency">USD</span></li>
                    <li>~0.00<span class="currency">EUR</span></li>
                    <li>~0<span class="currency">JPY</span></li>
                  </ul>
                </div>
              </div>
            </div>
            <!-- end bets -->
            <!-- graph -->
            <div class="column graph">
              <div class="row field">
                <header>
                  <h1>
                    {#if day.active}
                      Current prediction volumes
                    {:else}
                      RESULTS
                    {/if}
                  </h1>
                </header>
                <div class="small">
                  <!-- chart -->
                </div>
              </div>
            </div>
            <!-- end graph -->
            <!-- my bets -->
            <div class="column my-bets">
              <div>
                <div class="my-bets field">
                  <header><h1>My predictions</h1></header>
                  <div>
                    <div>
                      <table class="table">
                        <thead>
                        <tr>
                          <th>Symbol</th>
                          <th>Amount</th>
                          <th>Potential Win</th>
                        </tr>
                        </thead>
                        <tbody>
                        {#each coins as coin}
                          <tr>
                            <td class="text-left">{coin}</td>
                            <td class="text-right">0.00 ATOM</td>
                            <td class="text-right">0.00 ATOM</td>
                          </tr>
                        {/each}
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
              </div>
              {#if day.active}
                <div class="mt-5">
                  <div class="row field">
                    <header><h1>Predict tomorrow's best crypto</h1></header>
                    <form on:submit={onSendPrediction}>
                      <label for="bet-select">I predict that the best performing crypto during tomorrow will be:</label>
                      <div class="select">
                        <select id="bet-select" name="coin">
                          {#each coins as coin}
                            <option value="{coin}">{coin}</option>
                          {/each}
                        </select>
                      </div>

                      <div class="mt-3">
                        <label for="bet-input">I'm supporting my prediction with this amount of ATOM:</label>
                        <br/>
                        <input required="required" class="input" type="text" id="bet-input" name="atom"/>
                      </div>

                      <div class="flex flex-grow mt-3">
                        <button type="submit" class="button is-link flex-grow">
                          Send prediction
                        </button>
                      </div>
                    </form>
                  </div>
                </div>
              {/if}
            </div>
            <!-- end my bets -->
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .main-heading {
    font-weight: bolder;
    margin-left: 50px;
  }

  .row {
    display: flex;
    flex-wrap: wrap;
    flex: 1 1 auto;
    margin-right: -12px;
    margin-left: -12px;
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

  .day-container {
    padding-bottom: 24px;
  }

  .day-divider {
    position: relative;
    min-width: 96px;
  }

  .day-divider-dot {
    z-index: 2;
    border-radius: 50%;
    box-shadow: 0 2px 1px -1px rgba(0, 0, 0, .2), 0 1px 1px 0 rgba(0, 0, 0, .14), 0 1px 3px 0 rgba(0, 0, 0, .12);
    height: 38px;
    left: calc(50% - 19px);
    width: 38px;

    height: 24px;
    left: calc(50% - 12px);
    width: 24px;
  }

  .day-divider-dot-inner {
    height: inherit;
    margin: 0;
    width: inherit;
    border-radius: 50%;
    background-color: hsl(0, 0%, 96%);
    border-color: hsl(0, 0%, 96%);
  }

  .day-divider-dot-inner--active {
    background-color: hsl(171, 100%, 41%);
    border-color: hsl(171, 100%, 41%);
  }

  .day-body {
    max-width: calc(100% - 96px);
    position: relative;
    height: 100%;
    -webkit-box-flex: 1;
    -ms-flex: 1 1 auto;
    flex: 1 1 auto;
  }

  .day-card {
    background-color: #fff;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-shadow: 0 5px 10px rgba(0, 0, 0, .1);
    display: flex;
    padding-left: 0;
  }

  .day-card:before {
    border-top: 1px solid #ccc;
    content: "";
    display: block;
    width: 50px;
    position: absolute;
    left: -50px;
    top: 12px;
  }

  .day-card .column:not(:last-child) {
    border-right: 1px solid #ddd;
  }

  .day-card .column:first-child {
    padding-left: 45px;
  }

  .day-card > .column {
    width: 33%;
    flex: 1;
  }

  .day-card .column {
    padding: 35px;
  }

  .day-card .day-state {
    padding-bottom: 20px;
  }

  .day-card .column .row {
    margin-bottom: 20px;
  }

  .day-card .field header {
    color: #666;
    display: block;
    font-weight: 700;
    text-transform: uppercase;
    width: 100%;
  }

  .day-card .field .status {
    font-size: 1.8rem;
    line-height: 1.2em;
    margin-top: .2em;
  }

  .day-card .field header h1.day {
    color: #666;
  }

  .day-card .field header .date {
    color: #aaa;
    display: inline;
    font-size: .9rem;
    margin-left: 10px;
  }

  .day-card .field header h1 {
    color: #b3b3b3;
    display: inline;
    font-size: 1.1rem;
  }

  .grand-prize .atom {
    border: 2px solid hsl(204, 71%, 53%);
    border-radius: 4px;
    font-size: 3em;
    margin: 10px 0;
    padding: 20px 30px;
  }

  .grand-prize .currency {
    color: #999;
    font-size: .7em;
    margin-left: .3em;
  }

  .grand-prize .fiat {
    color: #666;
    font-size: 1.1rem;
    list-style-type: none;
    padding: 0;
    text-align: right;
  }

  .day-card .small {
    max-width: 60vw;
  }
</style>
