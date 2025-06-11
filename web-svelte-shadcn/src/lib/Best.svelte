<script lang="ts">
    import { best, timespansList } from "./count.svelte.js";
    import * as Table from "$lib/components/ui/table/index.js";
</script>

<div
    class="fixed left-1/2 top-7/10 md:top-8/10 -translate-1/2 flex flex-col items-center text-s md:text-xl lg:text-2xl"
>
    <Table.Root class="text-s md:text-xl lg:text-2xl">
        <Table.Body>
            {#each timespansList as time}
                <Table.Row>
                    <Table.Cell class="w-[100px]">{time}</Table.Cell>
                    <Table.Cell class="w-[100px] text-right"
                        >{best.get(time)}</Table.Cell
                    >
                </Table.Row>
            {/each}
        </Table.Body>
    </Table.Root>
    <br />
</div>

<style>
</style>
