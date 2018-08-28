package bindatafs_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/tools/godoc/vfs"

	"github.com/go-serve/bindatafs"
	"github.com/go-serve/bindatafs/examples/example1"
)

const ASSETS_PATH = "./examples/example1/assets/"

func TestFileSystem(t *testing.T) {
	var vfsFS vfs.FileSystem = bindatafs.New("assets://", nil, nil, nil)
	_ = vfsFS // just to prove bindatafs.FileSystem implements http.FileSystem
}

func msgNotFound(op, pathname string) string {
	return fmt.Sprintf("%s %s: no such file or directory", op, pathname)
}

func fileInfoEqual(src, target os.FileInfo) (err error) {
	if want, have := src.Name(), target.Name(); want != have {
		err = fmt.Errorf("Name(): expected %#v, got %#v", want, have)
		return
	}
	if want, have := src.IsDir(), target.IsDir(); want != have {
		err = fmt.Errorf("IsDir(): expected %#v, got %#v", want, have)
		return
	}
	if src.IsDir() {
		if want, have := int64(0), target.Size(); want != have {
			err = fmt.Errorf("Size(): expected %#v, got %#v", want, have)
			return
		}
		if want, have := os.ModeDir, target.Mode()&os.ModeType; want != have {
			err = fmt.Errorf("Mode():\nexpected %b\ngot      %b", want, have)
			return
		}
		if want, have := os.FileMode(0777), target.Mode()&os.ModePerm; want != have {
			err = fmt.Errorf("Mode():\nexpected %b\ngot      %b", want, have)
			return
		}
		if want, have := int64(0), target.ModTime().Unix(); want != have {
			err = fmt.Errorf("Modtime(): expected %#v, got %#v", want, have)
			return
		}
	} else {
		if want, have := src.Size(), target.Size(); want != have {
			err = fmt.Errorf("Size(): expected %#v, got %#v", want, have)
			return
		}
		if want, have := os.FileMode(0444), target.Mode()&os.ModePerm; want != have {
			err = fmt.Errorf("Mode():\nexpected %b\ngot      %b", want, have)
			return
		}
		if want, have := src.ModTime().Unix(), target.ModTime().Unix(); want != have {
			err = fmt.Errorf("Modtime(): expected %#v, got %#v", want, have)
			return
		}
	}
	return
}

