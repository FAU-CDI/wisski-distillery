package home

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

type UpdateInstanceList struct {
	component.Base
	Home *Home
}

var (
	_ component.Cronable = (*UpdateInstanceList)(nil)
)

func (*UpdateInstanceList) TaskName() string {
	return "instance list"
}

func (ul *UpdateInstanceList) Cron(ctx context.Context) error {
	names, err := ul.Home.instanceMap(ctx)
	if err != nil {
		return err
	}

	ul.Home.instanceNames.Set(names)
	return nil
}

type UpdateRedirect struct {
	component.Base
	Home *Home
}

var (
	_ component.Cronable = (*UpdateRedirect)(nil)
)

func (ur *UpdateRedirect) TaskName() string {
	return "redirect"
}

func (ur *UpdateRedirect) Cron(ctx context.Context) error {
	redirect, err := ur.Home.loadRedirect(ctx)
	if err != nil {
		return err
	}

	ur.Home.redirect.Set(&redirect)
	return nil
}

type UpdateHome struct {
	component.Base
	Home *Home
}

var (
	_ component.Cronable = (*UpdateHome)(nil)
)

func (ur *UpdateHome) TaskName() string {
	return "home render"
}

func (ur *UpdateHome) Cron(ctx context.Context) error {
	bytes, err := ur.Home.homeRender(ctx)
	if err != nil {
		return err
	}

	ur.Home.homeBytes.Set(bytes)
	return nil
}
