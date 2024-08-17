package postgres_repo

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func Migrate(db *sql.DB) error {
	migrationFiles, err := filepath.Glob("./repo/postgres_repo/migrations/*.sql")
	if err != nil {
		return err
	}

	for _, file := range migrationFiles {
		if err := applyMigragtion(db, file); err != nil {
			return fmt.Errorf("%s: :%v", file, err)
		}
	}

	return nil
}

func applyMigragtion(db *sql.DB, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return err
	}

	return nil
}
