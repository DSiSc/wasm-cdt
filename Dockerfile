FROM gcc
## 设置工作目录
WORKDIR /workspace
## 把当前（宿主机上）目录下的文件都复制到docker上刚创建的目录下
COPY ./llvm/bin/ /usr/local/bin/
COPY ./llvm/lib/ /usr/local/lib/
COPY ./wasm-cdt /usr/local/wasm-cdt/wasm-cdt
COPY ./misc/	/usr/local/wasm-cdt/misc/
## 启动需要执行的文件
ENTRYPOINT ["/usr/local/wasm-cdt/wasm-cdt"]
CMD ["--help"]
