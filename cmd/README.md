# synth CLI

Sample CLI using `go-synth` library with [bun.sh](https://bun.sh) executor.

## Tool Versions

Using [mise](https://mise.jdx.dev/getting-started.html#quickstart).

```console
# view requested config
mise ls

# install tool versions required
mise install
```

Synth `example/network.ts` with `bun`:

```console
go build -o synth
./synth -file example/network.ts -deps "terraconstructs:0.0.8" -src "cdktf.out/stacks/network-stack" -out result/network
```

Will create directory `result/` containing all the synthesized files.
