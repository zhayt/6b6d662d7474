package mssql

import (
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDatabase() (testcontainers.Container, *sqlx.DB, error) {
	// Create MSSQL Server container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/mssql/server:2019-latest",
		ExposedPorts: []string{"1433/tcp"},
		WaitingFor:   wait.ForListeningPort("1433/tcp"),
		Env: map[string]string{
			"ACCEPT_EULA":     "Y",
			"SA_PASSWORD":     "Test123!",
			"MSSQL_PID":       "Express",
			"MSSQL_COLLATION": "Cyrillic_General_CI_AS",
		},
	}

	// Start MSSQL Server container
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return dbContainer, nil, err
	}

	// Get host and port of MSSQL Server container
	port, err := dbContainer.MappedPort(context.Background(), "1433")
	if err != nil {
		return dbContainer, nil, err
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer, nil, err
	}

	// Create db connection string and connect
	connString := fmt.Sprintf("sqlserver://sa:Test123!@%s:%s?database=master", host, port.Port())

	db, err := sqlx.Connect("sqlserver", connString)
	if err != nil {
		return dbContainer, nil, err
	}

	// Create test database
	_, err = db.Exec("CREATE DATABASE test_db")
	if err != nil {
		return dbContainer, nil, err
	}

	// Connect to the test database
	connString = fmt.Sprintf("sqlserver://sa:Test123!@%s:%s?database=test_db", host, port.Port())

	db, err = sqlx.Connect("sqlserver", connString)
	if err != nil {
		return dbContainer, nil, err
	}

	qr := `CREATE TABLE R_CURRENCY (
    ID INT IDENTITY(1,1) PRIMARY KEY,
    TITLE nvarchar(60) NOT NULL,
    CODE nvarchar(3) NOT NULL,
    VALUE numeric(18,2) NOT NULL,
    A_DATE date not null
);`

	_, err = db.Exec(qr)
	if err != nil {
		return dbContainer, nil, err
	}
	return dbContainer, db, nil
}
