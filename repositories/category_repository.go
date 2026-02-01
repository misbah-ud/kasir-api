package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name FROM categories ORDER BY id"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (repo *CategoryRepository) Create(category *models.Category) error {
	query := "INSERT INTO categories (name) VALUES ($1) RETURNING id"
	return repo.db.QueryRow(query, category.Name).Scan(&category.ID)
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name FROM categories WHERE id = $1"

	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("category tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE categories SET name = $1 WHERE id = $2"
	result, err := repo.db.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}

	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("category tidak ditemukan")
	}

	return nil
}
