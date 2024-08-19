package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL!")

	return &MySQLStorage{db: db}
}

func (s *MySQLStorage) InitializeDatabase() (*sql.DB, error) {
	// initialize the tables
	if err := s.createUsersTable(); err != nil {
		return nil, err
	}

	if err := s.createProjectsTable(); err != nil {
		return nil, err
	}

	if err := s.createTasksTable(); err != nil {
		return nil, err
	}
	if err := s.applySchemaChanges(); err != nil {
		return nil, err
	}

	if err := s.createPasswordResetTokensTable(); err != nil {
		return nil, err
	}

	if err := s.createTokenBlacklistTable(); err != nil {
		return nil, err
	}
	return s.db, nil
}

func (s *MySQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			email VARCHAR(255) NOT NULL,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			phone VARCHAR(20),
			address VARCHAR(255),

			PRIMARY KEY (id),
			UNIQUE KEY (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createProjectsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createTasksTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			status ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE') NOT NULL DEFAULT 'TODO',
			projectId INT UNSIGNED NOT NULL,
			AssignedToID INT UNSIGNED NOT NULL,
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id),
			FOREIGN KEY (AssignedToID) REFERENCES users(id),
			FOREIGN KEY (projectId) REFERENCES projects(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) createPasswordResetTokensTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS password_reset_tokens (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		user_id INT UNSIGNED NOT NULL,
		token VARCHAR(255) NOT NULL,
		expiration TIMESTAMP NOT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE KEY (token)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) applySchemaChanges() error {
	// Check if the 'phone' column exists
	columnExists, err := s.columnExists("users", "phone")
	if err != nil {
		return err
	}
	if !columnExists {
		_, err := s.db.Exec("ALTER TABLE users ADD COLUMN phone VARCHAR(20)")
		if err != nil {
			return err
		}
	}

	// Check if the 'address' column exists
	columnExists, err = s.columnExists("users", "address")
	if err != nil {
		return err
	}
	if !columnExists {
		_, err := s.db.Exec("ALTER TABLE users ADD COLUMN address VARCHAR(255)")
		if err != nil {
			return err
		}
	}

	// Apply similar changes to other tables if needed

	return nil
}

func (s *MySQLStorage) columnExists(tableName, columnName string) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM INFORMATION_SCHEMA.COLUMNS 
        WHERE TABLE_NAME = ? 
          AND COLUMN_NAME = ?
    `
	var count int
	err := s.db.QueryRow(query, tableName, columnName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *MySQLStorage) createTokenBlacklistTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE  IF NOT EXISTS token_blacklist (
		token VARCHAR(255) NOT NULL,
		createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (token)
	);
	`)

	return err
}
