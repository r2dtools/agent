legoVersion=4.18.0
legoArchive=lego.tar.gz

test:
	docker run --volume="$(shell pwd):/opt/r2dtools" agent-tests

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

build_test:
	docker build -t agent-tests . 
clean:
	cd build; \
	rm -rf config; \
	rm -f lego r2dtools LICENSE .version

serve:
	go run cmd/main.go serve

.PHONY: test
