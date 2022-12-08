//go:build test
// +build test

package std

func TestSetContext(c contextI) {
	ctx = c
}
