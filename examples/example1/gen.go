// This file only provides command to generate
// data locally

//go:generate go-bindata -ignore=assets.go -o assets.go -pkg=example1 -prefix=assets/ assets/...

package example1
