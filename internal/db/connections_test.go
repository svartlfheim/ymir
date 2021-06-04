package db

import (
	"net"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	ymirtestdb "github.com/svartlfheim/ymir/test/db"
)

type fakeProvider struct{}

func (_ *fakeProvider) GetDriverName() string {
	return "fakedriver"
}
func (_ *fakeProvider) GetUsername() string {
	return "testuser"
}
func (_ *fakeProvider) GetPassword() string {
	return "testpass"
}
func (_ *fakeProvider) GetHost() string {
	return "testhost"
}
func (_ *fakeProvider) GetPort() string {
	return "1111"
}
func (_ *fakeProvider) GetDatabase() string {
	return "testdb"
}
func (_ *fakeProvider) GetSchema() string {
	return "testschema"
}

type fakePostgresProvider struct{}

func (_ *fakePostgresProvider) GetDriverName() string {
	return DriverPostgres
}
func (_ *fakePostgresProvider) GetUsername() string {
	return "testuser"
}
func (_ *fakePostgresProvider) GetPassword() string {
	return "testpass"
}
func (_ *fakePostgresProvider) GetHost() string {
	return "testhost"
}
func (_ *fakePostgresProvider) GetPort() string {
	return "5432"
}
func (_ *fakePostgresProvider) GetDatabase() string {
	return "testdb"
}
func (_ *fakePostgresProvider) GetSchema() string {
	return "testschema"
}

type usablePostgresProvider struct {
	driver string
	user   string
	pass   string
	host   string
	port   string
	db     string
	schema string
}

func (_ *usablePostgresProvider) GetDriverName() string {
	return DriverPostgres
}
func (p *usablePostgresProvider) GetUsername() string {
	return p.user
}
func (p *usablePostgresProvider) GetPassword() string {
	return p.pass
}
func (p *usablePostgresProvider) GetHost() string {
	return p.host
}
func (p *usablePostgresProvider) GetPort() string {
	return p.port
}
func (p *usablePostgresProvider) GetDatabase() string {
	return p.db
}
func (p *usablePostgresProvider) GetSchema() string {
	return p.schema
}

func Test_NewConnectionForPostgres_CantConnect(t *testing.T) {
	// This contains details which should not match any known db server
	p := &fakePostgresProvider{}

	conn, err := NewPostgresConnection(p)

	assert.Nil(t, conn)
	assert.IsType(t, &net.OpError{}, err)
}

func Test_NewConnectionForPostgres_Successful(t *testing.T) {
	ymirtestdb.RunTestWithPostgresDB(ymirtestdb.PostgresDbOptions{},
		t,
		func(t *testing.T, dbCfg ymirtestdb.PostgresTestDb) {
			p := &usablePostgresProvider{
				user:   dbCfg.UserName,
				pass:   dbCfg.Password,
				host:   dbCfg.Host,
				port:   dbCfg.Port,
				db:     dbCfg.DbName,
				schema: dbCfg.Schema,
			}
			conn, err := NewPostgresConnection(p)

			assert.Nil(t, err)
			assert.NotNil(t, conn)
		})
}
