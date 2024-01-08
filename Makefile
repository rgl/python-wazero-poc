SHELL=bash
.SHELLFLAGS=-euo pipefail -c

run: example
	./example

run-wazero: wazero
	mkdir -p output && echo 'from host'>output/from-host.txt
	./wazero \
		run \
		-cachedir=.cache \
		-mount=$$PWD/lib/python3.12:/usr/local/lib/python3.12:ro \
		-mount=$$PWD/output:/output \
		python.wasm \
		-- \
		<main.py

example: python.wasm lib *.py *.go go.*
	go build

python.wasm lib:
	rm -rf python.wasm lib python-3.12.1-wasi_sdk-20.zip
	wget https://github.com/brettcannon/cpython-wasi-build/releases/download/v3.12.1/python-3.12.1-wasi_sdk-20.zip
	unzip python-3.12.1-wasi_sdk-20.zip

wazero:
	wget https://github.com/tetratelabs/wazero/releases/download/v1.5.0/wazero_1.5.0_linux_amd64.tar.gz
	tar xf wazero_1.5.0_linux_amd64.tar.gz wazero

docker-build:
	docker build --progress=plain --tag=python-wazero .

docker-build-no-cache:
	docker build --no-cache --progress=plain .

docker-run: docker-build
	docker run python-wazero

clean:
	rm -rf example python*wasi*.zip python.wasm lib .cache wazero*
