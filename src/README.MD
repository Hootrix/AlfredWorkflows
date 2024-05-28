

# 打包

go build -ldflags='-s -w -extldflags "-static -fpic"'  -o timestamp-plus.bin