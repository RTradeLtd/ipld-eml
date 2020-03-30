package main

import (
	"log"
	"os"

	"github.com/RTradeLtd/ipld-eml/analysis"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "eml-util"
	app.Usage = "cli util for managing ipld-eml objects, and benchmarking"
	app.Authors = append(app.Authors, &cli.Author{
		Name:  "Alex Trottier",
		Email: "postables@rtradetechnologies.com",
	})
	app.Commands = cli.Commands{
		{
			Name:    "generate-fake-emails",
			Aliases: []string{"gen-fake-emails", "gfe"},
			Usage:   "generates fake emails to use for benchmarking",
			Action: func(c *cli.Context) error {
				if err := os.Mkdir(c.String("outdir"), os.ModePerm); !os.IsExist(err) {
					return err
				}
				parts, err := analysis.GenerateMessages(c.Int("count"))
				if err != nil {
					return err
				}
				return analysis.WritePartsToDisk(parts, c.String("outdir"))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "outdir",
					Usage: "directory to store emails is",
					Value: "outdir",
				},
				&cli.IntFlag{
					Name:  "count",
					Usage: "number of emails to generate",
					Value: 100,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
