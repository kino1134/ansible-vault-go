# ビルド可能なクロスコンパイル環境は、以下で表示可能
# go tool dit list [-json]

# 自環境
go build -o bin/vault_go

# Window x64
GOOS=windows GOARCH=amd64 go build -o bin/vault_go.exe
