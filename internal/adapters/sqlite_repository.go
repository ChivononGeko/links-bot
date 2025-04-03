package adapters

import (
	"certificate/internal/domain"
	"certificate/internal/ports"
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (ports.RegistrationRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(
		"file://migrations/",
		"sqlite://"+dbPath,
	)
	if err != nil {
		slog.Error("Ошибка при создании мигратора: %v", "error", err)
		return nil, err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("Ошибка при применении миграций: %v", "error", err)
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

// Создание записи с токеном
func (r *SQLiteRepository) Create(token string) error {
	_, err := r.db.Exec("INSERT INTO registrations (token) VALUES (?)", token)
	return err
}

// Получение токена по значению
func (r *SQLiteRepository) GetByToken(token string) (*domain.Registration, error) {
	row := r.db.QueryRow("SELECT id, token, used FROM registrations WHERE token = ?", token)
	reg := &domain.Registration{}
	err := row.Scan(&reg.ID, &reg.Token, &reg.Used)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

// Отметка токена как использованного и запись кто это сделал
func (r *SQLiteRepository) MarkTokenUsed(token, name, phone string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Обновляем статус токена на "использованный"
	_, err = tx.Exec("UPDATE registrations SET used = TRUE WHERE token = ?", token)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Записываем информацию о пользователе, который использовал токен
	_, err = tx.Exec("INSERT INTO token_usage (token, username, phone) VALUES (?, ?, ?)", token, name, phone)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Фиксируем изменения в БД
	return tx.Commit()
}

// Получить данные пользователя, который использовал токен
func (r *SQLiteRepository) GetTokenUsage(token string) (*domain.TokenUsage, error) {
	row := r.db.QueryRow("SELECT id, token, username, phone FROM token_usage WHERE token = ?", token)
	usage := &domain.TokenUsage{}
	err := row.Scan(&usage.ID, &usage.Token, &usage.Username, &usage.Phone)
	if err != nil {
		return nil, err
	}
	return usage, nil
}

// Получить список использованных токенов
func (r *SQLiteRepository) GetUsedTokens() ([]domain.Registration, error) {
	rows, err := r.db.Query("SELECT id, token, used FROM registrations WHERE used = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(&reg.ID, &reg.Token, &reg.Used); err != nil {
			return nil, err
		}
		tokens = append(tokens, reg)
	}

	return tokens, nil
}

// Получить список неиспользованных токенов
func (r *SQLiteRepository) GetUnusedTokens() ([]domain.Registration, error) {
	rows, err := r.db.Query("SELECT id, token, used FROM registrations WHERE used = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(&reg.ID, &reg.Token, &reg.Used); err != nil {
			return nil, err
		}
		tokens = append(tokens, reg)
	}

	return tokens, nil
}
