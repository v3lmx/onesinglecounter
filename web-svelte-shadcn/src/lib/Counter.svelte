<script lang="ts">
    import { increment, reset, getCount } from "./count.svelte.js";
    import { Button } from "$lib/components/ui/button/index.js";
</script>

<div class="fixed left-1/2 top-3/10 -translate-1/2 flex flex-col items-center">
    <!-- <div class="flex flex-col justify-center items-center m-10"> -->
    <h1 id="counter" class="m-10 text-6xl md:text-8xl lg:text:12xl">
        {getCount()}
    </h1>
    <div>
        <Button id="increment" class="m-1 md:m-3" onclick={increment}
            >increment</Button
        >
        <Button id="reset" class="m-1: md:m-3" onclick={reset}>reset</Button>
    </div>
</div>

<style>
</style>
