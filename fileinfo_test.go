package bindatafs

import (
	"os"
	"testing"
)

func Test_fileInfo(t *testing.T) {
	var i os.FileInfo = &fileInfo{}
	_ = i
	t.Log("*bindatafs.FileInfo{} implements os.FileInfo interface")
}

func Test_dirInfo(t *testing.T) {
	var i os.FileInfo = &dirInfo{}
	_ = i
	t.Log("*bindatafs.DirInfo{} implements os.FileInfo interface")
}
