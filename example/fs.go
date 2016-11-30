package example

import "github.com/go-serve/bindatafs"

// FileSystem returns a Filesystem implementation for the given assets
func FileSystem() bindatafs.FileSystem {
	return bindatafs.New("assets://", Asset, AssetDir, AssetInfo)
}
