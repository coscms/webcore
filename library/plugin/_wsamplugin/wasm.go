package wasmplugin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/wasmerio/wasmer-go/wasmer"
	"github.com/webx-top/echo"
)

func Load(files ...string) (err error) {
	if len(files) == 0 {
		pluginGlob := filepath.Join(echo.Wd(), `plugins`) + echo.FilePathSeparator + `*.wasm`
		files, err = filepath.Glob(pluginGlob)
		if err != nil {
			return
		}
	}
	for _, file := range files {
		im := NewWASMImports()
		var w *WASM
		log.Infof(`load plugin: %s`, file)
		w, err = LoadWASMFile(file, im)
		if err != nil {
			return
		}
		if start, sterr := w.Exports.GetWasiStartFunction(); sterr == nil {
			if _, err = start(); err != nil {
				return
			}
		} else if !strings.Contains(sterr.Error(), `WASI start function was not found`) {
			return sterr
		}
	}
	return
}

func LoadWASMFile(wasmFile string, imports *WASMImports) (*WASM, error) {
	w := &WASM{file: wasmFile}
	wasmBytes, err := os.ReadFile(wasmFile)
	if err != nil {
		return nil, err
	}
	w.engine = wasmer.NewEngine()
	w.store = wasmer.NewStore(w.engine)

	// Compiles the module
	w.module, err = wasmer.NewModule(w.store, wasmBytes)
	if err != nil {
		return nil, err
	}
	var wasiEnv *wasmer.WasiEnvironment
	wasiEnv, err = wasmer.NewWasiStateBuilder("wasi-program").
		// Choose according to your actual situation
		// Argument("--foo").
		// Environment("ABC", "DEF").
		// MapDirectory("./", ".").
		Finalize()
	if err != nil {
		return nil, err
	}
	// Instantiates the module
	if wasiEnv != nil {
		w.importObject, err = wasiEnv.GenerateImportObject(w.store, w.module)
		if err != nil {
			return nil, err
		}
	} else {
		w.importObject = wasmer.NewImportObject()
	}
	if imports != nil {
		for namespace, externCreators := range *imports {
			if externCreators == nil {
				continue
			}
			extern := map[string]wasmer.IntoExtern{}
			for name, create := range *externCreators {
				extern[name] = create(w.store)
			}
			w.importObject.Register(namespace, extern)
		}
	}
	w.Instance, err = wasmer.NewInstance(w.module, w.importObject)

	// Gets the `sum` exported function from the WebAssembly instance.
	// sum, _ := w.Instance.Exports.GetFunction("sum")

	return w, err
}

type IntoExternCreator func(*wasmer.Store) wasmer.IntoExtern

type WASMIECs map[string]IntoExternCreator

func (w *WASMIECs) Add(name string, creator IntoExternCreator) *WASMIECs {
	(*w)[name] = creator
	return w
}

func (w *WASMIECs) Del(names ...string) *WASMIECs {
	for _, name := range names {
		delete(*w, name)
	}
	return w
}

func NewWASMImports() *WASMImports {
	return &WASMImports{}
}

type WASMImports map[string]*WASMIECs

func (w *WASMImports) Add(ns string, name string, creator IntoExternCreator) *WASMImports {
	if _, ok := (*w)[ns]; !ok {
		(*w)[ns] = &WASMIECs{}
	}
	(*w)[ns].Add(name, creator)
	return w
}

func (w *WASMImports) AddItems(ns string, creators *WASMIECs) *WASMImports {
	(*w)[ns] = creators
	return w
}

func (w *WASMImports) Del(namespaces ...string) *WASMImports {
	for _, ns := range namespaces {
		delete(*w, ns)
	}
	return w
}

func (w *WASMImports) DelCreator(ns string, names ...string) *WASMImports {
	if _, ok := (*w)[ns]; ok {
		(*w)[ns].Del(names...)
	}
	return w
}

type WASM struct {
	file         string
	engine       *wasmer.Engine
	store        *wasmer.Store
	module       *wasmer.Module
	importObject *wasmer.ImportObject
	*wasmer.Instance
}

func (w *WASM) Close() {
	w.Instance.Close()
	w.module.Close()
	w.store.Close()
}
