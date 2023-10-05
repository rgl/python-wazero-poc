run: example
	./example

example: python.wasm lib *.py *.go go.*
	go build

python.wasm lib:
	rm -rf python.wasm lib python-3.12.0-wasi_sdk-20.zip
	wget https://github.com/brettcannon/cpython-wasi-build/releases/download/v3.12.0/python-3.12.0-wasi_sdk-20.zip
	unzip python-3.12.0-wasi_sdk-20.zip

docker-build:
	docker build --progress=plain --tag=python-wazero .

docker-build-no-cache:
	docker build --no-cache --progress=plain .

docker-run: docker-build
	docker run python-wazero

clean:
	rm -rf example python*wasi*.zip python.wasm lib wazero-*
