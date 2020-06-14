package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
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
				},
				Action: diffAction,
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
	return diff.Exec(database, schema)
}
