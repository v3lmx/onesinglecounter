<script>
    import { best, timespansList } from "./count.svelte.js";
</script>

<div id="best-box" class="box">
    <h4>highest count</h4>
    {#each timespansList as time}
        <span id={"best-" + time}>{time}: {best.get(time)}</span>
    {/each}
    <br />
</div>

<style>
</style>
