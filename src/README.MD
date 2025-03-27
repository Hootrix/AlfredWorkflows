

# 打包方式

## 方式一：直接使用 go build

```bash
cd ./src/timestamp_plus
go build -ldflags='-s -w -extldflags "-static -fpic"'  -o timestamp-plus.bin
```

## 方式二：使用 make

```bash
# 打包所有二进制文件
make all

# 或者单独打包特定模块
make timestamp_plus
make translate
make code
```