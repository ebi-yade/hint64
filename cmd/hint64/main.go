package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/ebi-yade/hint64"
)

var cli = struct {
	NumString string `arg:"" name:"num-string" help:"A human-readable expression of a number."`
	// TODO: version flag
}{}

func main() {
	if err := main_(); err != nil {
		slog.Error(fmt.Sprint("Error: +v", err))
		os.Exit(1)
	}
}

func main_() error {
	kong.Parse(&cli)

	num, err := hint64.Parse(cli.NumString)
	if err != nil {
		return fmt.Errorf("error hnum.ParseInt64: %w", err)
	}
	fmt.Printf("%d\n", num)

	return nil
}
