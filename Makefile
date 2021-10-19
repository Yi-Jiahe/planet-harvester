build:
	[ -d bin ] || mkdir bin
	cd bin; \
	go build ../src/server/websocket