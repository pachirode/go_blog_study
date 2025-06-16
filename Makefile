# ==============================================================================
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
OUTPUT_DIR := $(ROOT_DIR)/_output

VERSION_PACKAGE=github.com/marmotedu/Miniblog/pkg/version
ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git describe --tags --always --match='v*')
endif

GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
	-X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# ==============================================================================
.PHONY: build, format, tidy, clean

build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -ldflags '$(GO_LDFLAGS)' -o $(OUTPUT_DIR)/miniblog $(ROOT_DIR)/cmd/miniblog/main.go

format: # 格式化 Go 源码.
	@gofmt -s -w ./

tidy: # 自动添加/移除依赖包.
	@go mod tidy

clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)
