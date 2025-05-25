package db

import (
	"database/sql"
	"fmt"
	"userservice/internal/models"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(host, port, user, password, dbname string) (*PostgresDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) CreateUser(user *models.User) error {
	query := `
        INSERT INTO users (name, email, updated_at)
        VALUES ($1, $2, CURRENT_TIMESTAMP)
        RETURNING id, updated_at`

	return p.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID, &user.UpdatedAt)
}

func (p *PostgresDB) UpdateUser(user *models.User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
        RETURNING updated_at`

	return p.db.QueryRow(query, user.Name, user.Email, user.ID).Scan(&user.UpdatedAt)
}

func (p *PostgresDB) GetUser(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, updated_at FROM users WHERE id = $1`

	err := p.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PostgresDB) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, name, email, updated_at FROM users ORDER BY id`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
