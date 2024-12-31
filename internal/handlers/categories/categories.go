package categories

import (
    "encoding/json"
	"fmt"
    "net/http"
	"strconv"
	"github.com/go-chi/chi/v5"

	"github.com/blobfish465/common-circle-web-forum/internal/api"
	"github.com/blobfish465/common-circle-web-forum/internal/dataaccess/categories"
    "github.com/blobfish465/common-circle-web-forum/internal/database"
	"github.com/pkg/errors"
)

const (
	ErrRetrieveDatabase           = "Failed to retrieve database in %s"
	ErrRetrieveCategories          = "Failed to retrieve Categories in %s"
	ErrEncodeView                 = "Failed to encode Categories in %s"
)

// Retrieves all categories from the database and returns them as a JSON response.
func HandleListCategories(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, "HandleListCategories"))
	}

	categoriesList, err := categories.List(db)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveCategories, "HandleListCategories"))
	}

	data, err := json.Marshal(categoriesList)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, "HandleListCategories"))
	}

	return &api.Response{
		Payload: api.Payload{
			Data: data,
		},
		Messages: []string{"Categories retrieved successfully."},
	}, nil
}


// Retrieves the category details of a category
func HandleGetCategoryByID(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
    categoryIDStr := chi.URLParam(r, "id")
    if categoryIDStr == "" {
        return nil, fmt.Errorf("missing category ID")
    }

    categoryID, err := strconv.Atoi(categoryIDStr)
    if err != nil {
        return nil, fmt.Errorf("invalid category ID: %w", err)
    }

    db, err := database.GetDB()
    if err != nil {
        return nil, errors.Wrap(err, "failed to retrieve database")
    }

    category, err := categories.GetCategoryByID(db, categoryID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve category: %w", err)
    }

    data, err := json.Marshal(category)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal category data: %w", err)
    }

    return &api.Response{
        Payload: api.Payload{
            Data: data,
        },
        Messages: []string{"Category retrieved successfully"},
    }, nil
}


