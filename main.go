/*
Copyright 2020 somen440

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	version = "v1.0.2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:      "reflect",
				Usage:     "reflect schema from database.",
				ArgsUsage: "<database>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "driver",
						Value: "mysql",
					},
					&cli.StringFlag{
						Name:    "user",
						Aliases: []string{"u"},
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
					},
					&cli.StringFlag{
						Name:    "host",
						Value:   "localhost",
						Aliases: []string{"H"},
					},
					&cli.StringFlag{
						Name:    "port",
						Value:   "3306",
						Aliases: []string{"P"},
					},
					&cli.StringFlag{
						Name:    "output",
						Value:   "schema",
						Aliases: []string{"o"},
					},
				},
				Action: reflectAction,
			},
			{
				Name:  "generate",
				Usage: "generate migration template.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "output",
						Value: "migrations",
					},
				},
				Action: generateAction,
			},
			{
				Name:      "diff",
				Usage:     "diff schema <-> database.",
				ArgsUsage: "<database> <schema>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "driver",
						Value: "mysql",
					},
					&cli.StringFlag{
						Name:    "user",
						Aliases: []string{"u"},
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
					},
					&cli.StringFlag{
						Name:    "host",
						Value:   "localhost",
						Aliases: []string{"H"},
					},
					&cli.StringFlag{
						Name:    "port",
						Value:   "3306",
						Aliases: []string{"P"},
					},
					&cli.StringFlag{
						Name:    "ignores",
						Aliases: []string{"I"},
					},
				},
				Action: diffAction,
			},
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					fmt.Println(version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func reflectAction(c *cli.Context) error {
	database := c.Args().Get(0)

	conf := &DBConfig{
		Driver:   DbDriver(c.String("driver")),
		User:     c.String("user"),
		Password: c.String("password"),
		Host:     c.String("host"),
		Port:     c.String("port"),
		Database: database,
	}
	reflect := NewReflect(conf, c.String("output"))
	return reflect.Exec()
}

func generateAction(c *cli.Context) error {
	generate := NewGenerate(c.String("output"))
	return generate.Exec()
}

func diffAction(c *cli.Context) error {
	database := c.Args().Get(0)
	schema := c.Args().Get(1)

	conf := &DBConfig{
		Driver:   DbDriver(c.String("driver")),
		User:     c.String("user"),
		Password: c.String("password"),
		Host:     c.String("host"),
		Port:     c.String("port"),
		Database: database,
	}
	diff := NewDiff(conf)
	return diff.Exec(database, schema, strings.Split(c.String("ignores"), ","))
}
