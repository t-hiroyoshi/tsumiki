package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/t-hiroyoshi/tsumiki/container"
	"github.com/t-hiroyoshi/tsumiki/version"
	"github.com/urfave/cli"
)

func InstallAction(c *cli.Context) error {
	client, err := container.NewContainerClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	for _, p := range c.Args() {
		fmt.Printf("Installing %s...\n", p)

		info, err := client.ImagePull(ctx, p)
		if err != nil {
			return err
		}

		err = version.AddPackage(*info)
		if err != nil {
			return err
		}

		fmt.Printf("Installed %s\n", p)
	}

	return nil
}

func UninstallAction(c *cli.Context) error {
	client, err := container.NewContainerClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	for _, p := range c.Args() {
		pkg, err := version.GetPackage(p)
		if err != nil {
			return err
		}
		if pkg == nil {
			return errors.New(fmt.Sprintf("Not installed package %v", p))
		}

		err = client.ImageRemove(ctx, *pkg)
		if err != nil {
			return err
		}

		err = version.RemovePackage(*pkg, client.Host())
		if err != nil {
			return err
		}
	}

	return nil
}

func ListActions(c *cli.Context) error {
	pkgs, err := version.GetPackages()
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		fmt.Printf("%s:%s ", pkg.Tag)
	}

	return nil
}
