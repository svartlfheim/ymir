package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/ymir/internal/db"
	"gopkg.in/yaml.v3"
)

var happyYAML string = `#empty line to make it more readable
server:
  port: 9898

git:
  github:
    access_token: "somegithubtoken"

db:
  driver: "somedriver"
  options:
    fs:
      path: /some/fake/path.json
    postgres:
      user: "fake_user"
      password: "fakepass"
      migrator_user: "fake_mig_user"
      migrator_password: "fake_mig_pass"
      db: "fake_db_name"
      schema: "fake_schema"
      host: "fake_host"
      port: "3333"
`

var happyCfg Ymir = Ymir{
	Server: ServerConfig{
		Port: "9898",
	},
	Git: GitConfig{
		Github: GithubConfig{
			AccessToken: "somegithubtoken",
		},
	},
	Db: DbConfig{
		Driver: "somedriver",
		Options: DbOptionsConfig{
			FS: FSDbOptionsConfig{
				Path: "/some/fake/path.json",
			},
			Postgres: PostgresOptionsConfig{
				User:             "fake_user",
				Password:         "fakepass",
				MigratorUser:     "fake_mig_user",
				MigratorPassword: "fake_mig_pass",
				Database:         "fake_db_name",
				Schema:           "fake_schema",
				Host:             "fake_host",
				Port:             "3333",
			},
		},
	},
}

func Test_ConfigUnmarshalsFromYAML(t *testing.T) {
	cfg := Ymir{}
	err := yaml.Unmarshal([]byte(happyYAML), &cfg)

	assert.Nil(t, err)

	assert.Equal(t, happyCfg, cfg)
}

func Test_PostgresOptionsConfig_Getters(t *testing.T) {
	cfg := PostgresOptionsConfig{
		User:             "fake_user",
		Password:         "fakepass",
		MigratorUser:     "fake_mig_user",
		MigratorPassword: "fake_mig_pass",
		Database:         "fake_db_name",
		Schema:           "fake_schema",
		Host:             "fake_host",
		Port:             "3333",
	}

	assert.Equal(t, db.DriverPostgres, cfg.GetDriverName())
	assert.Equal(t, "fake_user", cfg.GetUsername())
	assert.Equal(t, "fakepass", cfg.GetPassword())
	assert.Equal(t, "fake_host", cfg.GetHost())
	assert.Equal(t, "3333", cfg.GetPort())
	assert.Equal(t, "fake_db_name", cfg.GetDatabase())
	assert.Equal(t, "fake_schema", cfg.GetSchema())
	assert.Equal(t, "fake_mig_user", cfg.GetMigratorUsername())
	assert.Equal(t, "fake_mig_pass", cfg.GetMigratorPassword())

	cfg.MigratorUser = ""
	cfg.MigratorPassword = ""

	// These should fall back to the User and Password properties if not supplied
	assert.Equal(t, "fake_user", cfg.GetMigratorUsername())
	assert.Equal(t, "fakepass", cfg.GetMigratorPassword())
}
