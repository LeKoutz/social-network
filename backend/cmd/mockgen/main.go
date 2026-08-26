package main

import (
	"errors"
	"os"

	"forum/src/db"
	"forum/src/models"
	"forum/src/utils"
)

// mockGen generates mock data for the database
func mockGen(dbPath string) error {
	if err := db.InitDB(dbPath); err != nil {
		return errors.Join(utils.GetFunctionName(), err)
	}

	// Create some mock categories
	var general models.CategoryType
	general.Name = "General"
	general.Description = "General discussions"
	var tech models.CategoryType
	tech.Name = "Tech"
	tech.Description = "Technology related discussions"
	var random models.CategoryType
	random.Name = "Random"
	random.Description = "Random topics"

	var categories models.CategoriesType = models.CategoriesType{general, tech, random}

	for _, cat := range categories {
		if err := cat.Add(); err != nil {
			utils.LogDebug(err)
		}
	}

	return nil
}

func main() {
	db_path := "./db.db"
	if len(os.Args) == 2 {
		db_path = os.Args[1]
	}
	mockGen(db_path)
}
