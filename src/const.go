package forum

import (
	"forum/src/models"
	"forum/src/utils"

	"golang.org/x/crypto/bcrypt"
)

const (
	WebsiteName = "Forum"
)

// HashPassword returns the bcrypt hash of a password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// MockGen generates mock data for the database
func MockGen(dbPath string) error {
	if err := models.InitDB(dbPath); err != nil {
		return err
	}

	// Create some mock categories
	categories := []models.Category{
		{Name: "General", Description: "General discussions"},
		{Name: "Tech", Description: "Technology related discussions"},
		{Name: "Random", Description: "Random topics"},
	}

	for _, cat := range categories {
		if !cat.DoesCategoryExist() {
			if err := models.AddCategory(cat); err != nil {
				utils.LogDebug(err)
			}
		}
	}

	return nil
}
