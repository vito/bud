package framework

// Flag is used by many of the framework generators
type Flag struct {
	Embed  bool
	Minify bool
	Hot    string // Hot reload address
}