package hint64

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/kong"
)

type KongFlag int64

var KongTypeMapper = kong.TypeMapper(
	reflect.TypeOf(KongFlag(0)),
	kong.MapperFunc(func(ctx *kong.DecodeContext, target reflect.Value) error {
		token, err := ctx.Scan.PopValue("hint64")
		if err != nil {
			return err
		}

		str, ok := token.Value.(string)
		if !ok {
			return fmt.Errorf("expected an int but got %q (%T)", token, token.Value)
		}

		num, err := Parse(str)
		if err != nil {
			// TODO: wrap original errors after modifying error messages
			return fmt.Errorf("expected a valid human-readable 64-bit integer but got %q", str)
		}

		target.Set(reflect.ValueOf(KongFlag(num)))
		return nil
	}),
)
