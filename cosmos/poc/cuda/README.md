# Cuda support

## Install required libs
```bash

```

## Development 
To execute cuda code run.

```bash
nvcc -shared -o librank.so rank.cu --compiler-options '-fPIC' && sudo cp librank.so /usr/lib/
go run *.go
```