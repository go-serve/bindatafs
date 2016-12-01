package bindatafs_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/go-serve/bindatafs"
	"github.com/go-serve/bindatafs/examples/example1"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

func IndexFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello Index\n")
}

func Example() {

	// create vfs.FileSystem implementation for
	// the go-bindata generated assets
	assetsfs := bindatafs.New(
		"assets://",
		example1.Asset,
		example1.AssetDir,
		example1.AssetInfo,
	)

	// serve the files with http
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(httpfs.New(assetsfs))))
	mux.Handle("/", http.HandlerFunc(IndexFunc))

	// production: uncomment this
	//http.ListenAndServe(":8080", mux)

	// below are for testings, can be removed for production

	// test the mux with httptest server
	server := httptest.NewServer(mux)
	defer server.Close()

	// examine the index
	resp, _ := http.Get(server.URL)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)

	// examine an asset
	resp, _ = http.Get(server.URL + "/assets/hello.txt")
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)

	// Output:
	// Hello Index
	// Hello World
}
