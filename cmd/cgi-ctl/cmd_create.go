package main

import (
	"fmt"
	"github.com/reddec/trusted-cgi/cmd/internal"
	internal_app "github.com/reddec/trusted-cgi/internal"
	"log"
	"os"
)

type create struct {
	remoteLink
	Args struct {
		Name string `name:"name" description:"project directory" required:"yes"`
	} `positional-args:"yes"`
}

func (cmd *create) Execute(args []string) error {
	ctx, closer := internal.SignalContext()
	defer closer()
	log.Println("login...")
	token, err := cmd.Token(ctx)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	err = os.MkdirAll(cmd.Args.Name, 0755)
	if err != nil {
		return fmt.Errorf("prepare directory: %w", err)
	}

	err = os.Chdir(cmd.Args.Name)
	if err != nil {
		return fmt.Errorf("change dir: %w", err)
	}

	log.Println("creating...")
	info, err := cmd.Project().Create(ctx, token)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	log.Println("created", info.UID)
	log.Println("saving info....")

	var cf controlFile
	cf.URL = cmd.URL
	err = cf.Save(controlFilename)
	if err != nil {
		return fmt.Errorf("save control file: %w", err)
	}
	err = appendIfNoLineFile(internal_app.CGIIgnore, controlFilename)
	if err != nil {
		return fmt.Errorf("update cgiignore file: %w", err)
	}

	err = info.Manifest.SaveAs(internal_app.ManifestFile)
	if err != nil {
		return fmt.Errorf("save manifest: %w", err)
	}
	log.Println("done")
	return nil
}
