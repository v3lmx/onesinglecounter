const timespans = Object.freeze({
    MINUTE: Symbol("minute"),
    HOUR: Symbol("hour"),
    DAY: Symbol("day"),
    WEEK: Symbol("week"),
    MONTH: Symbol("month"),
    YEAR: Symbol("year"),
    ALLTIME: Symbol("alltime"),
});
const timespansList = ["minute", "hour", "day", "week", "month", "year", "alltime"]
const messageReset = "res"
const messageIncrement = "inc"
const messageCurrent = "current"
const messageBest = "best"

// let ws = new WebSocket("ws://localhost:8000/connect")
let ws = new WebSocket("wss://private.v3lmx.com/osc/api/connect")
ws.binaryType = 'arraybuffer';
ws.onopen = (event) => {
    current()
    best()
};
ws.onmessage = (event) => {
    if (event.data == messageReset) {
        var el = document.getElementById("counter");
        el.innerText = 0;
    }
    if (event.data == messageIncrement) {
        var el = document.getElementById("counter");
        el.innerText++;
    }
    let e = event.data.split(":")
    if (e[0] == messageCurrent) {
        var el = document.getElementById("counter");
        el.innerText = e[1];
    }
    if (e[0] == messageBest) {
        for (let i = 1; i < e.length - 1; i += 2) {
            let time = e[i]
            let value = e[i + 1]
            var el = document.getElementById("value-best-" + time);
            el.innerText = value;

        }
    }
    updateBest()

    console.log(event.data)
};

function reset() {
    ws.send(messageReset)
}

async function increment() {
    ws.send(messageIncrement)
}

function current() {
    ws.send(messageCurrent)
}

function best() {
    ws.send(messageBest)
}

function updateBest() {
    const current = document.getElementById("counter").innerText;

    let minuteBest = document.getElementById("value-best-minute");
    if (parseInt(current) <= parseInt(minuteBest.innerText)) {
        return
    }
    minuteBest.innerText = current;

    let hourBest = document.getElementById("value-best-hour");
    if (parseInt(current) <= parseInt(hourBest.innerText)) {
        return
    }
    hourBest.innerText = current;

    let dayBest = document.getElementById("value-best-day");
    if (parseInt(current) <= parseInt(dayBest.innerText)) {
        return
    }
    dayBest.innerText = current;

    let weekBest = document.getElementById("value-best-week");
    if (parseInt(current) <= parseInt(weekBest.innerText)) {
        return
    }
    weekBest.innerText = current;

    let monthBest = document.getElementById("value-best-month");
    if (parseInt(current) <= parseInt(monthBest.innerText)) {
        return
    }
    monthBest.innerText = current;

    let yearBest = document.getElementById("value-best-year");
    if (parseInt(current) <= parseInt(yearBest.innerText)) {
        return
    }
    yearBest.innerText = current;

    let alltimeBest = document.getElementById("value-best-alltime");
    if (parseInt(current) <= parseInt(alltimeBest.innerText)) {
        return
    }
    alltimeBest.innerText = current;

}

