package repository

const ModulesTableName = "modules"
const AuditLogsTableName = "audit_logs"
const ModuleVersionsTableName = "module_versions"

type DbDriver string

const PostgresDriver DbDriver = "postgres"

func wrapTransactionError(e error) ErrDbTransaction {
	return ErrDbTransaction{
		Wrapped: e,
	}
}

func wrapQueryError(e error) ErrDbQuery {
	return ErrDbQuery{
		Wrapped: e,
	}
}

func wrapHydrationError(t string, e error) ErrDbHydration {
	return ErrDbHydration{
		Type:    t,
		Wrapped: e,
	}
}
