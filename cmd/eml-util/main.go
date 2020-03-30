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
				os.Mkdir(c.String("outdir"), os.ModePerm)
				return analysis.GenerateMessages(
					c.String("outdir"),
					c.Int("email.count"),
					c.Int("emoji.count"),
					c.Int("paragraph.count"),
				)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "outdir",
					Usage: "directory to store emails is",
					Value: "outdir",
				},
				&cli.IntFlag{
					Name:  "email.count",
					Usage: "number of emails to generate",
					Value: 100,
				},
				&cli.IntFlag{
					Name:  "emoji.count",
					Usage: "number of emojis to add in email",
					Value: 100,
				},
				&cli.IntFlag{
					Name:  "paragraph.count",
					Usage: "number of paragraphs to generate",
					Value: 10,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
