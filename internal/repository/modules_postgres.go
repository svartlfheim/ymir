package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/ymir/internal/registry"
)

type postgresDbModule struct {
	Id        string `db:"id"`
	Name      string `db:"name"`
	Namespace string `db:"namespace"`
	Provider  string `db:"provider"`
}

func (pM *postgresDbModule) ToDomainModel() registry.Module {
	return registry.Module{
		Id:        pM.Id,
		Name:      pM.Name,
		Namespace: pM.Namespace,
		Provider:  pM.Provider,
	}
}

func (pM *postgresDbModule) Populate(m registry.Module) {
	pM.Id = m.Id
	pM.Name = m.Name
	pM.Namespace = m.Namespace
	pM.Provider = m.Provider
}

type postgresDbModuleVersion struct {
	Id            string         `db:"id"`
	Version       string         `db:"version"`
	ModuleId      string         `db:"module_id"`
	SourceRef     string         `db:"source_ref"`
	ArchiveId     sql.NullString `db:"archive_id"`
	RepositoryUrl string         `db:"repository_url"`
	Status        string         `db:"status"`
	EventsJSON    string         `db:"meta"`
}

func (pMV *postgresDbModuleVersion) ToDomainModel() registry.ModuleVersion {
	archiveId := ""

	if pMV.ArchiveId.Valid {
		archiveId = pMV.ArchiveId.String
	}
	return registry.ModuleVersion{
		Id:            pMV.Id,
		ModuleId:      pMV.ModuleId,
		Version:       pMV.Version,
		Source:        pMV.SourceRef,
		DownloadURL:   archiveId,
		RepositoryURL: pMV.RepositoryUrl,
		Status:        registry.VersionStatus(pMV.Status),
	}
}

func (pMV *postgresDbModuleVersion) Populate(mv registry.ModuleVersion) {
	pMV.Id = mv.Id
	pMV.ModuleId = mv.ModuleId
	pMV.RepositoryUrl = mv.RepositoryURL
	pMV.SourceRef = mv.Source
	pMV.Version = mv.Version
	pMV.Status = string(mv.Status)

	// maybe this is right?
	if mv.DownloadURL != "" {
		pMV.ArchiveId = sql.NullString{String: mv.DownloadURL, Valid: true}
	}
}

type PostgresModules struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func (s *PostgresModules) startTransaction() (*sqlx.Tx, error) {
	tx, err := s.db.Beginx()

	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, ErrDbTransaction{
			Wrapped: err,
		}
	}

	return tx, nil
}

func (s *PostgresModules) ById(id string) (m registry.Module, err error) {
	dbModule := &postgresDbModule{}
	q := fmt.Sprintf(`SELECT
	*
FROM 
	%s
WHERE
	id = $1;`,
		ModulesTableName)

	err = s.db.Get(dbModule, q, id)

	if err == sql.ErrNoRows {
		return m, registry.ErrResourceNotFound{
			Type: "Module",
			URI:  id,
		}
	} else if err != nil {
		return m, wrapQueryError(err)
	}

	return dbModule.ToDomainModel(), nil
}

func (s *PostgresModules) ByFQN(fqn registry.ModuleFQN) (m registry.Module, err error) {
	dbModule := &postgresDbModule{}
	q := fmt.Sprintf(`SELECT
	*
FROM 
	%s
WHERE
	provider = $1 AND
	namespace = $2 AND
	name = $3;`,
		ModulesTableName)

	err = s.db.Get(dbModule, q, fqn.Provider, fqn.Namespace, fqn.Name)

	if err == sql.ErrNoRows {
		return m, registry.ErrResourceNotFound{
			Type: "Module",
			URI:  fqn.String(),
		}
	} else if err != nil {
		return m, wrapQueryError(err)
	}

	return dbModule.ToDomainModel(), nil
}

func (s *PostgresModules) buildModulesFilterClause(f registry.ModuleFilters) (clause string, params map[string]interface{}) {
	if f.Namespace == "" && f.Provider == "" {
		return
	}

	clauseParts := []string{}
	params = map[string]interface{}{}

	if f.Namespace != "" {
		clauseParts = append(clauseParts, " namespace = :namespace")
		params["namespace"] = f.Namespace
	}

	if f.Provider != "" {
		clauseParts = append(clauseParts, " provider = :provider")
		params["provider"] = f.Provider
	}

	clause = "WHERE " + strings.Join(clauseParts, " AND ")

	return
}

