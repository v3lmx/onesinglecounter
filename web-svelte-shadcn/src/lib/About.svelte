<script lang="ts">
    import * as Popover from "$lib/components/ui/popover/index.ts";
    import { ExternalLink } from "lucide";

    let text = `how high will it go?

the counter is the same for everyone, to be incremented by one, or reset to 0.

"highest count" keeps track of the, well, highest counts achieved.`;
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
