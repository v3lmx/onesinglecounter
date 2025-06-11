import { SvelteMap } from "svelte/reactivity";

let count = $state(0);
export const best = new SvelteMap();

let counterStyle = $derived.by(() => {
    if (count < 10) {
        return "default"
    }
    if (count < 25) {
        return "10"
    }
    if (count < 50) {
        return "25"
    }
    if (count == 69) {
        return "69"
    }
    if (count < 100) {
        return "50"
    }
    if (count < 200) {
        return "100"
    }
    if (count == 420) {
        return "420"
    }
    if (count < 500) {
        return "200"
    }
    if (count < 1000) {
        return "500"
    }
    if (count < 10000) {
        return "1000"
    }
    if (count < 100000) {
        return "10000"
    }
    if (count < 1000000) {
        return "100000"
    }
    return "1000000"
});


export function getCount() {
    return count;
}
export function getCounterStyle() {
    return counterStyle;
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

const wsUrl = import.meta.env.VITE_WS_URL || "ws://localhost:8000/connect";
const ws = new WebSocket(wsUrl);

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

ws.onopen = () => {
    current();
};

ws.onmessage = (event) => {
    const e = event.data.split(":");
    if (e[0] === messageCurrent) {
        count = e[1];
        for (const time of timespansList) {
            if (Number(count) > Number(best.get(time))) {
                best.set(time, count);
            }
        }
    }
    if (e[0] === "alltime") {
        for (let i = 0; i < e.length - 1; i += 2) {
            const time = e[i];
            const value = e[i + 1];
            best.set(time, value);
        }
    }
};
