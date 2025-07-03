# onesinglecounter

This project is an experiment to create a single globally available counter. Users can either increment the counter by one, or reset it to zero. 

The counter is live at [onesinglecounter.com](https://onesinglecounter.com).

It is also an excuse to experiment with technologies and concepts:
- **concurrency patterns**, eg. using channels to store values
- `sync` primitives like `atomic` values, `Cond` to syncronize events across goroutines
- observability setups using `VictoriaMetrics`, `VictoriaLogs` and `Grafana`
- `shadcn` components 
- `svelte` 5

## Acknowledgements

This project is inspired by [One Million Checkboxes](https://eieio.games/blog/one-million-checkboxes/) by [Eieio Games](https://eieio.games/).

## Running the project

### Dependencies 

#### Using `nix`

If you use `direnv`, just :

`direnv allow` 

otherwise, :

`nix develop`

#### Manually

You will need to install the following dependencies:
- `go` 1.23 or later
- `pnpm` 9 or later

And, optionally :
- [mprocs](https://github.com/pvolok/mprocs) to run all the processes
- [wgo](https://github.com/bokwoon95/wgo) to autoreload the server

### Running

The first time you run the project, you will have to install pnpm dependencies:

```
cd web
pnpm i
```

#### Using `mprocs`

Just:

`mprocs`

#### Manually

##### Server

```
cd server
go run cmd/counter/main.go
```

##### Web

```
cd web
pnpm run dev
```
