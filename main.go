package main

import (
	"fmt"

	"github.com/tupyy/vwap/internal/conf"
)

// CommitID contains the SHA1 Git commit of the build.
// It's evaluated during compilation.
var CommitID string

func main() {
	config := conf.Get()
	fmt.Println(config)
}
