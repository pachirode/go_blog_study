# ==============================================================================
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
OUTPUT_DIR := $(ROOT_DIR)/_output
# ==============================================================================
.PHONY: build, format, tidy, clean

build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -o $(OUTPUT_DIR)/miniblog $(ROOT_DIR)/cmd/miniblog/main.go

format: # 格式化 Go 源码.
	@gofmt -s -w ./

tidy: # 自动添加/移除依赖包.
	@go mod tidy

clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)
