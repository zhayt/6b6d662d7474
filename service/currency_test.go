package service

import (
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zhayt/kmf-tt/storage"
	"github.com/zhayt/kmf-tt/storage/mssql"
	"go.uber.org/zap"
	"log"
	"testing"
)

func TestCurrencyService_SaveCurrency(t *testing.T) {
	type fields struct {
		currency    *storage.Storage
		l           *zap.Logger
		externalAPI string
	}
	type args struct {
		ctx  context.Context
		date string
	}

	dbContainer, db, err := SetupTestDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer dbContainer.Terminate(context.Background())

	l := zap.NewExample()

	repo := &storage.Storage{mssql.NewStorage(db, l)}

	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "15.04.2021"}},
		{name: "invalid date",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "15.44.2021"}, wantErr: true},
		{name: "invalid date",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "asd"}, wantErr: true},
		{name: "success",
			fields: fields{currency: repo, l: l, externalAPI: "https://nationalbank.kz/rss/get_rates.cfm?fdate="},
			args:   args{context.TODO(), "15.04.2025"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CurrencyService{
				currency:    tt.fields.currency,
				l:           tt.fields.l,
				externalAPI: tt.fields.externalAPI,
			}
			if err := s.SaveCurrency(tt.args.ctx, tt.args.date); (err != nil) != tt.wantErr {
				t.Errorf("SaveCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(err)
		})
	}
}

func Test_makeURL(t *testing.T) {
	type args struct {
		api  string
		date string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"success", args{api: "https://nationalbank.kz/rss/get_rates.cfm?fdate=", date: "24.02.2015"}, "https://nationalbank.kz/rss/get_rates.cfm?fdate=24.02.2015"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeURL(tt.args.api, tt.args.date); got != tt.want {
				t.Errorf("makeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
