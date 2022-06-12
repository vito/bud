package app

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/vfs"
)

func Load(fsys fs.FS, injector *di.Injector, module *gomod.Module, flag *framework.Flag) (*State, error) {
	if err := vfs.Exist(fsys, "bud/internal/app/web"); err != nil {
		return nil, err
	}
	return (&loader{
		fsys:     fsys,
		injector: injector,
		module:   module,
		flag:     flag,
		imports:  imports.New(),
	}).Load()
}

type loader struct {
	fsys     fs.FS
	injector *di.Injector
	module   *gomod.Module
	flag     *framework.Flag

	imports *imports.Set
	bail.Struct
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover2(&err, "app: unable to load state")
	state = new(State)
	state.Provider = l.loadProvider()
	l.imports.AddStd("os", "context")
	l.imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	l.imports.Add(l.module.Import("bud/internal/app/web"))
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadProvider() *di.Provider {
	jsVM := di.ToType("github.com/livebud/bud/package/js", "VM")
	fn := &di.Function{
		Name:   "loadWeb",
		Target: l.module.Import("bud", "program"),
		Params: []di.Dependency{
			di.ToType("github.com/livebud/bud/package/gomod", "*Module"),
			di.ToType("context", "Context"),
		},
		Results: []di.Dependency{
			di.ToType(l.module.Import("bud/internal/app/web"), "*Server"),
			&di.Error{},
		},
		Aliases: di.Aliases{
			di.ToType("io/fs", "FS"): di.ToType("github.com/livebud/bud/package/overlay", "*FileSystem"),
			jsVM:                     di.ToType("github.com/livebud/bud/package/js/v8client", "*Client"),
			di.ToType("github.com/livebud/bud/runtime/view", "Renderer"): di.ToType("github.com/livebud/bud/runtime/view", "*Server"),
		},
	}
	if l.flag.Embed {
		fn.Aliases[jsVM] = di.ToType("github.com/livebud/bud/package/js/v8", "*VM")
	}
	provider, err := l.injector.Wire(fn)
	if err != nil {
		// Intentionally don't wrap the error. The error gets swallowed up too
		// easily
		l.Bail(fmt.Errorf("app: unable to wire. %s", err))
	}
	// Add imports
	for _, im := range provider.Imports {
		l.imports.AddNamed(im.Name, im.Path)
	}
	return provider
}