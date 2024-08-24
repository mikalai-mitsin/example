package main

import (
	"os"

	"github.com/mikalai-mitsin/example"
	"github.com/mikalai-mitsin/example/internal/pkg/containers"
	"github.com/urfave/cli/v2"
)

var (
	configPath = ""
)

func main() {
	app := &cli.App{
		Name:    example.Name,
		Usage:   "service",
		Version: example.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				EnvVars:     []string{"EXAMPLE_CONFIG_PATH"},
				TakesFile:   true,
				Value:       configPath,
				Destination: &configPath,
				HasBeenSet:  false,
			},
		},
		Action: runApp,
		Commands: []*cli.Command{
			{
				Name:      "migrate",
				Usage:     "Run migrations",
				Action:    runMigrations,
				ArgsUsage: "",
			},
			{
				Name:      "grpc",
				Usage:     "Run gRPC server",
				Action:    runGRPC,
				ArgsUsage: "",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

// runApp - run app
func runApp(context *cli.Context) error {
	app := containers.NewGRPCContainer(configPath)
	app.Run()
	return nil
}

// runGRPC - run grpc api
func runGRPC(context *cli.Context) error {
	app := containers.NewGRPCContainer(configPath)
	app.Run()
	return nil
}

// runMigrations - migrate database
func runMigrations(context *cli.Context) error {
	app := containers.NewMigrateContainer(configPath)
	app.Run()
	return nil
}
