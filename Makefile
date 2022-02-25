hypixel-cache: *.go
	go build -ldflags "-X main.mode=release" -o hypixel-cache
