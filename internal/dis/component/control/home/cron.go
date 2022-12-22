package home

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

type UpdateInstanceList struct {
	component.Base
	Dependencies struct {
		Home *Home
	}
}

var (
	_ component.Cronable = (*UpdateInstanceList)(nil)
)

func (*UpdateInstanceList) TaskName() string {
	return "instance list"
}

func (ul *UpdateInstanceList) Cron(ctx context.Context) error {
	names, err := ul.Dependencies.Home.instanceMap(ctx)
	if err != nil {
		return err
	}

	ul.Dependencies.Home.instanceNames.Set(names)
	return nil
}

type UpdateRedirect struct {
	component.Base
	Dependencies struct {
		Home *Home
	}
}

var (
	_ component.Cronable = (*UpdateRedirect)(nil)
)

func (ur *UpdateRedirect) TaskName() string {
	return "redirect"
}

func (ur *UpdateRedirect) Cron(ctx context.Context) error {
	redirect, err := ur.Dependencies.Home.loadRedirect(ctx)
	if err != nil {
		return err
	}

	ur.Dependencies.Home.redirect.Set(&redirect)
	return nil
}

type UpdateHome struct {
	component.Base
	Dependencies struct {
		Home *Home
	}
}

var (
	_ component.Cronable = (*UpdateHome)(nil)
)

func (ur *UpdateHome) TaskName() string {
	return "home render"
}

func (ur *UpdateHome) Cron(ctx context.Context) error {
	bytes, err := ur.Dependencies.Home.homeRender(ctx)
	if err != nil {
		return err
	}

	ur.Dependencies.Home.homeBytes.Set(bytes)
	return nil
}
