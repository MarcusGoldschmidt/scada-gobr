package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"testing"
	"time"
)

func randomPort() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer l.Close()
	_, port, err := net.SplitHostPort(l.Addr().String())
	return port, err
}

type SqlContainer struct {
	sqlDB     *sql.DB
	container testcontainers.Container
	gormDb    *gorm.DB
}

func (d *SqlContainer) GormDb() *gorm.DB {
	return d.gormDb
}

func (d *SqlContainer) Db() *sql.DB {
	return d.sqlDB
}

func (d *SqlContainer) Close() error {
	if d.container != nil {
		return d.container.Terminate(context.Background())
	}

	return nil
}

func SetupPostgresContainer(t *testing.T) *SqlContainer {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "timescale/timescaledb:latest-pg14",
		AutoRemove:   true,
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Fatalf("failed to start container: %v", err)
		return nil
	}

	port, err := container.MappedPort(ctx, "5432/tcp")

	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
		return nil
	}

	uri := fmt.Sprintf("host=%s user=postgres database=postgres password=postgres port=%d", "localhost", port.Int())

	t.Log(uri)

	sqlDB, err := sql.Open("pgx", uri)

	if err != nil {
		t.Fatalf("failed to open database: %v", err)
		return nil
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
		return nil
	}

	return &SqlContainer{gormDb: gormDb, sqlDB: sqlDB, container: container}
}
