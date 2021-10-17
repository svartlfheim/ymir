package config

import (
	"github.com/svartlfheim/ymir/internal/db"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type FSDbOptionsConfig struct {
	Path string `yaml:"path"`
}

type PostgresOptionsConfig struct {
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	MigratorUser     string `yaml:"migrator_user" split_words:"true"`
	MigratorPassword string `yaml:"migrator_password" split_words:"true"`
	Database         string `yaml:"db"`
	Schema           string `yaml:"schema"`
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
}

func (c PostgresOptionsConfig) GetDriverName() string {
	return db.DriverPostgres
}

func (c PostgresOptionsConfig) GetUsername() string {
	return c.User
}

func (c PostgresOptionsConfig) GetPassword() string {
	return c.Password
}

func (c PostgresOptionsConfig) GetHost() string {
	return c.Host
}

func (c PostgresOptionsConfig) GetPort() string {
	return c.Port
}

func (c PostgresOptionsConfig) GetDatabase() string {
	return c.Database
}

func (c PostgresOptionsConfig) GetSchema() string {
	return c.Schema
}

func (c PostgresOptionsConfig) GetMigratorUsername() string {
	if c.MigratorUser != "" {
		return c.MigratorUser
	}

	return c.User
}

func (c PostgresOptionsConfig) GetMigratorPassword() string {
	if c.MigratorPassword != "" {
		return c.MigratorPassword
	}

	return c.Password
}

// type postgresConnectionProvider interface {
// 	Username() string
// 	Password() string
// 	Host() string
// 	Port() string
// 	Database() string
// 	Schema() string
// 	MigratorUsername() string
// 	MigratorPassword() string
// }

type DbOptionsConfig struct {
	FS       FSDbOptionsConfig     `yaml:"fs"`
	Postgres PostgresOptionsConfig `yaml:"postgres"`
}

type GithubConfig struct {
	AccessToken string `yaml:"access_token" split_words:"true"`
}

type GitConfig struct {
	Github GithubConfig `yaml:"github"`
}

type DbConfig struct {
	Driver  string          `yaml:"driver"`
	Options DbOptionsConfig `yaml:"options"`
}

type Ymir struct {
	Server ServerConfig `yaml:"server"`
	Db     DbConfig     `yaml:"db"`
	Git    GitConfig    `yaml:"git"`
}
