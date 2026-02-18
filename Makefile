PROTO_DIR := ./proto/admin
OUT_DIR   := ./proto/admin
# Busca automaticamente todos os arquivos .proto na pasta
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Garante que o PATH inclua os binários do Go
export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: proto clean

# Comando principal: make proto
proto: $(PROTO_FILES)
	@echo "Gerando arquivos Go a partir dos protos..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Concluído."