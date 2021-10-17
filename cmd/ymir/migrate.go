package ymir

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/ymir/internal/db"
	"github.com/svartlfheim/ymir/internal/output"
)

type migrator interface {
	Up(t string) error
	Down(t string) error
	ListMigrations() ([]*gomigrator.MigrationRecord, error)
}

var postgresMigrations gomigrator.MigrationList = gomigrator.NewMigrationList(
	[]gomigrator.Migration{
		{
			Id:   "create-modules-table",
			Name: "create modules table",
			Execute: func(tx *sqlx.Tx) (sql.Result, error) {
				createTable := `CREATE TABLE modules(
	id uuid NOT NULL,
	name TEXT NOT NULL,
	namespace TEXT NOT NULL,
	provider TEXT NOT NULL,
	PRIMARY KEY(id),
	UNIQUE(name, namespace, provider)
);
CREATE INDEX idx_modules_full_name ON modules(name, namespace, provider);`

				return tx.Exec(createTable)
			},
			Rollback: func(tx *sqlx.Tx) (sql.Result, error) {
				dropTable := `DROP TABLE modules;`

				return tx.Exec(dropTable)
			},
		},
		{
			Id:   "create-module-versions-table",
			Name: "create module versions table",
			Execute: func(tx *sqlx.Tx) (sql.Result, error) {
				createTable := `CREATE TABLE module_versions(
	id uuid NOT NULL,
	version TEXT NOT NULL,
	module_id uuid NOT NULL,
	source_ref TEXT NOT NULL,
	archive_id TEXT DEFAULT NULL,
	repository_url TEXT NOT NULL,
	status TEXT NOT NULL,
	meta JSON DEFAULT '{}'::json,
	PRIMARY KEY(id),
	UNIQUE(version, module_id),
	CONSTRAINT fk_module FOREIGN KEY(module_id) REFERENCES modules(id)
);
CREATE INDEX idx_module_versions_status ON module_versions(status);
CREATE INDEX idx_module_versions_module_id ON module_versions(module_id);
CREATE INDEX idx_module_versions_module_id_and_version ON module_versions(module_id, version);`

				return tx.Exec(createTable)
			},
			Rollback: func(tx *sqlx.Tx) (sql.Result, error) {
				dropTable := `DROP TABLE module_versions;`

				return tx.Exec(dropTable)
			},
		},
		{
			Id:   "create-module-audit-log-table",
			Name: "create module audit log table",
			Execute: func(tx *sqlx.Tx) (sql.Result, error) {
				createTable := `CREATE TABLE audit_logs(
	id uuid NOT NULL,
	action TEXT NOT NULL,
	response_status TEXT NOT NULL,
	occurred_at timestamp with time zone,
	meta JSONB DEFAULT '{}'::jsonb,
	PRIMARY KEY(id)
);
CREATE INDEX idx_audit_logs_response_status ON audit_logs(response_status);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);`

				return tx.Exec(createTable)
			},
			Rollback: func(tx *sqlx.Tx) (sql.Result, error) {
				dropTable := `DROP TABLE audit_logs;`

				return tx.Exec(dropTable)
			},
		},
	},
)

func shouldMigrateAll(c YmirCommand) bool {
	val, err := c.cobra.LocalFlags().GetBool("all")

	if err != nil {
		// Flag didn't exist? weird...
		l := c.GetLogger()

		l.Error().Err(err).Msg("flag all was not defined on the migrate commands")
		return false
	}

	return val
}

func buildMigrator(c YmirCommand, applyer string) (migrator, error) {
	cfg := c.GetConfig()
	// Pretty sure there is a better way...
	// I'll sort this later, just wanted to test the actual migrations
	provider := cfg.Db.Options.Postgres
	provider.User = cfg.Db.Options.Postgres.GetMigratorUsername()
	provider.Password = cfg.Db.Options.Postgres.GetMigratorPassword()
	conn, err := db.NewPostgresConnection(provider)

	if err != nil {
		return nil, err
	}

	l := c.GetLogger()

	return gomigrator.NewMigrator(conn, postgresMigrations, gomigrator.Opts{
		Schema: cfg.Db.Options.Postgres.Schema,
		Applyer: applyer,
	}, l)
}

func migrateUp(c YmirCommand) error {
	l := c.GetLogger()
	applyer, err := c.cobra.LocalFlags().GetString("migrator")

	if err != nil {
		l.Error().Err(err).Msg("the migrator flag is required")
		return err
	}

	if applyer == "" {
		l.Error().Msg("the migrator flag cannot be an empty string")
		return errors.New("the migrator flag is required")
	}

	migrator, err := buildMigrator(c, applyer)

	if err != nil {
		return err
	}

	migrator.Up(c.GetArg(0, gomigrator.MigrateToLatest))
	return nil
}

func migrateDown(c YmirCommand) error {
	l := c.GetLogger()
	applyer, err := c.cobra.LocalFlags().GetString("migrator")

	if err != nil {
		l.Error().Err(err).Msg("the migrator flag is required")
		return err
	}

	if applyer == "" {
		l.Error().Msg("the migrator flag cannot be an empty string")
		return errors.New("the migrator flag is required")
	}

	migrator, err := buildMigrator(c, applyer)

	if err != nil {
		return err
	}

	migrateAll := shouldMigrateAll(c)

	target, err := c.GetRequiredArg(0)

	if err != nil && !migrateAll {
		return errors.New("target or --all must be supplied to down")
	} else if err != nil && migrateAll {
		target = gomigrator.MigrateToNothing
	}

	migrator.Down(target)
	return nil
}

func migrateList(c YmirCommand) error {
	migrator, err := buildMigrator(c, "list-command")

	if err != nil {
		return err
	}

	migs, err := migrator.ListMigrations()

	if err != nil {
		return err
	}

	rows := [][]string{}

	for _, m := range migs {
		rows = append(rows, []string{m.Id, m.Status})
	}

	tf := buildTableFactory()

	err = tf.CreateAndPrint([]string{"Migration", "Status"}, rows, output.WithIndexColumn())

	return err
}
