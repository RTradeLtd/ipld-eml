package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/RTradeLtd/go-temporalx-sdk/client"
	ipldeml "github.com/RTradeLtd/ipld-eml"
	"github.com/RTradeLtd/ipld-eml/analysis"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := cli.NewApp()
	app.Name = "eml-util"
	app.Usage = "cli util for managing ipld-eml objects, and benchmarking"
	app.Authors = append(app.Authors, &cli.Author{
		Name:  "Alex Trottier",
		Email: "postables@rtradetechnologies.com",
	})
	app.Commands = cli.Commands{
		{
			Name:    "convert",
			Aliases: []string{"conv", "c"},
			Usage:   "read emails from directory uploading to ipfs",
			Action: func(c *cli.Context) error {
				cl, err := client.NewClient(client.Opts{
					ListenAddress: c.String("endpoint"),
					Insecure:      c.Bool("insecure"),
				})
				if err != nil {
					return err
				}
				converter := ipldeml.NewConverter(ctx, cl)
				res, err := converter.AddFromDirectory(c.String("email.dir"))
				if err != nil {
					return err
				}
				formatted := ""
				for name, hash := range res {
					formatted = fmt.Sprintf("%sfile: %s\thash: %s\n", formatted, name, hash)
					if err != nil {
						return err
					}
				}
				return ioutil.WriteFile(c.String("save.file"), []byte(formatted), os.FileMode(0642))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "save.file",
					Usage: "file to save name -> hash informatino",
					Value: "converted_results.txt",
				},
				&cli.StringFlag{
					Name:  "endpoint",
					Usage: "temporalx endpoint to connect to",
					Value: "localhost:9090",
				},
				&cli.StringFlag{
					Name:  "email.dir",
					Usage: "directory containing emails, must not be a recursive directory",
					Value: "outdir",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "establish an insecure connection to temporalx",
					Value: true,
				},
			},
		},
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
