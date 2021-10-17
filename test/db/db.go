package ymirtestdb

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	ymirtestdocker "github.com/svartlfheim/ymir/test/docker"
	ymirtestsuites "github.com/svartlfheim/ymir/test/suites"
)

type PrepareDbFunc func(*testing.T, PostgresTestDb)

type PostgresDbOptions struct {
	UserName    string
	Password    string
	DbName      string
	HostName    string
	Version     string
	KillTimeout int
	Prepare PrepareDbFunc
}

type PostgresTestDb struct {
	Host     string
	Port     string
	Schema   string
	UserName string
	Password string
	DbName   string
}

func addDefaultPostgresOpts(opts PostgresDbOptions) PostgresDbOptions {
	if opts.Version == "" {
		opts.Version = "13"
	}

	if opts.KillTimeout == 0 {
		opts.KillTimeout = 120
	}

	if opts.HostName == "" {
		opts.HostName = "ymir_test_db"
	}

	if opts.UserName == "" {
		opts.UserName = "testuser"
	}

	if opts.Password == "" {
		opts.Password = "testpass"
	}

	if opts.DbName == "" {
		opts.DbName = "testdb"
	}

	return opts
}

func DbConnFromPostgresDb(t *testing.T, dbCfg PostgresTestDb) *sqlx.DB {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s search_path=%s sslmode=disable",
		dbCfg.UserName,
		dbCfg.Password,
		dbCfg.DbName,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Schema,
	)

	conn, err := sqlx.Connect("postgres", connString)

	if err != nil {
		t.Error("could not initialise db connection")
		t.FailNow()
	}

	return conn
}

func RunTestWithPostgresDB(dbOpts PostgresDbOptions, t *testing.T, toRun func(*testing.T, PostgresTestDb)) {
	ymirtestsuites.SkipIfIntegrationNotEnabled(t)

	dbOpts = addDefaultPostgresOpts(dbOpts)

	var dbConn *sql.DB
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	nw, err := ymirtestdocker.FindTestNetwork(pool.Client)

	if err != nil {
		log.Fatalf("Could not connect find test network: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Networks: []*dockertest.Network{
			{
				Network: nw,
			},
		},
		Hostname:   dbOpts.HostName,
		Repository: "postgres",
		Tag:        dbOpts.Version,
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbOpts.Password),
			fmt.Sprintf("POSTGRES_USER=%s", dbOpts.UserName),
			fmt.Sprintf("POSTGRES_DB=%s", dbOpts.DbName),
			"listen_addresses = '*'",
		},
		Cmd: []string{
			// Useful for debugging tests
			"postgres",
			"-c",
			"log_connections=ON",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbOpts.UserName, dbOpts.Password, dbOpts.HostName, "5432", dbOpts.DbName)

	// log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(uint(dbOpts.KillTimeout)) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = time.Duration(dbOpts.KillTimeout) * time.Second
	if err = pool.Retry(func() error {
		dbConn, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return dbConn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	
	testDb := PostgresTestDb{
		Host:     dbOpts.HostName,
		Port:     "5432",
		Schema:   "public",
		UserName: dbOpts.UserName,
		Password: dbOpts.Password,
		DbName:   dbOpts.DbName,
	}

	if dbOpts.Prepare != nil {
		dbOpts.Prepare(t, testDb)
	}

	//Run tests
	toRun(t, testDb)

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
