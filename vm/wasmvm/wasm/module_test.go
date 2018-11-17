package wasm_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"dan-road-vbft/vm/wasmvm/wasm"
)

func TestReadModule(t *testing.T) {
	fnames, err := filepath.Glob(filepath.Join("testdata", "*.wasm"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fname := range fnames {
		name := fname
		t.Run(filepath.Base(name), func(t *testing.T) {
			raw, err := ioutil.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}

			r := bytes.NewReader(raw)
			m, err := wasm.ReadModule(r, nil)
			if err != nil {
				t.Fatalf("error reading module %v", err)
			}
			if m == nil {
				t.Fatalf("error reading module: (nil *Module)")
			}
		})
	}
}
