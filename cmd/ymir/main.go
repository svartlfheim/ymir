package ymir

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/svartlfheim/clapp"
	"github.com/svartlfheim/ymir/internal/config"
	"github.com/svartlfheim/ymir/internal/output"
	"github.com/svartlfheim/ymir/internal/repository"
	"github.com/svartlfheim/ymir/pkg/gopoint"
)

var versionHashRef string = "unknown"
var version string = "unknown"

type YmirCommandHandlerFunc func(cmd YmirCommand) error

type YmirCommand struct {
	cobra *cobra.Command
	args  []string
}

func (c *YmirCommand) GetConfig() *config.Ymir {
	return clapp.ConfigFromContext(c.cobra.Context()).(*config.Ymir)
}

func (c *YmirCommand) GetLogger() zerolog.Logger {
	return clapp.LoggerFromContext(c.cobra.Context())
}

func (c *YmirCommand) GetOutput() *output.Fmt {
	return output.NewFmt(os.Stdout)
}

func (c *YmirCommand) GetArg(i int, def string) string {
	if len(c.args) > i {
		return c.args[i]
	}

	return def
}

func (c *YmirCommand) GetRequiredArg(i int) (string, error) {
	if len(c.args) > i {
		return c.args[i], nil
	}

	return "", ErrNoArgAtIndex{
		index: i,
	}
}

func (c *YmirCommand) GetAllArgs() []string {
	return c.args
}

func configPath() string {
	if val, found := os.LookupEnv("YMIR_CONFIG_FILE"); found {
		return val
	}

	if homeDir, err := os.UserHomeDir(); err != nil {
		return fmt.Sprintf("%s/ymir.yaml", homeDir)
	}

	return "./ymir.yaml"
}

func buildHandler(f YmirCommandHandlerFunc) clapp.HandlerFunc {
	return func(c *cobra.Command, args []string) error {
		ymirCommand := YmirCommand{
			cobra: c,
			args:  args,
		}

		return f(ymirCommand)
	}

}

var ymirConfig config.Ymir = config.Ymir{
	Server: config.ServerConfig{
		Port: "8888",
	},
	Db: config.DbConfig{
		Driver: string(repository.PostgresDriver),
		Options: config.DbOptionsConfig{
			Postgres: config.PostgresOptionsConfig{
				Database: "postgres",
				Schema:   "ymir",
				Host:     "postgres",
				Port:     "5432",
			},
		},
	},
	Git: config.GitConfig{
		Github: config.GithubConfig{
			AccessToken: "",
		},
	},
}

