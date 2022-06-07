package public_test

import (
	"bytes"
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/public"
	"github.com/livebud/bud/internal/golden"
	"github.com/livebud/bud/internal/is"
)

type Flag = framework.Flag

func TestEmpty(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	state, err := public.Load(fsys, &Flag{
		Embed: false,
	})
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(state, nil)
}

func check(actual []byte, expected []byte) ([]byte, error) {
	if !bytes.Equal(actual, expected) {
		return nil, errors.New("bytes not equal")
	}
	return []byte{}, nil
}

func TestLinkEmptyPublic(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"public": &fstest.MapFile{Mode: fs.ModeDir},
	}
	state, err := public.Load(fsys, &Flag{
		Embed: false,
	})
	is.NoErr(err)
	golden.State(t, state)
}

func TestEmbedEmptyPublic(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"public": &fstest.MapFile{Mode: fs.ModeDir},
	}
	state, err := public.Load(fsys, &Flag{
		Embed: true,
	})
	is.NoErr(err)
	golden.State(t, state)
}
