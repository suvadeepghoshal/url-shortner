run:
	@go run cmd/main.go

build:
	@go build -o bin/app cmd/main.go
# 	need to add @upx https://upx.github.io
	bin/app
	@echo "compiled you application with all its assets to a single binary => bin/app"