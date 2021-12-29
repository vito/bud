package program

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/di"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

type Generator struct {
	Module   *mod.Module
	Injector *di.Injector
}

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	if err := gen.SkipUnless(g.Module, "bud/command/command.go"); err != nil {
		return err
	}
	// Add the imports
	imports := imports.New()
	imports.AddStd("os", "errors", "context", "runtime", "path/filepath")
	// imports.AddStd("fmt")
	imports.AddNamed("console", "gitlab.com/mnm/bud/log/console")
	imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")
	imports.AddNamed("plugin", "gitlab.com/mnm/bud/plugin")
	imports.AddNamed("mod", "gitlab.com/mnm/bud/go/mod")
	imports.Add(g.Module.Import("bud/command"))
	// Inject command
	provider, err := g.Injector.Wire(&di.Function{
		Name:   "loadCLI",
		Target: g.Module.Import("bud", "program"),
		Params: []di.Dependency{
			&di.Type{Import: "gitlab.com/mnm/bud/go/mod", Type: "*Module"},
			&di.Type{Import: "gitlab.com/mnm/bud/gen", Type: "*FileSystem"},
		},
		Results: []di.Dependency{
			&di.Type{Import: g.Module.Import("bud", "command"), Type: "*CLI"},
			&di.Error{},
		},
	})
	if err != nil {
		return err
	}
	for _, im := range provider.Imports {
		imports.AddNamed(im.Name, im.Path)
	}
	code, err := generator.Generate(State{
		Imports:  imports.List(),
		Provider: provider,
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}