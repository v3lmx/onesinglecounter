<script>
    import "./app.css";
    import Counter from "./lib/Counter.svelte";
    import Best from "./lib/Best.svelte";
    import About from "./lib/About.svelte";
    import Header from "./lib/Header.svelte";
    import Toggle from "./lib/Toggle.svelte";
    import { ModeWatcher } from "mode-watcher";
</script>

<main>
    <ModeWatcher />
    <Toggle />
    <Header />
    <Counter />
    <Best />
    <About />
</main>
