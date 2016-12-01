package bindatafs_test

import (
	"fmt"
	"net/http"

	"github.com/go-serve/bindatafs"
	"github.com/go-serve/bindatafs/examples/example1"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

func exampleIndex(w http.ResponseWriter, r *http.Request) {
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
	mux.Handle("/", http.HandlerFunc(exampleIndex))

	// serve the mux
	http.ListenAndServe(":8080", mux)
}