func (s *PostgresModules) All(_ registry.ChunkingOptions, f registry.ModuleFilters) (ms []registry.Module, err error) {
	where, params := s.buildModulesFilterClause(f)
	q := fmt.Sprintf(`SELECT
	*
FROM 
	%s m
%s
ORDER BY m.provider ASC, m.namespace ASC, m.name ASC;`, ModulesTableName, where)

	rows, err := s.db.NamedQuery(q, params)

	if err != nil {
		return ms, wrapQueryError(err)
	}

	ms = []registry.Module{}

	for rows.Next() {
		dbM := &postgresDbModule{}
		err := rows.StructScan(dbM)

		if err != nil {

			// Add id ideally
			s.logger.Error().Err(err).Msg("failed to scan row")
			return []registry.Module{}, wrapHydrationError("Module", err)
		}

		ms = append(ms, dbM.ToDomainModel())
	}

	return ms, nil
}

func (s *PostgresModules) VersionById(id string) (mv registry.ModuleVersion, err error) {
	dbModuleVersion := &postgresDbModuleVersion{}
	q := fmt.Sprintf(`
SELECT
	*
FROM 
	%s
WHERE
	id = $1;`,
		ModuleVersionsTableName)

	err = s.db.Get(dbModuleVersion, q, id)

	if err == sql.ErrNoRows {
		return mv, registry.ErrResourceNotFound{
			Type: "ModuleVersion",
			URI:  id,
		}
	} else if err != nil {
		return mv, wrapQueryError(err)
	}

	return dbModuleVersion.ToDomainModel(), nil
}

func (s *PostgresModules) VersionsByModule(moduleId string, _ registry.ChunkingOptions) (mVs []registry.ModuleVersion, err error) {
	q := fmt.Sprintf(`
SELECT
	*
FROM 
	%s
WHERE
	module_id = $1;`,
		ModuleVersionsTableName)

	rows, err := s.db.Queryx(q, moduleId)

	if err != nil {
		return mVs, wrapQueryError(err)
	}

	mVs = []registry.ModuleVersion{}

	for rows.Next() {
		dbM := &postgresDbModuleVersion{}
		err := rows.StructScan(dbM)

		if err != nil {

			// Add id ideally
			s.logger.Error().Err(err).Msg("failed to scan row")
			return []registry.ModuleVersion{}, wrapHydrationError("Module", err)
		}

		mVs = append(mVs, dbM.ToDomainModel())
	}

	return mVs, nil
}

func (s *PostgresModules) VersionsByModuleFQN(fqn registry.ModuleFQN, _ registry.ChunkingOptions) (mVs []registry.ModuleVersion, err error) {
	q := fmt.Sprintf(`
SELECT
	*
FROM 
	%s
WHERE
	module_id = (
		SELECT
			id
		FROM
			%s
		WHERE
			provider = $1 AND
			namespace = $2 AND
			name = $3
		);`,
		ModuleVersionsTableName, ModulesTableName)

	rows, err := s.db.Queryx(q, fqn.Provider, fqn.Namespace, fqn.Name)

	if err != nil {
		return mVs, wrapQueryError(err)
	}

	mVs = []registry.ModuleVersion{}

	for rows.Next() {
		dbM := &postgresDbModuleVersion{}
		err := rows.StructScan(dbM)

		if err != nil {

			// Add id ideally
			s.logger.Error().Err(err).Msg("failed to scan row")
			return []registry.ModuleVersion{}, wrapHydrationError("Module", err)
		}

		mVs = append(mVs, dbM.ToDomainModel())
	}

	return mVs, nil
}

func (s *PostgresModules) VersionByModuleAndValue(moduleId string, version string) (mv registry.ModuleVersion, err error) {
	dbModuleVersion := &postgresDbModuleVersion{}
	q := fmt.Sprintf(`
SELECT
	*
FROM 
	%s
WHERE
	module_id = $1 AND
	version = $2;
`, ModuleVersionsTableName)

	err = s.db.Get(dbModuleVersion, q, moduleId, version)

	if err == sql.ErrNoRows {
		return mv, registry.ErrResourceNotFound{
			Type: "ModuleVersion",
			URI:  fmt.Sprintf("%s@%s", moduleId, version),
		}
	} else if err != nil {
		return mv, wrapQueryError(err)
	}

	return dbModuleVersion.ToDomainModel(), nil
}

