# Cuda support

## Install required libs
```bash

```

## Development 
To execute cuda code run.

```bash
nvcc -shared -o librank.so rank.cu --compiler-options '-fPIC -frounding-math -fsignaling-nans' && sudo cp librank.so /usr/lib/
go run *.go
```


https://gcc.gnu.org/wiki/FloatingPointMath

While creating the shared libraries, position independent code should be produced. This helps the shared library 
 to get loaded as any address instead of some fixed address. For this -fPIC option is used.
 
'-frounding-math' is round-to-zero for all floating point to integer conversions, and round-to-nearest for all other arithmetic truncations.  
 
 