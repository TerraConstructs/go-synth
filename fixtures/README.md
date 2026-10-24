# Fixtures

Sample configurations to test `go-synth` against.

## local-package

Sample `main.ts` referencing a local library (used in executor tests)

## TerraConstructs/Base

Sample `main.ts` creating an AWS Network with [terraconstructs](https://www.npmjs.com/package/terraconstructs) (formerly `@envtio/base`)

```console
synth -file main.ts -deps "terraconstructs:0.2.19" -src "cdktf.out/stacks/sample-stack"
```
