package main

import (
	"io"
	"os"

	"github.com/JustDaile/go-denest-structs/pkg/core"
)

func main() {
	var bs []byte
	fi, _ := os.Stdin.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		bs, _ = io.ReadAll(os.Stdin)
	} else {
		os.Stdout.Write([]byte("must pipe data\n"))
		os.Exit(1)
	}

	denester := core.NewStructDenester(bs)
	if err := denester.Process(os.Stdout); err != nil {
		panic(err)
	}
}
