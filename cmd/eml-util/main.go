package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	app.Flags = []cli.Flag{
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
	}
	app.Commands = cli.Commands{
		{
			Name:    "benchmark",
			Aliases: []string{"bench", "b"},
			Usage:   "run specialized benchmark tool, calculating space savings",
			Description: `
this adds the "samples/generated" directory to ipfs,
calculating the total deduplicated size to compare to on disk storage.
input file is expected to a list of hashes **only**.
`,
			Action: func(c *cli.Context) error {
				cl, err := client.NewClient(client.Opts{
					ListenAddress: c.String("endpoint"),
					Insecure:      c.Bool("insecure"),
				})
				if err != nil {
					return err
				}
				contents, err := ioutil.ReadFile(c.String("input.file"))
				if err != nil {
					return err
				}
				hashes := strings.Split(string(contents), "\n")
				var parsed = make([]string, len(hashes))
				var max int
				for i, hash := range hashes {
					if hash == "" {
						continue
					}
					parsed[i] = hash
					max = i
				}
				converter := ipldeml.NewConverter(ctx, cl)
				size, err := converter.CalculateEmailSize(true, parsed[:max]...)
				if err != nil {
					return err
				}
				fmt.Println("total deduplicated size: ", size)
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "input.file",
					Usage: "file to get hash information from",
					Value: "converted_results.txt",
				},
			},
		},
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
				var numFiles int
				fh, err := ioutil.ReadDir(c.String("email.dir"))
				if err != nil {
					return err
				}
				for _, f := range fh {
					if !f.IsDir() {
						numFiles++
					}
				}
				converter := ipldeml.NewConverter(ctx, cl)
				res, err := converter.AddFromDirectory(c.String("email.dir"))
				if err != nil {
					return err
				}
				formatted := ""
				for name, hash := range res {
					if !c.Bool("only.hash") {
						formatted = fmt.Sprintf("%sfile: %s\thash: %s\n", formatted, name, hash)
					} else {
						formatted = fmt.Sprintf("%s%s\n", formatted, hash)
					}
				}
				return ioutil.WriteFile(c.String("save.file"), []byte(formatted), os.FileMode(0642))
			},
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "only.hash",
					Aliases: []string{"oh", "o"},
					Usage:   "whether or not to only store hash information",
					Value:   true,
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
					true,
					c.Int("email.count"),
					c.Int("emoji.count"),
					c.Int("paragraph.count"),
				)
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
