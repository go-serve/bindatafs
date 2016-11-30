// This file only provides command to generate
// data locally

//go:generate go-bindata -o assets.go -ignore=gen.go -ignore=assets.go -pkg=example -prefix=assets/ assets/...

package example
