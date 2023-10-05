package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/liamg/memoryfs"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed python.wasm
var python []byte

//go:embed main.py
var script string

const pythonLibRelativePath = "lib/python3.12"

func main() {
	// ensure the python library exists.
	pythonLibPath, err := filepath.Abs(pythonLibRelativePath)
	if err != nil {
		log.Fatalf("error finding python library: %v", err)
	}
	if stat, err := os.Stat(pythonLibPath); err != nil {
		log.Fatalf("error finding python library: %v", err)
	} else if !stat.IsDir() {
		log.Fatalf("the python library path is not a directory: %v", err)
	}

	ctx := context.Background()

	// NB unfortunately, there is no way to have an embedded read-only cache.
	//    see https://github.com/tetratelabs/wazero/issues/1733
	// NB unfortunately, this also means that in the long-run, it will contain
	//    stalled data. so we should probably disable this entirely, and let it
	//    compile every time.
	cache, err := wazero.NewCompilationCacheWithDir(".cache")
	if err != nil {
		log.Fatalf("error creating the compilation cache: %v", err)
	}
	defer cache.Close(ctx)

	r := wazero.NewRuntimeWithConfig(ctx,
		wazero.
			NewRuntimeConfigCompiler().
			WithCompilationCache(cache))
	defer r.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	module, err := r.CompileModule(ctx, python)
	if err != nil {
		log.Fatalf("error compiling python: %v", err)
	}

	mfs := memoryfs.New()
	err = mfs.WriteFile("from-host.txt", []byte("from host"), 0)
	if err != nil {
		log.Fatalf("error creating the test file: %v", err)
	}
	files, err := mfs.Glob("*")
	if err != nil {
		log.Fatalf("error glob: %v", err)
	}
	for _, name := range files {
		log.Printf("go fs file: %s", name)
	}

	// run the python script.
	// NB this is equivalent to:
	//		./wazero \
	// 			run \
	// 			-cachedir=.cache \
	// 			-mount=$PWD/lib/python3.12:/usr/local/lib/python3.12:ro \
	// 			-mount=$PWD/output:/output \
	// 			python.wasm \
	// 			-- \
	// 			-c "$(cat main.py)"
	// see https://github.com/tetratelabs/wazero/blob/v1.5.0/cmd/wazero/wazero.go
	moduleConfig := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithFSConfig(wazero.NewFSConfig().
			// TODO even thou we have a working FS, why isn't the python script seeing this FSMount?
			WithFSMount(mfs, "/output").
			WithReadOnlyDirMount(pythonLibPath, fmt.Sprintf("/usr/local/%s", pythonLibRelativePath))).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs("python", "-c", script)

	_, err = r.InstantiateModule(ctx, module, moduleConfig)
	if err != nil {
		log.Panicf("failed to run the python script: %v", err)
	}
}
