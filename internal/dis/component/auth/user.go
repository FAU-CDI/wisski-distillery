//spellchecker:words auth
package auth

//spellchecker:words bytes context encoding base errors image reflect strings github wisski distillery internal component models passwordx pquerna totp pkglib password golang crypto bcrypt
import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.tkw01536.de/pkglib/password"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

func (auth *Auth) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: models.User{},
	}
}

// Users returns all users in the database.
func (auth *Auth) Users(ctx context.Context) (users []*AuthUser, err error) {
	// query the user table
	table, err := sql.OpenInterface[models.User](ctx, auth.dependencies.SQL, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}

	// find all the users
	dUsers, err := table.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	// and map them to high-level user objects
	users = make([]*AuthUser, len(dUsers))
	for i, user := range dUsers {
		users[i] = &AuthUser{
			User: user,
			auth: auth,
		}
	}

	return users, nil
}

// User returns a single user.
// If the user does not exist, returns [ErrUserNotFound].
func (auth *Auth) User(ctx context.Context, name string) (au *AuthUser, err error) {
	// quick and dirty check for the empty username (which is not allowed)
	if name == "" {
		return nil, ErrUserNotFound
	}

	table, err := sql.OpenInterface[models.User](ctx, auth.dependencies.SQL, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}

	au = &AuthUser{}
	au.User, err = table.Where(&models.User{User: name}).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	au.auth = auth

	return
}