var app clapp.App = clapp.App{
	Config:     &ymirConfig,
	ConfigPath: configPath(),
	Fs:         afero.NewOsFs(),
	Logger:     zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel),
	RootCommand: clapp.Command{
		Name: "ymir",
		Descriptions: clapp.Descriptions{
			Short: "Ymir is a self-hosted terraform registry.",
			Long: `Ymir provides an API, and CLI interface to manage terraform modules. 
It allows the creation of new modules and versions of modules to be distributed from this registry.

Module source can be fetched from a configured git repository, and a resulting archive is stored in the chosen storage system. This archive will be served via the module registry protocol when requested.`,
		},
		Handle: buildHandler(func(cmd YmirCommand) error {
			cfg := clapp.ConfigFromContext(cmd.cobra.Context())

			json, err := json.MarshalIndent(cfg, "", "\t")

			if err != nil {
				fmt.Print("Error marshalling to JSON!\n")

				return err
			}

			fmt.Printf("%s\n", string(json))

			return nil
		}),
		Children: []clapp.Command{
			{
				Name: "version",
				Descriptions: clapp.Descriptions{
					Short: "Show the version.",
					Long:  `Shows the version of the app, and the commit hash it was built from.`,
				},
				Handle: func(_ *cobra.Command, _ []string) error {
					fmt.Printf("Version: %s (%s)\n", version, versionHashRef)
					return nil
				},
			},
			{
				Name: "serve",
				Descriptions: clapp.Descriptions{
					Short: "Run the HTTP server",
					Long: `Runs an HTTP server, that exposes the endpoints necessary acccording to the module registry protocol defined by terraform.
					
Coming soon, is a UI to manage module versions.`,
				},
				LocalFlags: []clapp.Flag{
					{
						Name:        "port",
						Short:       "p",
						Description: "The port that the http server should run on.",
						ValueRef:    &ymirConfig.Server.Port,
						Type:        clapp.StringFlag,
						Required:    false,
					},
				},
				Handle: buildHandler(serve),
			},
			{
				Name: "module",
				Descriptions: clapp.Descriptions{
					Short: "Contains commands associated with managing modules.",
					Long: `See help for available commands.
					
These will allow you to manage modules from the cli.`,
				},
				Children: []clapp.Command{
					{
						Name:   "add",
						Handle: buildHandler(module_add),
						Descriptions: clapp.Descriptions{
							Short: "Add a new module.",
							Long: `A module can be added interactively or from a supplied file, which will be parsed as json.
If the file option is not supplied, you will be prompted for the values to fields.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "file",
								Short:       "f",
								Description: "Location of a json file to import.",
								ValueRef:    gopoint.ToString(""),
								Type:        clapp.StringFlag,
								Required:    false,
							},
						},
					},
					{
						Name:   "list",
						Handle: buildHandler(module_list),
						Descriptions: clapp.Descriptions{
							Short: "List all of the available modules.",
							Long: `Output a list of all available modules.
Output can be tabular, or JSON depending on options provided.

Modules can be filtered by provider, and/or namespace.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "output",
								Short:       "o",
								Description: "The output style to use, one of: json, table. Default: table",
								ValueRef:    gopoint.ToString(""),
								Required:    false,
								Type:        clapp.StringFlag,
							},
							{
								Name:        "provider",
								Short:       "p",
								Description: "List only modules for the supplied provider. Works inclusively with namespace.",
								ValueRef:    gopoint.ToString(""),
								Required:    false,
								Type:        clapp.StringFlag,
							},
							{
								Name:        "namespace",
								Short:       "n",
								Description: "List only modules for the supplied namespace. Works inclusively with provider.",
								ValueRef:    gopoint.ToString(""),
								Required:    false,
								Type:        clapp.StringFlag,
							},
						},
					},
					{
						Name:   "show",
						Handle: buildHandler(module_show),
						Descriptions: clapp.Descriptions{
							Short: "Show details about a single module.",
							Long: `Show information about a single module, defined by an ID or a ModuleFQN.
An ID or an FQN must be supplied.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "output",
								Short:       "o",
								Description: "The output style to use, one of: json, table. Default: table",
								ValueRef:    gopoint.ToString(""),
								Required:    false,
								Type:        clapp.StringFlag,
							},
						},
					},
					{
						Name:   "delete",
						Handle: buildHandler(module_delete),
						Descriptions: clapp.Descriptions{
							Short: "Delete a module.",
							Long: `Deletes a module by it's ID, or ModuleFQN.
You will be prompted for confirmation unless the force option is supplied.

All versions related to this module will be deleted.

Danger, Will Robinson!`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "yes",
								Short:       "y",
								Description: "Whether to force deletion, without prompt.",
								ValueRef:    gopoint.ToBool(false),
								Required:    false,
								Type:        clapp.BoolFlag,
							},
							{
								Name:        "versions",
								Description: "Whether to delete all versions of the module. Delete will fail if versions exist and this option isn't supplied.",
								ValueRef:    gopoint.ToBool(false),
								Required:    false,
								Type:        clapp.BoolFlag,
							},
						},
					},
				},
			},
			{
				Name: "module-version",
				Descriptions: clapp.Descriptions{
					Short: "Contains commands associated with managing module versions.",
					Long: `See help for available commands.
					
These will allow you to manage module versions from the cli.`,
				},
				Children: []clapp.Command{
					{
						Name:   "add",
						Handle: buildHandler(module_version_add),
						Descriptions: clapp.Descriptions{
							Short: "Add a new version to the specified module.",
							Long: `A module ID or ModuleFQN must be supplied.

You can add a version interactively, or by specifying a json file.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "file",
								Short:       "f",
								Description: "Location of a json file to import.",
								ValueRef:    gopoint.ToString(""),
								Type:        clapp.StringFlag,
								Required:    false,
							},
						},
					},
					{
						Name:   "list",
						Handle: buildHandler(module_version_list),
						Descriptions: clapp.Descriptions{
							Short: "List all of the available versions of a module.",
							Long: `Output a list of all available versions for a module.
A module ID or ModuleFQN must be provided.

Output can be tabular, or JSON depending on options provided.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "output",
								Short:       "o",
								Description: "The output style to use, one of: json, table. Default: table",
								ValueRef:    gopoint.ToString(""),
								Required:    false,
								Type:        clapp.StringFlag,
							},
						},
					},
					{
						Name:   "show",
						Handle: buildHandler(module_version_show),
						Descriptions: clapp.Descriptions{
							Short: "Show details about a single module.",
							Long: `Show information about a single module, defined by an ID or a ModuleFQN.
An ID or an FQN must be supplied.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "output",
								Short:       "o",
								Description: "The output style to use, one of: json, table. Default: table",
								ValueRef:    gopoint.ToString(""),
								Required:    false,
								Type:        clapp.StringFlag,
							},
						},
					},
					{
						Name:   "delete",
						Handle: buildHandler(module_version_delete),
						Descriptions: clapp.Descriptions{
							Short: "Delete a version of a module.",
							Long: `Deletes a single module, defined by an ID or a ModuleFQN.

A ModuleVersion ID or ModuleVersionFQN must be supplied.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "yes",
								Short:       "y",
								Description: "Whether to force deletion, without prompt.",
								ValueRef:    gopoint.ToBool(false),
								Required:    false,
								Type:        clapp.BoolFlag,
							},
						},
					},
				},
			},
			{
				Name: "migrate",
				Descriptions: clapp.Descriptions{
					Short: "Commands to migrate the database schema for the configured driver.",
					Long: `This wil only be applicable when using an external database driver (e.g. postgres). 
The filesystem, and inmemory drivers cannot be migrated.

See subcommands.`,
				},
				// PersistentFlags: []clapp.Flag{
				// 	{
				// 		Name: "migrator",
				// 		Short: "m",
				// 		Description: "The name of the user/system performing the migration.",
				// 		Required: true,
				// 		ValueRef: gopoint.ToString(""),
				// 		Type: clapp.StringFlag,
				// 	},
				// },
				Children: []clapp.Command{
					{
						Name: "up",
						Descriptions: clapp.Descriptions{
							Short: "Apply pending migrations up to the defined target.",
							Long: `You must supply the name of the migration to go to.

Alternatiely you can supply the --all flag to run all.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "migrator",
								Short:       "m",
								Description: "The name of the user/system performing the migration.",
								Required:    true,
								ValueRef:    gopoint.ToString(""),
								Type:        clapp.StringFlag,
							},
						},
						Handle: buildHandler(migrateUp),
					},
					{
						Name: "down",
						Descriptions: clapp.Descriptions{
							Short: "Rollback migrations to the supplied target.",
							Long: `The supplied target will also be rolled back.
The target must be supplied, to prevent accidental data loss.

The option --all can be passed to completely reset the database. The tables to track migrations themselves will remain.`,
						},
						LocalFlags: []clapp.Flag{
							{
								Name:        "all",
								Description: "Tells the migrator to rollback all migrations.",
								ValueRef:    gopoint.ToBool(false),
								Type:        clapp.BoolFlag,
							},
							{
								Name:        "migrator",
								Short:       "m",
								Description: "The name of the user/system performing the migration.",
								Required:    true,
								ValueRef:    gopoint.ToString(""),
								Type:        clapp.StringFlag,
							},
						},
						Handle: buildHandler(migrateDown),
					},
					{
						Name: "list",
						Descriptions: clapp.Descriptions{
							Short: "List all migrations and their current state.",
							Long:  ``,
						},
						Handle: buildHandler(migrateList),
					},
				},
			},
		},
	},
}

func Execute() error {
	return clapp.Run(app, clapp.NewCobraExecutor())
}
