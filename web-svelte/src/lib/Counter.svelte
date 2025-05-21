<script>
    import {
        increment,
        reset,
        getCount,
        getCounterStyle,
    } from "./count.svelte.js";
</script>

<div id="counter-box" class="box">
    <h1 id="counter" class={"counter-" + getCounterStyle()}>{getCount()}</h1>
    <div>
        <button id="increment" onclick={increment}>increment</button>
        <button id="reset" onclick={reset}>reset</button>
    </div>
</div>

<style>
    #counter {
        height: 40vh;
        max-height: 40vh;
        line-height: 40vh;
        margin: 5vh;
    }
    button {
        background-color: var(--rich-black);
        color: white;
        margin: 1em;
        border-radius: 1em;
        border: none;
        height: 3em;
        font-size: 1em;
        width: 7em;
        cursor: pointer;
    }

    h1 {
        margin: 0;
    }

    .counter-default {
        font-size: 3em;
    }
    .counter-10 {
        font-size: 3.5em;
    }
    .counter-25 {
        font-size: 4em;
    }
    .counter-50 {
        font-size: 5em;
    }
    .counter-69 {
        font-size: 5em;
        background-image: var(--rainbow);
        color: transparent;
        background-clip: text;
    }
    .counter-100 {
        font-size: 7em;
    }
    .counter-200 {
        font-size: 8em;
        animation: vertical-shaking 0.5s infinite;
    }
    .counter-420 {
        font-size: 8em;
        color: #08520c;
        animation: vertical-shaking 0.5s infinite;
    }
    .counter-500 {
        font-size: 9em;
        animation: vertical-shaking 0.3s infinite;
    }
    .counter-1000 {
        font-size: 10em;
        background-image: var(--rainbow);
        color: transparent;
        background-clip: text;
        animation: vertical-shaking 0.3s infinite;
    }
    .counter-10000 {
        font-size: 15em;
        background-image: var(--rainbow);
        color: transparent;
        background-clip: text;
        animation: vertical-shaking 0.2s infinite;
    }
    .counter-100000 {
        font-size: 20em;
        background-image: var(--rainbow);
        color: transparent;
        background-clip: text;
        animation: vertical-shaking-more 0.1s infinite;
    }
    .counter-1000000 {
        font-size: 40em;
        background-image: var(--rainbow);
        color: transparent;
        background-clip: text;
        animation: vertical-shaking-more 0.05s infinite;
    }

    @keyframes vertical-shaking {
        0% {
            transform: translateY(0);
        }
        25% {
            transform: translateY(5px);
        }
        50% {
            transform: translateY(-5px);
        }
        75% {
            transform: translateY(5px);
        }
        100% {
            transform: translateY(0);
        }
    }
    @keyframes vertical-shaking-more {
        0% {
            transform: translateY(0);
        }
        25% {
            transform: translateY(10px);
        }
        50% {
            transform: translateY(-10px);
        }
        75% {
            transform: translateY(10px);
        }
        100% {
            transform: translateY(0);
        }
    }
</style>
