legoVersion=4.17.4
legoArchive=lego.tar.gz

test:
	go test ./...

build_agent:
	go build -o ./build/r2dtools -v cmd/main.go

build_lego:
	wget "https://github.com/go-acme/lego/releases/download/v${legoVersion}/lego_v$(legoVersion)_linux_amd64.tar.gz" -O $(legoArchive); \
	tar -xvzf $(legoArchive) -C build lego; \
	rm $(legoArchive)

build: build_agent build_lego
	mkdir -p build/config; \
    cp .version LICENSE build/; \
    cp config/params.yaml build/config/

clean:
	cd build; \
	rm -rf config; \
	rm -f lego r2dtools LICENSE .version

serve:
	go run cmd/main.go serve
