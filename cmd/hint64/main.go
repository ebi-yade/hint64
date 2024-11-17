package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/ebi-yade/hint64"
)

var cli = struct {
	Number hint64.KongFlag `arg:"" name:"number" help:"A human-readable expression of a number."`

	// TODO: version flag
}{}

func main() {
	kong.Parse(&cli,
		hint64.KongTypeMapper, // Important!
	)
	fmt.Printf("%d\n", cli.Number)
}
