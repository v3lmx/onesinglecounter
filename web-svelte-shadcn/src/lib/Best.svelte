<script>
    import { best, timespansList } from "./count.svelte.js";
</script>

<div
    class="fixed left-1/2 top-8/10 -translate-1/2 flex flex-col items-center text-s md:text-xl lg:text-2xl"
>
    {#each timespansList as time}
        <span id={"best-" + time}>{time}: {best.get(time)}</span>
    {/each}
    <br />
</div>

<style>
</style>
