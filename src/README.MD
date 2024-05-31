

# 打包 e.g.

cd ./src/timestamp-plus
go build -ldflags='-s -w -extldflags "-static -fpic"'  -o timestamp-plus.bin