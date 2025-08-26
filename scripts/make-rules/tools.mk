# ==============================================================================
# 工具相关
# ==============================================================================

TOOLS ?= protoc-plugins protoc-go-inject-tag protolint

tools.verify: $(addprefix tools.verify., $(TOOLS))

tools.install: $(addprefix tools.install., $(TOOLS))

tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

install.protoc-plugins:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0

install.protoc-go-inject-tag:
	@$(GO) install github.com/favadi/protoc-go-inject-tag@latest

install.protolint:
	@$(GO) install github.com/yoheimuta/protolint/cmd/protolint@latest

.PHONY: tools.verify tools.install tools.install.% tools.verify.% \
		install.protoc-plugins install.protoc-go-inject-tag install.protolint
