# Alfred Workflows 执行脚本构建

# 构建参数
LDFLAGS=-ldflags='-s -w -extldflags "-static -fpic"'

# 所有目标
.PHONY: all timestamp_plus translate code raycast_timestamp clean help

# 默认目标：构建所有二进制文件
all: timestamp_plus translate code

# 构建 timestamp_plus 模块
timestamp_plus:
	cd cmd/timestamp_plus && go build $(LDFLAGS) -o ../../bin/timestamp-plus.bin

# 构建 translate 模块
translate:
	cd cmd/translate && go build $(LDFLAGS) -o ../../bin/translate.bin

# 构建 code 模块
code:
	cd cmd/code && go build $(LDFLAGS) -o ../../bin/code.bin

# 清理所有生成的二进制文件
clean:
	rm -f bin/timestamp-plus.bin
	rm -f bin/translate.bin
	rm -f bin/code.bin

# 安装并启动 Raycast timestamp_plus 插件
raycast_timestamp:
	cd raycast/timestamp_plus && npm install && npm run dev

# 显示帮助信息
help:
	@echo "可用的 make 命令："
	@echo "  make all             - 构建所有模块"
	@echo "  make timestamp_plus  - 只构建 timestamp_plus 模块"
	@echo "  make translate       - 只构建 translate 模块"
	@echo "  make code            - 只构建 code 模块"
	@echo "  make raycast_timestamp - 安装并启动 Raycast timestamp_plus 插件"
	@echo "  make clean           - 清理所有生成的二进制文件"
	@echo "  make help            - 显示此帮助信息"
