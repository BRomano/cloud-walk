package main

import (
	"cloud-walk/internal/domain/service"
	"cloud-walk/internal/infra/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func loadFile(filePath string) ([]byte, error) {
	loggerContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return loggerContent, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := service.NewLogParser(repository.LogParserFactory)

	app := &cli.App{
		Name:  "cli command to walk-cloud parse log",
		Usage: "help for more",
		Commands: []*cli.Command{
			{
				Name:    "statistic",
				Aliases: []string{"s"},
				Usage:   "print statistics from the given file",
				Action: func(cCtx *cli.Context) error {
					file, err := loadFile(cCtx.Args().First())
					if err != nil {
						return err
					}
					statistics, err := service.GetMatchesStatistics(ctx, repository.Quake3Arena, file)
					if err != nil {
						return err
					}
					prettyJsonn, err := json.MarshalIndent(statistics, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(prettyJsonn))
					return nil
				},
			},
			{
				Name:    "death",
				Aliases: []string{"d"},
				Usage:   "print deths from the given file",
				Action: func(cCtx *cli.Context) error {
					file, err := loadFile(cCtx.Args().First())
					if err != nil {
						return err
					}
					statistics, err := service.GetKillsByMeans(ctx, repository.Quake3Arena, file)
					if err != nil {
						return err
					}
					prettyJsonn, err := json.MarshalIndent(statistics, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(prettyJsonn))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
