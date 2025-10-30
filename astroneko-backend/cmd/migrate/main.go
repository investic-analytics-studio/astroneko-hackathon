package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"astroneko-backend/configs"
	"astroneko-backend/pkg/databases/gorm"
)

func main() {
	// Initialize configuration
	configs.InitViper("./configs")
	config := configs.GetViper()

	// Connect to database
	var db *gorm.DB
	var err error

	if config.App.Env == "local" {
		db, err = gorm.ConnectToPostgreSQL(
			config.Postgres.Host,
			config.Postgres.Port,
			config.Postgres.Username,
			config.Postgres.Password,
			config.Postgres.DbName,
			config.Postgres.SSLMode,
		)
	} else {
		db, err = gorm.ConnectToCloudSQL(
			config.Postgres.InstanceConnectionName,
			config.Postgres.Username,
			config.Postgres.Password,
			config.Postgres.DbName,
		)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer gorm.DisconnectPostgres(db.Postgres)

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("All migrations completed successfully!")
}

func createMigrationsTable(db *gorm.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
	`
	return db.Postgres.Exec(createTableSQL).Error
}

func runMigrations(db *gorm.DB) error {
	migrationsDir := "./migrations"

	// Get list of migration files
	migrationFiles, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %v", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Run pending migrations
	for _, file := range migrationFiles {
		version := strings.TrimSuffix(file, ".sql")

		if _, applied := appliedMigrations[version]; applied {
			fmt.Printf("Migration %s already applied, skipping...\n", version)
			continue
		}

		fmt.Printf("Running migration: %s\n", version)

		if err := runMigration(db, filepath.Join(migrationsDir, file)); err != nil {
			return fmt.Errorf("failed to run migration %s: %v", version, err)
		}

		// Record migration as applied
		if err := recordMigration(db, version); err != nil {
			return fmt.Errorf("failed to record migration %s: %v", version, err)
		}

		fmt.Printf("Migration %s completed successfully\n", version)
	}

	return nil
}

func getMigrationFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, d.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure correct order
	sort.Strings(files)
	return files, nil
}

func getAppliedMigrations(db *gorm.DB) (map[string]bool, error) {
	var versions []string
	err := db.Postgres.Raw("SELECT version FROM schema_migrations").Scan(&versions).Error
	if err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, version := range versions {
		applied[version] = true
	}

	return applied, nil
}

func runMigration(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	sql := string(content)
	if strings.TrimSpace(sql) == "" {
		return nil
	}

	return db.Postgres.Exec(sql).Error
}

func recordMigration(db *gorm.DB, version string) error {
	return db.Postgres.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version).Error
}
