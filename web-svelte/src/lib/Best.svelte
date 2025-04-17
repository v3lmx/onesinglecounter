<script>
import { best, refreshBest, timespansList } from "./count.svelte.js";
</script>

<div id="best-box" class="box">
    <h4>Highest count</h4>
    {#each timespansList as time}
        <span id={"best-" + time}>{time}: {best.get(time)}</span>
    {/each}
    <br />
    <button id="refresh-best" onclick={refreshBest}>refresh</button>
</div>
