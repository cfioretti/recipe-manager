package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func init() {
	InjectTestProperties()
}

type TestDatabase struct {
	Container testcontainers.Container
	DB        *sql.DB
	Port      string
}

func SetupTestDb(t *testing.T) (*TestDatabase, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"0/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "test",
			"MYSQL_DATABASE":      "test_db",
		},
		WaitingFor: wait.ForSQL("0/tcp", "mysql", func(host string, port nat.Port) string {
			return fmt.Sprintf("root:test@tcp(%s:%s)/%s", host, port.Port(), "test_db")
		}),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	port, err := container.MappedPort(ctx, "0/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("root:test@tcp(localhost:%s)/test_db", port.Port())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return &TestDatabase{
		Container: container,
		DB:        db,
		Port:      port.Port(),
	}, nil
}

func runMigrations(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(filename), "../../migrations")

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (td *TestDatabase) Cleanup(ctx context.Context) error {
	td.DB.Close()
	return td.Container.Terminate(ctx)
}

func InjectTestProperties() {
	viper.SetConfigName("props-test")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../../configs/")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read test properties config: %w", err))
	}
}
