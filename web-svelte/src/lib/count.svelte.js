import { SvelteMap } from "svelte/reactivity";

let count = $state(0);
export const best = new SvelteMap();

export function getCount() {
    return count;
}

export const timespansList = [
    "minute",
    "hour",
    "day",
    "week",
    "month",
    "year",
    "alltime",
];

const messageReset = "res";
const messageIncrement = "inc";
const messageCurrent = "current";
const messageBest = "best";

for (const time of timespansList) {
    best.set(time, 0);
}

const ws = new WebSocket("ws://localhost:10001/connect");
ws.binaryType = "arraybuffer";

export function reset() {
    ws.send(messageReset);
}

export function increment() {
    ws.send(messageIncrement);
}

function current() {
    ws.send(messageCurrent);
}

export function refreshBest() {
    ws.send(messageBest);
}

ws.onopen = (event) => {
    current();
    // getBest();
};

ws.onmessage = (event) => {
    if (event.data === messageReset) {
        count = 0;
    }
    if (event.data === messageIncrement) {
        count++;
        for (const time of timespansList) {
            if (count > best.get(time)) {
                best.set(time, count);
            }
        }
    }
    const e = event.data.split(":");
    if (e[0] === messageCurrent) {
        count = e[1];
        for (const time of timespansList) {
            if (Number(count) > best.get(time)) {
                best.set(time, count);
            }
        }
    }
    if (e[0] === messageBest) {
        for (let i = 1; i < e.length - 1; i += 2) {
            const time = e[i];
            const value = e[i + 1];
            best.set(time, value);
        }
    }

    console.log(event.data);
};
