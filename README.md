## What

**Scrubeer** is a tool to scrub potential PII from a Java heapdump. It works by replacing **all** primitive arrays with
either `x` or `0.0` for floats and doubles.

Static strings are kept, which basically means stack frames.

## Why

To remove PII from heapdumps, but still allow for developers to perform analysis.

## Usage

```
Usage:
  scrubeer [OPTIONS]

Application Options:
  -i, --input-file=
  -o, --output-file=
  -k, --keep=

Help Options:
  -h, --help         Show this help message
```

It's possible to keep specific types by passing a `-k` / `--keep` flag, which will keep the scrubber from operating on
specific types. The recognized types are:

- `bool` or `boolean`
- `char`
- `float`
- `double`
- `byte`
- `short`
- `int`
- `long`

These types can be passed individually (`-k char -k float`) or comma separated (`-k char,float`)

## How

If you're not setup to build golang code, you can build and run the docker image:

```shell
make docker
docker run --rm -it -v $(pwd):/work scrubeer:local -i input.hprof -o output.hprof
```

## Apologies

The code isn't organized great, and error handling is for other people. It can be improved if someone needs it, but
since this is largely a manual process anyway, it's probably acceptable to not handle things cleanly.
