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
	cache, err := wazero.NewCompilationCacheWithDir(".")
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
	err = mfs.MkdirAll("output", 0)
	if err != nil {
		log.Fatalf("error creating the output directory: %v", err)
	}
	err = mfs.WriteFile("output/test.txt", []byte("test"), 0)
	if err != nil {
		log.Fatalf("error creating the test file: %v", err)
	}
	files, err := mfs.Glob("*/*")
	if err != nil {
		log.Fatalf("error glob: %v", err)
	}
	for _, name := range files {
		log.Printf("go fs file: %s", name)
	}

	moduleConfig := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin).
		WithRandSource(rand.Reader).
		WithFSConfig(wazero.NewFSConfig().
			// TODO even thou we have a working FS, why isn't the python script seeing this FSMount?
			WithFSMount(mfs, "/").
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