func (s *PostgresModules) VersionByFQN(fqn registry.ModuleVersionFQN) (mv registry.ModuleVersion, err error) {
	dbModuleVersion := &postgresDbModuleVersion{}
	q := fmt.Sprintf(`
SELECT
	*
FROM 
	%s AS mv
WHERE
	mv.version = $1 AND
	mv.module_id = (
		SELECT
			id
		from 
			%s AS m 
		WHERE
			name = $2 AND 
			namespace = $3 AND 
			provider = $4
	);
`, ModuleVersionsTableName, ModulesTableName)

	err = s.db.Get(dbModuleVersion, q, fqn.Version, fqn.ModuleFQN.Name, fqn.ModuleFQN.Namespace, fqn.ModuleFQN.Provider)

	if err == sql.ErrNoRows {
		return mv, registry.ErrResourceNotFound{
			Type: "ModuleVersion",
			URI:  fqn.String(),
		}
	} else if err != nil {
		s.logger.Info().Str("q", q).Msg("query")
		return mv, wrapQueryError(err)
	}

	return dbModuleVersion.ToDomainModel(), nil
}

func (s *PostgresModules) AddModule(mod registry.Module) (m registry.Module, err error) {
	tx, err := s.startTransaction()

	if err != nil {
		return m, err
	}

	dbModule := &postgresDbModule{}
	dbModule.Populate(mod)

	insert := fmt.Sprintf(`
INSERT INTO %s (id, name, namespace, provider) VALUES (:id, :name, :namespace, :provider);`,
		ModulesTableName)

	_, err = tx.NamedExec(insert, dbModule)

	if err != nil {
		return m, wrapTransactionError(err)
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			return m, wrapTransactionError(rollbackErr)
		}

		return m, wrapTransactionError(err)
	}

	m, err = s.ById(mod.Id)

	if err != nil {
		if _, ok := err.(registry.ErrResourceNotFound); !ok {
			return m, wrapTransactionError(errors.New("module was not persisted"))
		}

		return m, wrapQueryError(err)
	}

	return m, nil
}

func (s *PostgresModules) DeleteModule(mod registry.Module) (err error) {
	tx, err := s.startTransaction()

	if err != nil {
		return err
	}

	delete := fmt.Sprintf(`
DELETE FROM %s WHERE id = $1`,
		ModulesTableName)

	_, err = tx.Exec(delete, mod.Id)

	if err != nil {
		return wrapTransactionError(err)
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			return wrapTransactionError(rollbackErr)
		}

		return wrapTransactionError(err)
	}

	return nil
}

func (s *PostgresModules) AddVersion(new registry.ModuleVersion) (v registry.ModuleVersion, err error) {
	tx, err := s.startTransaction()

	if err != nil {
		return v, err
	}

	insert := fmt.Sprintf(`
INSERT INTO %s (
	id, 
	version,
	source_ref,
	archive_id,
	repository_url,
	status,
	module_id
) VALUES (
	:id, 
	:version, 
	:source_ref, 
	NULL,
	:repository_url,
	:status,
	:module_id
);`,
		ModuleVersionsTableName)

	dbVModule := &postgresDbModuleVersion{}
	dbVModule.Populate(new)
	dbVModule.Status = string(registry.VersionStatuses.Pending)
	dbVModule.ArchiveId = sql.NullString{}

	_, err = tx.NamedExec(insert, dbVModule)

	if err != nil {
		return v, wrapTransactionError(err)
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			return v, wrapTransactionError(rollbackErr)
		}

		return v, wrapTransactionError(err)
	}

	v, err = s.VersionById(new.Id)

	if err != nil {
		if _, ok := err.(registry.ErrResourceNotFound); !ok {
			return v, wrapTransactionError(errors.New("module version was not persisted"))
		}

		return v, wrapQueryError(err)
	}

	return v, nil
}

func (s *PostgresModules) DeleteVersionsForModule(mod registry.Module) (err error) {
	tx, err := s.startTransaction()

	if err != nil {
		return err
	}

	delete := fmt.Sprintf(`
DELETE FROM %s WHERE module_id = $1`,
		ModuleVersionsTableName)

	_, err = tx.Exec(delete, mod.Id)

	if err != nil {
		return wrapTransactionError(err)
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			return wrapTransactionError(rollbackErr)
		}

		return wrapTransactionError(err)
	}

	return nil
}

func (s *PostgresModules) DeleteModuleVersion(mv registry.ModuleVersion) error {
	tx, err := s.startTransaction()

	if err != nil {
		return err
	}

	delete := fmt.Sprintf(`
DELETE FROM %s WHERE id = $1`,
		ModuleVersionsTableName)

	_, err = tx.Exec(delete, mv.Id)

	if err != nil {
		return wrapTransactionError(err)
	}

	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()

		if rollbackErr != nil {
			return wrapTransactionError(rollbackErr)
		}

		return wrapTransactionError(err)
	}

	return nil
}
