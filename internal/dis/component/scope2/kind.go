package scope2

// Kind represents a kind of scope.
// The existing kinds of scopes are defined below.
type Kind string

const (
	KindNever  Kind = ""       // never fullfilled
	KindAlways Kind = "always" // fullfilled by anyone

	KindValidUser   Kind = "user"       // fullfilled by a valid distillery
	KindValidAdmin  Kind = "admin"      // fullfukked by a valid distillery admin
	KindSecureAdmin Kind = "admin-totp" // fullfilled by a distillery admin with totp setup

	KindInstanceUser        Kind = "instance-user"       // fullfilled by a user for a specific instance
	KindInstanceAdmin       Kind = "instance-admin"      // fullfilled by an admin for a specific instance
	KindInstanceSecureAdmin Kind = "instance-admin-totp" // fullfilled by an admin for the specific instance
)

// kindInfo holds information about specific kinds
type kindInfo struct {
	NeedsInstance bool
}

// valid kinds that are known
var kindInfos = map[Kind]kindInfo{
	KindNever:  {NeedsInstance: false},
	KindAlways: {NeedsInstance: false},

	KindValidUser:   {NeedsInstance: false},
	KindValidAdmin:  {NeedsInstance: false},
	KindSecureAdmin: {NeedsInstance: false},

	KindInstanceUser:        {NeedsInstance: true},
	KindInstanceAdmin:       {NeedsInstance: true},
	KindInstanceSecureAdmin: {NeedsInstance: true},
}

// Valid checks if a kind is valid
func (kind Kind) Valid() bool {
	_, ok := kindInfos[kind]
	return ok
}

// RequiresInstance checks if a kind requires an instance
func (kind Kind) RequiresInstance() bool {
	return kindInfos[kind].NeedsInstance
}
