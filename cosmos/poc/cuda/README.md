# Cuda support

## Development 
To execute gpu and cpu rank computing ran:

```bash
nvcc -fmad=false -shared -o librank.so rank.cu --compiler-options '-fPIC -frounding-math -fsignaling-nans'
sudo cp librank.so /usr/lib/
go run *.go
```
After executing check ranks. They should match.


While creating the shared libraries, position independent code should be produced. This helps the shared library 
 to get loaded as any address instead of some fixed address. For this `-fPIC` option is used.

## Determinism

//fma is disabled for current version.

https://gcc.gnu.org/wiki/FloatingPointMath
'-frounding-math' is round-to-zero for all floating point to integer conversions, and round-to-nearest for 
all other arithmetic truncations.  
 
 