# Alfred Workflows 根目录 Makefile
# 这个文件会自动转发所有命令到 src 目录下的 Makefile

# 获取命令行参数
ARGS = $(filter-out $@,$(MAKECMDGOALS))

# 默认目标：显示帮助信息
.PHONY: default
default: help

# 传递所有目标到 src 目录下的 Makefile
%:
	cd src && $(MAKE) $@ $(ARGS)