func TestFileSystem_Open(t *testing.T) {
	fs := example1.FileSystem()
	tests := []struct {
		desc string
		path string
		err  string
	}{
		{
			desc: "test open file",
			path: "hello.txt",
		},
		{
			desc: "test open sub-directory file",
			path: "hello/world.txt",
		},
		{
			desc: "test open directory",
			path: "hello",
			err:  msgNotFound("Open", "hello"),
		},
		{
			desc: "test open non-exists path",
			path: "notfound",
			err:  msgNotFound("Open", "notfound"),
		},
	}

	for i, test := range tests {
		t.Logf("test fs.Open %d: %s", i+1, test.desc)

		// get the file/dir in the bindatafs
		file, err := fs.Open(test.path)
		if test.err != "" {
			if err == nil {
				t.Errorf("expected error %#v, got nil", test.err)
			} else if want, have := test.err, err.Error(); want != have {
				t.Errorf("expected error %#v, got %#v", want, have)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		defer file.Close()

		// get the counter part in the source assets
		srcFile, err := os.Open(ASSETS_PATH + test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		defer srcFile.Close()
		srcFileBytes, err := ioutil.ReadAll(srcFile)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}

		if want, have := string(srcFileBytes), string(fileBytes); want != have {
			t.Errorf("unexpected content for %#v", test.path)
			t.Logf("expected:\n%s\ngot:\n%s", want, have)
		}

	}
}

func TestFileSystem_Stat(t *testing.T) {

	fs := example1.FileSystem()
	assetvfs := vfs.OS(ASSETS_PATH)

	tests := []struct {
		desc string
		path string
		err  string
	}{
		{
			desc: "test open file",
			path: "hello.txt",
		},
		{
			desc: "test open sub-directory file",
			path: "hello/world.txt",
		},
		{
			desc: "test open directory",
			path: "hello",
		},
		{
			desc: "test open non-exists path",
			path: "notfound",
			err:  msgNotFound("Stat", "notfound"),
		},
	}

	for i, test := range tests {
		t.Logf("test fs.Stat  %d: %s", i+1, test.desc)

		// get the file/dir in the bindatafs
		targetStat, err := fs.Stat(test.path)
		if test.err != "" {
			if err == nil {
				t.Errorf("expected error %#v, got nil", test.err)
			} else if want, have := test.err, err.Error(); want != have {
				t.Errorf("expected error %#v, got %#v", want, have)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if targetStat == nil {
			t.Errorf("targetStat is nil")
			continue
		}

		// get the counter part in the source assets
		srcFile, err := os.Open(ASSETS_PATH + test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}
		defer srcFile.Close()
		srcStat, err := srcFile.Stat()

		if err = fileInfoEqual(srcStat, targetStat); err != nil {
			t.Errorf("error: %s", err.Error())
		}

		// get the counter part in vfs.OS file system
		vfsStat, err := assetvfs.Stat(test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}
		if err = fileInfoEqual(vfsStat, targetStat); err != nil {
			t.Errorf("error: %s", err.Error())
		}

	}

	for i, test := range tests {
		t.Logf("test fs.Lstat %d: %s", i+1, test.desc)

		// get the file/dir in the bindatafs
		targetStat, err := fs.Lstat(test.path)
		if test.err != "" {
			if err == nil {
				t.Errorf("expected error %#v, got nil", test.err)
			} else if want, have := test.err, err.Error(); want != have {
				t.Errorf("expected error %#v, got %#v", want, have)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if targetStat == nil {
			t.Errorf("targetStat is nil")
			continue
		}

		// get the counter part in the source assets
		srcFile, err := os.Open(ASSETS_PATH + test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}
		defer srcFile.Close()
		srcStat, err := srcFile.Stat()

		if err = fileInfoEqual(srcStat, targetStat); err != nil {
			t.Errorf("error: %s", err.Error())
		}

		// get the counter part in vfs.OS file system
		vfsStat, err := assetvfs.Stat(test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}
		if err = fileInfoEqual(vfsStat, targetStat); err != nil {
			t.Errorf("error: %s", err.Error())
		}
	}

}

func TestFileSystem_Readdir(t *testing.T) {

	fs := example1.FileSystem()
	assetvfs := vfs.OS(ASSETS_PATH)

	tests := []struct {
		desc  string
		path  string
		err   string
		files map[string]string
	}{
		{
			desc: "test open file",
			path: "hello.txt",
			err:  msgNotFound("ReadDir", "hello.txt"),
		},
		{
			desc: "test open sub-directory file",
			path: "hello/world.txt",
			err:  msgNotFound("ReadDir", "hello/world.txt"),
		},
		{
			desc: "test open directory",
			path: "hello",
			files: map[string]string{
				"bar.txt":   "file",
				"world.txt": "file",
			},
		},
		{
			desc: "test open root directory",
			path: "",
			files: map[string]string{
				"hello":      "dir",
				"hello.txt":  "file",
				"index.html": "file",
			},
		},
		{
			desc: "test open non-exists path",
			path: "notfound",
			err:  msgNotFound("ReadDir", "notfound"),
		},
	}

	for i, test := range tests {
		t.Logf("test fs.ReadDir %d: %s", i+1, test.desc)

		fsList, err := fs.ReadDir(test.path)
		if test.err != "" {
			if err == nil {
				t.Errorf("expected error %#v, got nil", test.err)
			} else if want, have := test.err, err.Error(); want != have {
				t.Errorf("expected %#v, got %#v", want, have)
			}
			continue
		}

		if want, have := len(test.files), len(fsList); want != have {
			t.Errorf("expected len(fsList) to be %d, got %d", want, have)
		}

		vfsList, err := assetvfs.ReadDir(test.path)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		vfsMap := make(map[string]os.FileInfo)
		for _, fi := range vfsList {
			vfsMap[fi.Name()] = fi
		}

		for _, fi := range fsList {
			if _, ok := test.files[fi.Name()]; !ok {
				t.Errorf("unexpected entity: %s", fi.Name())
			}
			if _, ok := vfsMap[fi.Name()]; !ok {
				t.Errorf("unexpected entity: %s", fi.Name())
			}
		}

	}

}

func TestFileSystem_RootType(t *testing.T) {
	fs := bindatafs.New("hello", nil, nil, nil)
	if want, have := vfs.RootType(""), fs.RootType("hello"); want != have {
		t.Logf("expected %#v, got %#v", want, have)
	}
}

func TestFileSystem_String(t *testing.T) {
	fs := bindatafs.New("hello", nil, nil, nil)
	if want, have := "hello", fs.String(); want != have {
		t.Logf("expected %#v, got %#v", want, have)
	}
}
