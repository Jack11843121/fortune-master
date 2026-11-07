APP := fortune

.PHONY: build run test clean build-web build-linux build-windows

build-web:
	cd web && npm install && npm run build

build: build-web
	go build -o $(APP) ./cmd/server

run: build
	./$(APP) -config config.yaml

test:
	go test ./...

# 交叉编译到 Linux（部署用）
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP)-linux ./cmd/server

# 交叉编译 Windows（本机就用 build）
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(APP).exe ./cmd/server

clean:
	rm -f $(APP) $(APP)-linux $(APP).exe