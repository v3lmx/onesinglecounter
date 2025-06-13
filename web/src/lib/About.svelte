<script lang="ts">
    import * as Popover from "$lib/components/ui/popover/index.ts";

    let text = `how high can we go?

the counter is the same for everyone, to be incremented by one, or reset to 0.

you can see the highest counts achieved for each current period.`;
</script>

<div class="fixed left-1/2 bottom-1 -translate-1/2">
    <Popover.Root>
        <Popover.Trigger>about</Popover.Trigger>
        <Popover.Content>
            <div style="white-space: pre-wrap;">{text}</div>
            <br />
            <div>
                <span>created by </span>
                <a
                    class="underline"
                    target="_blank"
                    rel="noopener noreferrer"
                    href="https://github.com/v3lmx"
                >
                    v3lmx ↗
                </a>
            </div>
            <br />
            <div>
                <span>see the code </span>
                <a
                    class="underline"
                    target="_blank"
                    rel="noopener noreferrer"
                    href="https://github.com/v3lmx/onesinglecounter"
                >
                    here ↗
                </a>
            </div>
            <br />
            <div>
                <span>inspired by </span>
                <a
                    class="underline"
                    target="_blank"
                    rel="noopener noreferrer"
                    href="https://eieio.games/blog/one-million-checkboxes/"
                >
                    One Million Checkboxes ↗
                </a>
            </div>
        </Popover.Content>
    </Popover.Root>
</div>
