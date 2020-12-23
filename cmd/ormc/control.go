package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dinic/ormc/pkg/config"
	"github.com/dinic/ormc/pkg/dbinfo"
	"github.com/dinic/ormc/pkg/model"
	"github.com/dinic/ormc/pkg/utils"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type Option struct {
	ConfigPath    string
	ConfigOutPath string
}

func (o *Option) SetDefaultConfigPath() error {
	if o.ConfigPath != "" {
		return nil
	}

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if utils.IsExist(currDir, ".ormc.yaml") {
		o.ConfigPath = filepath.Join(currDir, ".ormc.yaml")
		return nil
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if utils.IsExist(homedir, ".ormc.yaml") {
		o.ConfigPath = filepath.Join(homedir, ".ormc.yaml")
		return nil
	}

	return errors.New("not found ormc config")
}

var (
	opt Option

	cmdCode = cli.Command{
		Name:  "code",
		Usage: "gen go code from mysql",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "config, c",
				Usage:       "set code config path",
				Value:       "",
				Destination: &opt.ConfigPath,
			},
		},
		Action: func(c *cli.Context) error {
			if err := opt.SetDefaultConfigPath(); err != nil {
				log.Printf("code set default config failed, err: %s", err)
				return err
			}

			cf := config.Parse(opt.ConfigPath)
			database := make([]*dbinfo.DB, 0, len(cf.Database)+1)
			for _, mc := range cf.Database {
				db := dbinfo.NewDB(mc)
				db.Load()
				database = append(database, db)
			}

			renders := make([]*model.Renderer, 0, 12)
			for _, db := range database {
				mdb := model.NewDbModel(db)
				rs := mdb.Render()
				if len(rs) > 0 {
					renders = append(renders, rs...)
				}
			}

			for _, r := range renders {
				r.Render()
			}
			return nil
		},
	}

	cmdConfig = cli.Command{
		Name:  "config",
		Usage: "gen config, default save as ./.ormc.yaml",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "out, o",
				Usage:       "set gen config output path",
				Value:       "./.ormc.yaml",
				Destination: &opt.ConfigOutPath,
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)
