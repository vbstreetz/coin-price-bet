<script>
  import _ from 'lodash';
  import {onMount} from 'svelte';
  import {get, derived} from 'svelte/store';
  import moment from 'moment';
  import Chart from '../components/chart.svelte';
  import Clock from '../components/clock.svelte';
  import {fromMicro, toMicro, formatFiat, sleep, sl} from '../utils';
  import {address, info, blockchain, tryReloadBalance} from '../stores/blockchain';
  import {COINS, DAY_STATES} from '../config';

  export let day;
  export let label;

  const dayId = day.diff(moment.unix(get(info).firstDay), 'days');
  const dateLabel = day.format('YYYY-MM-DD');

  let dayInfo;
  let myDayInfo;
  let coinStats = [];
  let chartY;

  let canBet;
  let isDrawing;
  let isPayout;

  let isBetting;
  let isWithdrawingWins;

  address.subscribe(() => {
    Promise.all([
      loadDayInfo(),
      loadMyDayInfo()
    ])
  });

  $: {
    if (dayInfo) {
      let perfs = dayInfo.coinsPerf.map((perf, id) => ({id, perf}));
      if (!canBet) {
        perfs = _.orderBy(perfs, 'perf', 'desc');
      }
      coinStats = perfs.map(({id, perf}) => {
        return {
          id,
          perf,
          amount: !myDayInfo ? 0 : parseInt(myDayInfo.coinBetTotalAmount[id])
        };
      });

      chartY = dayInfo.coinsVolume.map(v => fromMicro(v));
    } else {
      coinStats = COINS.map((c, id) => ({
        id,
        perf: 0,
        amount: 0,
      }));
    }
  }

  onMount(async function () {
    await Promise.all([
      loadDayInfo(),
      loadMyDayInfo()
    ]);
  });

  async function loadDayInfo() {
    const info = await blockchain.query(`/coinpricebet/day-info/${dayId}`);
    info.grandPrizeAmount = parseInt(info.grandPrizeAmount);
    info.currencyPriceUSD = Math.pow(10, 6); // parseInt(info.currencyPriceUSD);
    info.coinsPerf = info.coinsPerf.map(n => parseInt(n));
    info.coinsVolume = info.coinsVolume.map(n => parseInt(n));
    info.state = parseInt(info.state);
    dayInfo = info;

    canBet = dayInfo.state === DAY_STATES.BET;
    isDrawing = dayInfo.state === DAY_STATES.DRAWING;
    isPayout = dayInfo.state === DAY_STATES.PAYOUT;
  }

  async function loadMyDayInfo() {
    const $address = get(address);
    if ($address) {
      const info = await blockchain.query(`/coinpricebet/day-info/${dayId}/${$address}`);
      info.totalBetAmount = parseInt(info.totalBetAmount);
      info.coinBetTotalAmount = info.coinBetTotalAmount.map(n => parseInt(n));
      myDayInfo = info;
    }
  }

  async function onSendPrediction(event) {
    event.preventDefault();

    isBetting = false;

    const coin = event.target.coin.value;
    const amount = event.target.amount.value;
    const $info = get(info);

    try {
      sl('info', `Supporting ${coin} win prediction with ${amount}BET`, 'Place bet?', async () => {
        await blockchain.tx('post', '/coinpricebet/place-bet', {
          amount: `${toMicro(amount).toString()}stake`,
          coinId: COINS.indexOf(coin)
        });
        sl('success', 'Waiting for confirmation...');
        await sleep(3000);
        await Promise.all([
          loadDayInfo(),
          loadMyDayInfo(),
          tryReloadBalance(),
        ]);
      });
    } catch (e) {
      sl('error', e);
    } finally {
      isBetting = false;
    }
  }

  async function onWithdrawWins(event) {
    isWithdrawingWins = true;
    try {
      await blockchain.tx('post', '/coinpricebet/payout', {
        dayId: dayId.toString()
      });
      sl('success', 'Done!');
      await sleep(3000);
      await Promise.all([
        loadDayInfo(),
        loadMyDayInfo(),
        tryReloadBalance(),
      ]);
    } finally {
      isWithdrawingWins = false;
    }
  }
</script>

<div class="day-container flex dark"> <!-- todo: find out why purgecss doesn't look up html.dark -->
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
              <div class="amount">{fromMicro(dayInfo.grandPrizeAmount)}<span class="currency">BET</span></div>
              <ul class="fiat">
                <li>~{formatFiat(fromMicro(dayInfo.grandPrizeAmount) * fromMicro(dayInfo.currencyPriceUSD), 'USD')}<span
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
              {#if canBet || chartY}
                <Chart y={chartY} x={COINS} type='bar' yLabel="BET"/>
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
                  <table class="table text-sm">
                    <thead>
                    <tr>
                      <th>Symbol</th>
                      {#if !canBet}
                        <th></th>
                      {/if}
                      <th align="right">Amount</th>
                    </tr>
                    </thead>
                    <tbody>
                    {#each coinStats as stat}
                      <tr>
                        <td>{COINS[stat.id]}</td>
                        {#if !canBet}
                          <td class="text-xs ml-3">{(fromMicro(stat.perf) * 100).toFixed(2)}%</td>
                        {/if}
                        <td align="right">{fromMicro(stat.amount)} <small class="text-xs">BET</small></td>
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
                    <label for="bet-input">I'm supporting my prediction with this amount of BET:</label>
                    <br/>
                    <input required="required" class="input" type="text" id="bet-input" name="amount"
                           placeholder="Enter amount of BET..."/>
                  </div>

                  <div class="flex flex-grow mt-3">
                    <button type="submit" class="button is-link flex-grow" disabled={isBetting || !$address}>
                      Send prediction
                    </button>
                  </div>
                </form>
              </div>
            </div>
          {:else if isDrawing}
            {#if myDayInfo}
              <div class="mt-4">Your total bet amount: {fromMicro(myDayInfo.totalBetAmount)} BET</div>
              <div>Your predicted total win amount: {fromMicro(myDayInfo.totalWinAmount)} BET</div>
            {/if}
          {:else }
            {#if myDayInfo}
              <div class="mt-4">Your total bet amount: {fromMicro(myDayInfo.totalBetAmount)} BET</div>
              <div>Your total win amount: {fromMicro(myDayInfo.totalWinAmount)} BET</div>
              <div class="flex flex-grow mt-3">
                <button type="submit" class="button is-link flex-grow"
                        disabled={isWithdrawingWins || myDayInfo.paid || !myDayInfo.totalWinAmount}
                        on:click={onWithdrawWins}>
                  Withdraw Wins
                </button>
              </div>
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

  .dark .day-divider-dot-inner {
    border-color: #333;
  }

  .day-divider-dot-inner--active {
    background-color: var(--dot-active-color);
    border-color: var(--dot-active-color);
  }

  .dark .day-divider-dot-inner--active {
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

  .dark .day-card {
    background: var(--dark-1);
    border: none;
    box-shadow: 0 2px 1px -1px rgba(0, 0, 0, 0.2), 0px 1px 1px 0px rgba(0, 0, 0, 0.14), 0px 1px 3px 0px rgba(0, 0, 0, 0.12);
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

  .dark .day-card:before {
    border-color: var(--border-color);
  }

  .day-card .column:not(:last-child) {
    border-right: 1px solid #ddd;
  }

  .dark .day-card .column:not(:last-child) {
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

  .grand-prize .amount {
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

  .table {
    width: 100%;
  }
</style>