// CreateUser creates a new user and returns it.
// The user is not associated to any WissKIs, and has no password set.
func (auth *Auth) CreateUser(ctx context.Context, name string) (user *AuthUser, err error) {
	user = &AuthUser{
		User: models.User{
			User: name,
		},
		auth: auth,
	}
	user.SetAdmin(false)
	user.SetEnabled(false)
	user.SetTOTPEnabled(false)

	table, err := sql.OpenInterface[models.User](ctx, auth.dependencies.SQL, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}

	if err := table.Create(ctx, &user.User); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// AuthUser represents an authorized user.
type AuthUser struct {
	auth *Auth
	models.User
}

func (au *AuthUser) String() string {
	if au == nil {
		return "User{nil}"
	}
	hasPassword := len(au.PasswordHash) > 0
	return fmt.Sprintf("User{Name:%q,Enabled:%t,HasPassword:%t,Admin:%t}", au.User.User, au.IsEnabled(), hasPassword, au.IsAdmin())
}

var (
	ErrTOTPEnabled  = errors.New("TOTP is enabled")
	ErrTOTPDisabled = errors.New("TOTP is disabled")
	ErrTOTPFailed   = errors.New("TOTP failed")
)

func (au *AuthUser) TOTP() (*otp.Key, error) {
	// TODO: make this private
	if au.TOTPURL == "" {
		return nil, ErrTOTPDisabled
	}

	key, err := otp.NewKeyFromURL(au.TOTPURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get key from url: %w", err)
	}
	return key, nil
}

// CheckTOTP validates the given totp passcode against the saved secret.
// If totp is not enabled, any passcode will pass the check.
func (au *AuthUser) CheckTOTP(passcode string) error {
	secret, err := au.TOTP()
	if err != nil {
		return err
	}

	if au.IsTOTPEnabled() && !totp.Validate(passcode, secret.Secret()) {
		return ErrTOTPFailed
	}
	return nil
}

// NewTOTP generates a new TOTP secret, returning a totp key.
func (au *AuthUser) NewTOTP(ctx context.Context) (*otp.Key, error) {
	if au.IsTOTPEnabled() {
		return nil, ErrTOTPEnabled
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "WissKI Distillery",
		AccountName: au.User.User,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	au.TOTPURL = key.URL()
	return key, au.Save(ctx)
}

func TOTPLink(secret *otp.Key, width, height int) (string, error) {
	// make an image
	img, err := secret.Image(width, height)
	if err != nil {
		return "", fmt.Errorf("failed to create QR code: %w", err)
	}

	// encode image as base64
	var buffer bytes.Buffer

	if err := png.Encode(&buffer, img); err != nil {
		return "", fmt.Errorf("failed to encode QR code to png: %w", err)
	}

	// return the image url
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

// EnableTOTP enables totp for the given user.
func (au *AuthUser) EnableTOTP(ctx context.Context, passcode string) error {
	secret, err := au.TOTP()
	if err != nil {
		return err
	}
	if !totp.Validate(passcode, secret.Secret()) {
		return ErrTOTPFailed
	}

	au.SetTOTPEnabled(true)
	return au.Save(ctx)
}

// DisableTOTP disables totp for the given user.
func (au *AuthUser) DisableTOTP(ctx context.Context) (err error) {
	au.SetTOTPEnabled(false)
	au.TOTPURL = ""
	return au.Save(ctx)
}

// SetPassword sets the password for this user and turns the user on.
func (au *AuthUser) SetPassword(ctx context.Context, password []byte) (err error) {
	au.PasswordHash, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate hash from password: %w", err)
	}
	au.SetEnabled(true)
	return au.Save(ctx)
}

// UnsetPassword removes the password from this user, and disables them.
func (au *AuthUser) UnsetPassword(ctx context.Context) error {
	au.PasswordHash = nil
	au.SetEnabled(false)
	return au.Save(ctx)
}

const MinPasswordLength = 8

var (
	ErrPolicyBlank    = errors.New("password is blank")
	ErrPolicyTooShort = fmt.Errorf("password is too short: minimum length %d", MinPasswordLength)
	ErrPolicyKnown    = errors.New("password is on the list of known passwords")
	ErrPolicyUsername = errors.New("password may not be identical to username")
)

// CheckPasswordPolicy checks if the given password would pass the password policy.
//
// The password policy checks that the password has a minimum length of [MinPasswordLength]
// and that it is not a common password.
// It also checks that password and username are not identical.
func (auth *Auth) CheckPasswordPolicy(candidate string, username string) error {
	if candidate == "" {
		return ErrPolicyBlank
	}

	if strings.EqualFold(candidate, username) {
		return ErrPolicyUsername
	}

	if len(candidate) < MinPasswordLength {
		return ErrPolicyTooShort
	}

	if err := password.CheckCommonPassword(func(common string) (bool, error) { return common == candidate, nil }, passwordx.Sources...); err != nil {
		return ErrPolicyKnown
	}

	return nil
}

func (au *AuthUser) CheckPasswordPolicy(candidate string) error {
	return au.auth.CheckPasswordPolicy(candidate, au.User.User)
}

var (
	ErrNoUser       = errors.New("user is nil")
	ErrUserDisabled = errors.New("user is disabled")
	ErrUserBlank    = errors.New("user has no password set")
)

// CheckPassword checks if this user can login with the provided password.
// Returns nil on success, an error otherwise.
func (au *AuthUser) CheckPassword(ctx context.Context, password []byte) error {
	if au == nil {
		return ErrNoUser
	}
	if !au.IsEnabled() {
		return ErrUserDisabled
	}

	if len(au.PasswordHash) == 0 {
		return ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword(au.PasswordHash, password); err != nil {
		return fmt.Errorf("wrong password: %w", err)
	}
	return nil
}

func (au *AuthUser) CheckCredentials(ctx context.Context, password []byte, passcode string) error {
	if err := au.CheckPassword(ctx, password); err != nil {
		return err
	}
	if err := au.CheckTOTP(passcode); err != nil && !errors.Is(err, ErrTOTPDisabled) {
		return err
	}
	return nil
}

// MakeAdmin makes this user an admin, and saves the update in the database.
// If the user is already an admin, does not return an error.
func (au *AuthUser) MakeAdmin(ctx context.Context) error {
	au.SetAdmin(true)
	return au.Save(ctx)
}

// MakeRegular removes admin rights from this user.
// If this user is not an dmin, does not return an error.
func (au *AuthUser) MakeRegular(ctx context.Context) error {
	au.SetAdmin(false)
	return au.Save(ctx)
}

// Save saves the given user in the database.
func (au *AuthUser) Save(ctx context.Context) error {
	table, err := sql.OpenInterface[models.User](ctx, au.auth.dependencies.SQL, au.auth)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	{
		_, err := table.Select("*").Updates(ctx, au.User)
		if err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}
	}
	return nil
}

// Delete deletes the user from the database.
func (au *AuthUser) Delete(ctx context.Context) error {
	// run all the user delete hooks
	for _, c := range au.auth.dependencies.UserDeleteHooks {
		if err := c.OnUserDelete(ctx, &au.User); err != nil {
			return fmt.Errorf("failed to run delete hook %q: %w", c.Name(), err)
		}
	}

	table, err := sql.OpenInterface[models.User](ctx, au.auth.dependencies.SQL, au.auth)
	if err != nil {
		return fmt.Errorf("failed to query table: %w", err)
	}

	if _, err := table.Where("pk = ?", au.User.Pk).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete user from db: %w", err)
	}
	return nil
}
