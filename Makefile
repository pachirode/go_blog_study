# 设置默认目标为 all
.DEFAULT_GOAL := all

# make 命令默认执行
all: tidy format build

include scripts/make-rules/all.mk

# ==============================================================================
# Usage
# ==============================================================================

define USAGE_OPTIONS

选项:

  BINS				要构建的二进制文件，默认为 cmd 文件夹下面的所有文件
					示例：make build BINS="miniblog"
  VERSION			编译到二进制文件中的版本信息
  V					设置为 1 启动纤细构建信息输出，默认值为 0

 endef
export USAGE_OPTIONS

# ==============================================================================
# Binaries
# ==============================================================================

build: go.tidy
	@$(MAKE) go.build

# ==============================================================================
# Testing
# ==============================================================================

test:
	@$(MAKE) go.test

cover:
	@$(MAKE) go.cover

# ==============================================================================
# Cleanup
# ==============================================================================

clean:
	@echo "==> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)


# ==============================================================================
# Lint / verification
# ==============================================================================

lint:
	@$(MAKE) go.lint

tidy:
	@$(MAKE) go.tidy

format: tools.verify.protolint
	@$(MAKE) go.format
	@protolint -fix -config_path ${PROJ_ROOT_DIR}/.protolint.yaml $(shell find $(APIROOT) -name *.proto)
