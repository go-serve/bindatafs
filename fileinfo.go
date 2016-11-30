package bindatafs

import (
	"os"
	"time"
)

// Type of an asset provided by FileSystem
type Type int

// fileInfo implements FileInfo
type fileInfo struct {
	name string
	os.FileInfo
}

// Name implements os.FileInfo
func (fi *fileInfo) Name() string {
	return fi.name
}

// Size gives length in bytes for regular files;
// system-dependent for others
func (fi *fileInfo) Size() int64 {
	return fi.FileInfo.Size()
}

// Mode gives file mode bits
func (fi *fileInfo) Mode() os.FileMode {
	return fi.FileInfo.Mode()&os.ModeType | 0444
}

// ModTime gives modification time
func (fi *fileInfo) ModTime() (t time.Time) {
	return fi.FileInfo.ModTime()
}

// IsDir is abbreviation for Mode().IsDir()
func (fi *fileInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

// Sys gives underlying data source (can return nil)
func (fi *fileInfo) Sys() interface{} {
	return nil
}

// dirInfo implements FileInfo for directory in the assets
type dirInfo struct {
	name string
}

// Name gives base name of the file
func (fi *dirInfo) Name() string {
	return fi.name
}

// Size gives length in bytes for regular files;
// system-dependent for others
func (fi *dirInfo) Size() int64 {
	return 0 // hard code 0 for now (originally system-dependent)
}

// Mode gives file mode bits
func (fi *dirInfo) Mode() os.FileMode {
	return os.ModeDir | 0777
}

// ModTime gives modification time
func (fi *dirInfo) ModTime() (t time.Time) {
	return time.Unix(0, 0)
}

// IsDir is abbreviation for Mode().IsDir()
func (fi *dirInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

// Sys gives underlying data source (can return nil)
func (fi *dirInfo) Sys() interface{} {
	return nil
}
