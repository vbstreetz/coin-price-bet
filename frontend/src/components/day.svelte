<script>
  import {onMount} from 'svelte';
  import {get} from 'svelte/store';
  import moment from 'moment';
  import Clock from '../components/clock.svelte';
  import {coinPriceBetBlockchain} from '../utils/blockchains';
  import {fromMicro, toMicro} from '../utils/cosmos';
  import {formatFiat} from '../utils';
  import sl from '../utils/sl';
  import {address, info} from '../stores/blockchains';
  import {COINS, DAY_STATES} from '../config';

  export let day;
  export let label;

  const daySinceEpoch = day.diff(moment.unix(get(info).firstDay), 'days');
  const dateLabel = day.format('YYYY-MM-DD');

  let dayInfo;
  let myDayInfo;

  let canBet;
  let isDrawing;
  let isPayout;

  onMount(async function () {
    Promise.all([
      loadDayInfo(),
      loadMyDayInfo()
    ])
  });

  async function loadDayInfo() {
    dayInfo = await coinPriceBetBlockchain.query(`/coinpricebet/day-info/${daySinceEpoch}`);

    canBet = dayInfo.state === DAY_STATES.BET;
    isDrawing = dayInfo.state === DAY_STATES.DRAWING;
    isPayout = dayInfo.state === DAY_STATES.PAYOUT;
  }

  async function loadMyDayInfo() {
    const $address= get(address);
    if ($address) {
      myDayInfo = await coinPriceBetBlockchain.query(`/coinpricebet/day-info/${daySinceEpoch}/${$address}`);
    }
  }

  async function onSendPrediction(event) {
    event.preventDefault();
    const coin = event.target.coin.value;
    const atom = event.target.atom.value;

    try {
      await coinPriceBetBlockchain.tx('post', '/coinpricebet/place-bet', {
        amount: `${toMicro(atom)} transfer/${$info.betchainTransferChannel}/uatom`,
        coinId: COINS.indexOf(coin)
      });
      sl('success', 'WAITING FOR CONFIRMATION...');
    } catch (e) {
      sl('error', e);
    }
  }
</script>

<div class="day-container flex">
  {#if dayInfo}
    <div class="day-divider flex justify-center">
      <div class="day-divider-dot">
        <div
          class="day-divider-dot-inner {canBet ? 'day-divider-dot-inner--active' : ''} flex items-center justify-center"></div>
      </div>
    </div>
    <div class="day-body">
      <div class="day-card">
        <!-- day-state -->
        <div class="column day-state">
          <div class="row field">
            <header>
              <h1 class="day">{label}</h1>
              <h3 class="date">({dateLabel})</h3>
            </header>
            <main class="status">
              {#if canBet}
                Predictions are open
              {:else if isDrawing}
                Drawing... <!-- Waiting for results -->
              {:else}
                Finalized
              {/if}
            </main>
          </div>
          <div class="row field">
            {#if canBet || isDrawing}
              <header><h1>Closing in</h1></header>
              <Clock/>
            {/if}
          </div>
          <div class="row field">
            <header>
              <h1>
                {#if canBet}
                  Current
                {/if}
                grand prize
              </h1>
            </header>
            <div class="grand-prize">
              <div class="atom">{fromMicro(dayInfo.grandPrizeAmount)}<span class="currency">ATOM</span></div>
              <ul class="fiat">
                <li>~{formatFiat(dayInfo.grandPrizeAmount * dayInfo.atomPriceCents / 100, 'USD')}<span
                  class="currency">USD</span></li>
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
                {#if canBet}
                  Current prediction volumes
                {:else}
                  Final prediction volumes
                {/if}
              </h1>
            </header>
            <div class="small">
              <!-- chart -->
              {#if canBet || dayInfo.coinsVolume.length}
                <ul>
                  {#each COINS as c, i}
                    <li>{c} - {fromMicro(dayInfo.coinsVolume[i] || 0)} ATOM</li>
                  {/each}
                </ul>
              {:else}
                There were no predictions during this day.
              {/if}
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
                    {#each COINS as coin, i}
                      <tr>
                        <td class="text-left">
                          {coin}
                          {#if !canBet}
                            <small class="text-xs">{fromMicro(dayInfo.coinsPerf[i] || 0).toFixed(2)}%</small>
                          {/if}
                        </td>
                        {#if myDayInfo}
                          <td class="text-right">{fromMicro(myDayInfo.coinBetTotalAmount[i] || 0)} ATOM</td>
                          <td class="text-right">{fromMicro(myDayInfo.coinPredictedWinAmount[i] || 0)} ATOM</td>
                        {/if}
                      </tr>
                    {/each}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
          {#if canBet}
            <div class="mt-5">
              <div class="row field">
                <header><h1>Predict tomorrow's best crypto</h1></header>
                <form on:submit={onSendPrediction}>
                  <label for="bet-select">I predict that the best performing crypto during tomorrow will be:</label>
                  <div class="select">
                    <select id="bet-select" name="coin">
                      {#each COINS as coin}
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
          {:else if isDrawing}
            {#if myDayInfo}
              <div class="mt-4">Your total bet amount: {fromMicro(myDayInfo.totalBetAmount)} ATOM</div>
              <div>Your predicted total win amount: {fromMicro(myDayInfo.totalWinAmount)} ATOM</div>
            {/if}
          {:else }
            {#if myDayInfo}
              <div class="mt-4">Your total bet amount: {fromMicro(myDayInfo.totalBetAmount)} ATOM</div>
              <div>Your total win amount: {fromMicro(myDayInfo.totalWinAmount)} ATOM</div>
            {/if}
          {/if}
        </div>
        <!-- end my bets -->
      </div>
    </div>
  {/if}
</div>

<style>
  .row {
    display: flex;
    flex-wrap: wrap;
    flex: 1 1 auto;
    margin-right: -12px;
    margin-left: -12px;
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

  :global(.dark) .day-divider-dot-inner {
    border-color: #333;
  }

  .day-divider-dot-inner--active {
    background-color: var(--dot-active-color);
    border-color: var(--dot-active-color);
  }

  :global(.dark) .day-divider-dot-inner--active {
    border-color: #333;
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

  :global(.dark) .day-card {
    background: var(--dark-1);
    border: none;
    box-shadow: 0px 2px 1px -1px rgba(0,0,0,0.2), 0px 1px 1px 0px rgba(0,0,0,0.14), 0px 1px 3px 0px rgba(0,0,0,0.12);
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

  :global(.dark) .day-card:before {
    border-color: var(--border-color);
  }

  .day-card .column:not(:last-child) {
    border-right: 1px solid #ddd;
  }

  :global(.dark) .day-card .column:not(:last-child) {
    border-color: #333;
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
    color: #888;
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
    color: #888;
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
    color: #888;
    font-size: 1.1rem;
    list-style-type: none;
    padding: 0;
    text-align: right;
  }

  .day-card .small {
    max-width: 60vw;
  }
</style>