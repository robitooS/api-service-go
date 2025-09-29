package db

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"time"
)

func RunMigrations(db *sql.DB) error {
	// Cria a tabela para guardar a versão da ultima migration que foi aplicada
	err := createSchemaTable(db)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela de migrations: %v", err)
	}

	err = applyPendingMigrations(db)
	if err != nil {
		return fmt.Errorf("erro ao adicionar migrations pendentes: %v", err)
	}

	return nil
}

func applyPendingMigrations(db *sql.DB) error {
	// Verificar o que falta de acordo com o que tem dentro da pasta e o que retorna do map
	applieds, err := getVersionsFromDB(db)
	if err != nil {
		return err
	}
	files, err := os.ReadDir("./migrations")
	if err != nil {
		return err
	}
	sort.Slice(files, func (i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if applieds[name] {
			fmt.Println("Migration já aplicada:", name)
			continue
		}
		sqlBytes, err := os.ReadFile("./migrations/" + name)
		if err != nil {
			return fmt.Errorf("não foi possível ler a migration %s - %v", name, err)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("não foi possível aplicar a migration %s no banco - %v", name, err)
		}

		_, err = db.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)", name, time.Now())
		if err != nil {
			return fmt.Errorf("não foi possível inserir o novo registro na tabela de migrations - %v", err)
		}
	}
	return nil
}

// Função para verificar a última versão que foi aplicada e se possui uma nova
func getVersionsFromDB(db *sql.DB) (map[string]bool, error) {
	versions := make(map[string]bool)

	rows, err := db.Query(`
	SELECT version FROM schema_migrations
	`)
	if err != nil {
		return nil, fmt.Errorf("erro ao procurar última versão da migration aplicada: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil  {
			return nil, fmt.Errorf("erro ao verificar a ultima versão aplicada do migration")
		}
		versions[version] = true
	}
	return versions, nil
}

func createSchemaTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `)
    return err
}