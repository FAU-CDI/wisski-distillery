//spellchecker:words passwordx
package passwordx

//spellchecker:words embed github pkglib password
import (
	"embed"

	"github.com/tkw1536/pkglib/password"
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
