package args

import (
	"github.com/jessevdk/go-flags"
)

func IsHelp(err error) bool {
	if err == nil {
		return false
	}
	flagError, ok := err.(*flags.Error)
	if !ok {
		return false
	}
	if flagError.Type != flags.ErrHelp {
		return false
	}
	return true
}
