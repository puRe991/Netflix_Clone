package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "" +
		"# comment\n" +
		"\n" +
		`DATABASE_URL="postgresql://streamflix:streamflix@localhost:5432/streamflix?sslmode=disable"` + "\n" +
		"JWT_SECRET=super-secret\n" +
		"export APP_ENV=development\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "APP_ENV", "PORT"} {
		os.Unsetenv(key)
	}

	loadDotEnv(envPath)
	defer func() {
		for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "APP_ENV"} {
			os.Unsetenv(key)
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := cfg.DatabaseURL; got != "postgresql://streamflix:streamflix@localhost:5432/streamflix?sslmode=disable" {
		t.Errorf("DatabaseURL = %q, want unquoted value from .env", got)
	}
	if got := cfg.JWTSecret; got != "super-secret" {
		t.Errorf("JWTSecret = %q, want %q", got, "super-secret")
	}
	if got := cfg.AppEnv; got != "development" {
		t.Errorf("AppEnv = %q, want %q (export prefix should be stripped)", got, "development")
	}
}

func TestLoadDotEnvDoesNotOverrideRealEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte(`DATABASE_URL="from-file"`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	os.Setenv("DATABASE_URL", "from-real-env")
	defer os.Unsetenv("DATABASE_URL")

	loadDotEnv(envPath)

	if got := os.Getenv("DATABASE_URL"); got != "from-real-env" {
		t.Errorf("DATABASE_URL = %q, want real environment variable to take precedence", got)
	}
}

func TestLoadDotEnvMissingFileIsNotAnError(t *testing.T) {
	loadDotEnv(filepath.Join(t.TempDir(), "does-not-exist.env"))
}
