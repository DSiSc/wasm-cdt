# wasm-cdt

justitita wasm contract development kit

[![Build Status](https://circleci.com/gh/DSiSc/wasm-cdt/tree/master.svg?style=shield)](https://circleci.com/gh/DSiSc/wasm-cdt/tree/master)
[![codecov](https://codecov.io/gh/DSiSc/wasm-cdt/branch/master/graph/badge.svg)](https://codecov.io/gh/DSiSc/wasm-cdt)

## Getting started

Running it then should be as simple as:

```
$ make all
```

### Testing

```
$ make test
```

### Build Docker Image
We use `docker` so that `wasm-cdt` command can run on any type system. You can compile from the source code or directly use [pre-build binary](https://github.com/DSiSc/wasm-cdt/releases/download/v0.1/wasm-build.7z) to build docker image.

1. Make build directory, structure:
    ```
     -\
       -llvm\
       -misc\
       -Dockerfile
       -wasm-cdt
    ```
    
    [llvm](https://github.com/DSiSc/wasm-cdt/tree/master/llvm): llvm binary files
    
    [misc](https://github.com/DSiSc/wasm-cdt/tree/master/misc): base header files
    
    [Dockerfile](https://github.com/DSiSc/wasm-cdt/blob/master/Dockerfile): Docker build file
    
    wasm-cdt: go build binary file.

2. Start Building：
    ```
     $ docker build -t wasm-cdt .
     $ docker run –rm -it wasm-cdt --help (verify your image)
    ```

## Write Wasm Contract

1. For your `C` file must have `invoke` method, and `invoke` is the only entry point of the contract. `invoke` method format:
    
    ```
    $ char *invoke(int argc, char *argv[])
    ```

2. You can use [justitia.h](https://github.com/DSiSc/wasm-cdt/blob/master/misc/include/justitia/justitia.h) in your contract to get info relate to chain:
    
    ```
    #include <stdio.h>
    #include <string.h>
    #include "justitia.h"
    ```

3. [demo.c](https://github.com/DSiSc/wasm-cdt/blob/master/example/contract.c) is our demo contract, you get some detail to write wasm contract.

4. To deploy a contract, you can use our `deploy` command.
    
    ```
    docker run --rm -it -v ${local_workspace}:/workspace wasm-cdt deploy -f ${wasm_file} -ac ${deploy_account} -ep ${endpoint_address}
    ```

5. To call wasm contract, you can use our rpc api [wasm-call](https://dsisc.github.io/slate/#wasmcall) and [wasm-transaction](https://dsisc.github.io/slate/#sendwasmtransaction).