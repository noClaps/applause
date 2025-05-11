package applause_test

import (
	"testing"

	"github.com/noclaps/applause"
)

func TestParse(t *testing.T) {
	type Args struct {
		name         string `help:"The name of the package"`
		skipPeer     bool   `type:"option" name:"skip-peer" help:"Skip counting peer dependencies" short:"p"`
		skipOptional bool   `type:"option" name:"skip-optional" help:"Skip counting optional dependencies" short:"o"`
		version      string `type:"option" help:"The version of the package being fetched"`
	}
	_ = applause.Parse(Args{})
}
