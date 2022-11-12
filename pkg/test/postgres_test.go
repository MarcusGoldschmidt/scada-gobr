package test

import "testing"

func TestPingPostgres(t *testing.T) {
	container := SetupPostgresContainer(t)
	defer container.Close()

	err := container.Db().Ping()
	if err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}
