//spellchecker:words passwordx
package passwordx

//spellchecker:words embed pkglib password
import (
	"embed"

	"go.tkw01536.de/pkglib/password"
)

//go:embed common
var commonEmbed embed.FS

var Sources []password.PasswordSource

func init() {
	var err error
	Sources, err = password.NewSources(commonEmbed, "**/*.txt")
	if err != nil {
		panic(err)
	}
	if len(Sources) == 0 {
		panic("no sources")
	}
}
