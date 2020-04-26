<script>
    import { writable } from 'svelte/store';

    const TOKENS_MAP = {
        BTC: 0,
        ETH: 1,
        LTC: 2,
        BAND: 3,
        ATOM: 4,
        LINK: 5,
        XTZ: 6
    };
    const TOKENS_LIST = Object.entries(TOKENS_MAP);
    const prices = writable({});
</script>

<style>
    button {
        margin-left: 10px;
    }
</style>

<div class="flex">
    {#each TOKENS_LIST as [token, id]}
        <button>{token} - {prices[id] || 0}</button>
    {/each}
</div>