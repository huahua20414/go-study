# 你可以直接执行 make 命令，也可以单独的命令复制到控制台。
# 注意，如果你是 Windows 并且不是在 WSL 下，要使用gitbash来执行make docker
# 要注意文件分隔符使用 Windows 的分隔符。
.PHONY: docker
docker:
	# remove file webook if exits
	@rm webook || true
	# Delete the existing mirrors.
	@#docker rmi -f huahua/webook:v0.0.4 || true
	# Create a linux executable file
	@GOOS="linux" GOARCH="arm" go build -tags=k8s -o webook .
	# Pull this image through the Dockerfile configuration
	@docker build --no-cache -t huahua/webook:v0.0.1 .
