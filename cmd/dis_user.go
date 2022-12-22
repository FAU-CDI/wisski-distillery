package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// DisUser is the 'dis_user' command
var DisUser wisski_distillery.Command = disUser{}

type disUser struct {
	CreateUser bool `short:"c" long:"create" description:"create a new user"`
	DeleteUser bool `short:"d" long:"delete" description:"delete a user"`

	InfoUser  bool `short:"i" long:"info" description:"show information about a user"`
	ListUsers bool `short:"l" long:"list" description:"list all users"`

	SetPassword   bool `short:"s" long:"set-password" description:"interactively set a user password"`
	UnsetPassword bool `short:"u" long:"unset-password" description:"delete a users password and block the account"`
	CheckPassword bool `short:"p" long:"check-password" description:"interactively check a user password"`

	Positionals struct {
		User string `positional-arg-name:"USER" description:"username to manage. May be omitted for some actions"`
	} `positional-args:"true"`
}

func (disUser) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "dis_user",
		Description: "manage distillery users",
	}
}

var errUserRequired = exit.Error{
	Message:  "`USER` argument is required",
	ExitCode: exit.ExitCommandArguments,
}

func (du disUser) AfterParse() error {
	var counter int
	for _, action := range []bool{
		du.CreateUser,
		du.InfoUser,
		du.DeleteUser,
		du.SetPassword,
		du.UnsetPassword,
		du.CheckPassword,
		du.ListUsers,
	} {
		if action {
			counter++
		}
	}

	if counter != 1 {
		return errNoActionSelected
	}

	if !du.ListUsers && du.Positionals.User == "" {
		return errUserRequired
	}

	return nil
}

func (du disUser) Run(context wisski_distillery.Context) error {
	switch {
	case du.InfoUser:
		return du.runInfo(context)
	case du.CreateUser:
		return du.runCreate(context)
	case du.DeleteUser:
		return du.runDelete(context)
	case du.SetPassword:
		return du.runSetPassword(context)
	case du.UnsetPassword:
		return du.runUnsetPassword(context)
	case du.CheckPassword:
		return du.runCheckPassword(context)
	case du.ListUsers:
		return du.runListUsers(context)
	}
	panic("never reached")
}

func (du disUser) runInfo(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().User(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}

	context.Println(user)
	return nil
}

func (du disUser) runCreate(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().CreateUser(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}

	context.Println(user)
	return nil
}

func (du disUser) runDelete(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().User(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}

	return user.Delete(context.Context)
}

func (du disUser) runSetPassword(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().User(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}

	var passwd string
	{
		context.Printf("Enter new password for user %s:", du.Positionals.User)
		passwd1, err := context.IOStream.ReadPassword()
		if err != nil {
			return err
		}
		context.Println()

		context.Printf("Enter the same password again:")
		passwd, err = context.IOStream.ReadPassword()
		if err != nil {
			return err
		}
		context.Println()

		if passwd != passwd1 {
			return errPasswordsNotIdentical
		}
		if len(passwd) == 0 {
			return errPasswordsNotIdentical
		}
	}

	return user.SetPassword(context.Context, []byte(passwd))
}

func (du disUser) runUnsetPassword(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().User(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}

	return user.UnsetPassword(context.Context)
}

func (du disUser) runCheckPassword(context wisski_distillery.Context) error {
	user, err := context.Environment.Auth().User(context.Context, du.Positionals.User)
	if err != nil {
		return err
	}

	context.Printf("Enter password for %s:", du.Positionals.User)

	candidate, err := context.IOStream.ReadPassword()
	if err != nil {
		return err
	}
	context.Println()

	return user.CheckPassword(context.Context, []byte(candidate))
}

func (du disUser) runListUsers(context wisski_distillery.Context) error {
	users, err := context.Environment.Auth().Users(context.Context)
	if err != nil {
		return err
	}
	for _, user := range users {
		context.Println(user)
	}
	return nil
}