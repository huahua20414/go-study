# 你可以直接执行 make 命令，也可以单独的命令复制到控制台。
# 注意，如果你是 Windows 并且不是在 WSL 下，要使用gitbash来执行make docker
# 要注意文件分隔符使用 Windows 的分隔符。
.PHONY: docker
docker:
	# remove file day01 if exits
	@rm day01 || true
	# Delete the existing mirrors.
	@docker rmi -f huahua/day01:v0.0.2
	# Create a linux executable file
	@GOOS="linux" GOARCH="arm" go build -o day01 .
	# Pull this image through the Dockerfile configuration
	@docker build --no-cache -t huahua/day01:v0.0.2 .
