<script>
    let isOpen = false;
    let text = `how high will it go?

the counter is the same for everyone, to be incremented by one, or reset to 0.

"highest count" keeps track of the, well, highest counts achieved.`;

    const open = function () {
        isOpen = true;
    };
    const close = function () {
        isOpen = false;
    };
</script>

<span id="about" on:click={open}>about</span>

{#if isOpen}
    <div class="modal-overlay" on:click={close}>
        <div class="modal" on:click|stopPropagation>
            <button class="close-button" on:click={close}>×</button>
            <div style="white-space: pre-wrap;">{text}</div>
            <br />
            <div>
                <span>created by </span>
                <a href="https://github.com/v3lmx"> v3lmx </a>
            </div>
            <br />
            <div>
                <span>inspired by </span>
                <a href="https://eieio.games/blog/one-million-checkboxes/">
                    One Million Checkboxes
                </a>
            </div>
        </div>
    </div>
{/if}

<style>
    #about {
        position: fixed;
        left: 50%;
        bottom: 1em;
        transform: translate(-50%, -50%);
        margin: 0 auto;
        cursor: pointer;
    }
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .modal {
        background: var(--rich-black);
        color: white;
        padding: 4em;
        border-radius: 5px;
        position: relative;
        text-align: center;
    }

    .close-button {
        position: absolute;
        color: white;
        top: 10px;
        right: 10px;
        background: none;
        border: none;
        font-size: 1.5em;
        cursor: pointer;
    }

    a {
        color: var(--silver-lake-blue);
        text-decoration: none;
    }
</style>
