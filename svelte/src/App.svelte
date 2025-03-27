<script>
    import Counter from "./lib/Counter.svelte";
    import Best from "./lib/Best.svelte";
</script>

<main>
    <div class="card">
        <Counter />
        <Best />
    </div>
</main>

<style>
    #counter {
        border-radius: 0.5em;
        padding: 20px;
        width: 10em;
        text-align: center;
    }

    :global .box {
        border-radius: 0.5em;
        font: 1.2em sans-serif;
        margin: 5vh;

        display: flex;
        align-items: center;
        justify-content: space-evenly;
        flex-direction: column;
    }

    #best-box {
        height: 30vh;
    }

    #counter-box {
        height: 60vh;
    }
</style>
