package uds_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/internal/uds"

	"github.com/matryer/is"
)

func TestListenTransport(t *testing.T) {
	is := is.New(t)
	socketPath := fmt.Sprintf("%s.sock", strings.ToLower(t.Name()))
	is.NoErr(os.RemoveAll(socketPath))
	listener, err := uds.Listen(socketPath)
	is.NoErr(err)
	defer listener.Close()
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi!"))
		}),
	}
	go func() {
		err := server.Serve(listener)
		if !errors.Is(err, http.ErrServerClosed) {
			is.NoErr(err)
		}
	}()
	client := &http.Client{Transport: uds.Transport(socketPath)}
	res, err := client.Get("http://unix/")
	is.NoErr(err)
	body, err := ioutil.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "hi!")
	is.NoErr(server.Close())
}