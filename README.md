# pede-lang

~~Simple. Fast. Easy to use.~~

Yet another dummy programming language started for the sake of going through the hassle of building a language from scratch and obtain knowledge. 

## Build from source
- `go` to build the pede compiler.
- `make` for simplicity.

```bash
make build
make run IN=examples/hello.pede OUT=hello
./hello
```

## Prerequisites
- `gcc` or `clang` required at build time. pede will use `clang` by default to link generated code.

## Usage

```bash
./pede build examples/hello.pede
./hello

./pede --log=DEBUG build -o arithmetics --os=darwin --arch=arm64 --cc=clang --keep-ir \
  examples/arithmetics.pede
./arithmetics
```

## Examples

```pede
# examples/hello.pede
hello = "Hello, World!"
print(hello)
```

```pede
# examples/arithmetics.pede
x = 3 + 4 * 2
y = 5
print("x = 3 + 4 * 2:")
print(x)
print("y = 5:")
print(y)
print("x + y")
print(x + y)
```
