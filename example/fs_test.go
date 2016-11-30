package example_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/tools/godoc/vfs/httpfs"

	"github.com/go-serve/bindatafs/example"
)

func assetContents(name string) (str string) {

	// get content of index file
	file, err := os.Open(name)
	if err != nil {
		panic(fmt.Sprintf("error opening %s, %s", name, err.Error()))
	}
	byteContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("error reading %s, %s", name, err.Error()))
	}
	str = string(byteContents)
	return
}

func TestFileSystem(t *testing.T) {
	fileSrvr := http.FileServer(httpfs.New(example.FileSystem()))

	tests := []struct {
		desc      string
		path      string
		code      int
		checkBody bool
		body      string
	}{
		{
			desc:      "getting normal text file hello.txt",
			path:      "/hello.txt",
			code:      http.StatusOK,
			checkBody: true,
			body:      assetContents("assets/hello.txt"),
		},
		{
			desc: "getting HTML file index.html",
			path: "/index.html",
			code: http.StatusMovedPermanently,
		},
		{
			desc:      "getting root directory",
			path:      "/",
			code:      http.StatusOK,
			checkBody: true,
			body:      assetContents("assets/index.html"),
		},
		{
			desc: "getting non-exists path",
			path: "/notfound.txt",
			code: http.StatusNotFound,
		},
		{
			desc: "getting directory index in hello",
			path: "/hello/",
			code: http.StatusOK,
		},
		{
			desc:      "getting normal text file hello/world.txt",
			path:      "/hello/world.txt",
			code:      http.StatusOK,
			checkBody: true,
			body:      assetContents("assets/hello/world.txt"),
		},
		{
			desc: "getting directory non-exists sub-directory path",
			path: "/hello/foo.txt",
			code: http.StatusNotFound,
		},
	}

	for i, test := range tests {

		t.Logf("test %d: %s", i, test.desc)

		// test getting text file asset
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", test.path, nil)
		if err != nil {
			t.Errorf("error reading %s, %s", test.path, err.Error())
			return
		}
		fileSrvr.ServeHTTP(w, r)
		if want, have := test.code, w.Code; want != have {
			t.Errorf("expected: %#v, got %#v", want, have)
		}
		if test.checkBody {
			if want, have := test.body, w.Body.String(); want != have {
				t.Errorf("\nexpected: %s\ngot:      %s", want, have)
			}
		}
	}

}
