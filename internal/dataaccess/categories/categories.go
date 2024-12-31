package categories

import (
	"github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/blobfish465/common-circle-web-forum/internal/models"
)

func List(db *database.Database) ([]models.Category, error) {
    rows, err := db.DB.Query("SELECT id, name FROM categories")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var category models.Category
        err := rows.Scan(&category.ID, &category.Name)
        if err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func GetCategoryByID(db *database.Database, id int) (*models.Category, error) {
    query := `SELECT id, name FROM categories WHERE id = $1`

    row := db.DB.QueryRow(query, id)

    var category models.Category
    err := row.Scan(&category.ID, &category.Name)
    if err != nil {
        return nil, err
    }

    return &category, nil
}
