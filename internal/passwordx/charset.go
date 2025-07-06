//spellchecker:words passwordx
package passwordx

//spellchecker:words github pkglib password
import "go.tkw01536.de/pkglib/password"

// Safe is a charset used for generating passwords that can be safely passed without having to be escaped.
const Safe = password.DefaultCharSet

// Printable is a charset that contains all printable ascii characters.
const Printable = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

// Identifier is a charset to be used to generate unique identifiers.
// These are typically used for snapshots and names.
const Identifier = password.DefaultCharSet
