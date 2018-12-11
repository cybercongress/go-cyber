# Cuda support

## Running with GPU rank calculation
To compile cyberd with cuda support first install latest [cuda toolkit](https://developer.nvidia.com/cuda-downloads).

Next compile **cbdrank** lib, copy it to `/usr/lib/` folder:

```bash
# project root
cd x/rank/cuda
nvcc -fmad=false -shared -o libcbdrank.so rank.cu --compiler-options '-fPIC -frounding-math -fsignaling-nans'
sudo cp libcbdrank.so /usr/lib/
sudo cp cbdrank.h /usr/lib/
```

Compile binaries, copy configs and run daemon
```bash
# project root
cp testnet/genesis.json .cyberd/config/genesis.json
cp testnet/config.toml .cyberd/config/config.toml
go build -tags cuda -o daemon ./cyberd
./cyberd
```

## Testing 
To test GPU and CPU rank computing determinism run:
```bash
# project root
cd x/rank/cuda
nvcc -fmad=false -shared -o libcbdrank.so rank.cu --compiler-options '-fPIC -frounding-math -fsignaling-nans'
sudo cp libcbdrank.so /usr/lib/
sudo cp cbdrank.h /usr/lib/
go build -tags cuda -o test *.go && ./test && rm test
```
After executing check ranks. They should match.

## Compilation args

While creating the shared libraries, position independent code should be produced. This helps the shared library 
 to get loaded as any address instead of some fixed address. For this `-fPIC` option is used.

See https://gcc.gnu.org/wiki/FloatingPointMath. '-frounding-math' is round-to-zero for all floating point to integer
 conversions, and round-to-nearest for all other arithmetic truncations. FMA is disabled for current version. 
 
 