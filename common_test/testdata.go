package common_test

import "embed"

// Allow contents of the testdata subdirectory to be embedded in unit test
// code.
//
//go:embed testdata
var TestData embed.FS
