<script>
    import { increment, reset, getCount } from "./count.svelte.js";
</script>

<h1>{getCount()}</h1>

<button id="increment" onclick={increment}>increment</button>
<button id="reset" onclick={reset}>reset</button>
