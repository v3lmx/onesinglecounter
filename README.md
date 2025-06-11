# onesinglecounter

This project is an experiment to create a single globally available counter. Users can either increment the counter by one, or reset it to zero. 

The counter will be live at [onesinglecounter.com](https://onesinglecounter.com).

It is also an excuse to experiment with technologies and concepts:
- **concurrency patterns**, eg. using channels to store values
- `sync` primitives like `atomic` values, `Cond` to syncronize events across goroutines
- observability setups using `VictoriaMetrics`, `VictoriaLogs` and `Grafana`
- `shadcn` components 
- `svelte` 5

