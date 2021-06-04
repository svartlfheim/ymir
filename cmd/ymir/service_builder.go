package ymir

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/svartlfheim/clapp"
	"github.com/svartlfheim/ymir/internal/cli"
	"github.com/svartlfheim/ymir/internal/config"
	"github.com/svartlfheim/ymir/internal/db"
	"github.com/svartlfheim/ymir/internal/output"
	"github.com/svartlfheim/ymir/internal/registry"
	"github.com/svartlfheim/ymir/internal/repository"
	"github.com/svartlfheim/ymir/internal/server"
)

func buildModuleRepository(cfg *config.Ymir, ctx context.Context, l zerolog.Logger) (registry.ModuleRepository, error) {
	switch cfg.Db.Driver {
	case string(repository.PostgresDriver):
		conn, err := db.NewPostgresConnection(cfg.Db.Options.Postgres)

		if err != nil {
			return nil, err
		}

		return repository.BuildModulesForPostgres(conn, clapp.LoggerFromContext(ctx)), nil
	default:
		return nil, repository.ErrDriverNotImplemented{
			Driver: cfg.Db.Driver,
		}
	}
}

func buildAuditor(cfg *config.Ymir, ctx context.Context, l zerolog.Logger) (*registry.Auditor, error) {
	var repo server.AuditLogRepository
	switch cfg.Db.Driver {
	case string(repository.PostgresDriver):
		conn, err := db.NewPostgresConnection(cfg.Db.Options.Postgres)

		if err != nil {
			return nil, err
		}

		repo = repository.BuildAuditLogsForPostgres(conn, clapp.LoggerFromContext(ctx))
	default:
		return nil, repository.ErrDriverNotImplemented{
			Driver: cfg.Db.Driver,
		}
	}

	return registry.NewAuditor(repo, l), nil

}

func buildTableFactory() *output.TableFactory {
	return output.NewTableFactory(os.Stdout)
}

func buildCommandBus(c YmirCommand) *registry.CommandBus {
	l := c.GetLogger()
	ctx := c.cobra.Context()

	moduleRepo, err := buildModuleRepository(c.GetConfig(), ctx, l)

	if err != nil {
		l.Fatal().Err(err).Msg("failed to build module repo")
	}

	cb := registry.NewCommandBus(
		registry.WithFS(clapp.FsFromContext(ctx)),
		registry.WithModuleRepo(moduleRepo),
		registry.WithLogger(l),
		registry.WithPrompter(cli.NewPrompter()),
		registry.WithCommandValidatorBuilder(registry.NewCommandValidator),
	)

	return cb
}
