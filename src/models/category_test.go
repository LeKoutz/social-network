package models

import (
	"forum/src/db"
	"forum/src/ferror"
	"forum/src/utils"
	"strings"
	"testing"
)

func setupTestCategoryDB(t *testing.T) {
	t.Helper()
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}
}

func TestCategoryValidateCategory(t *testing.T) {
	tests := []struct {
		name    string
		catName string
		wantErr error
	}{
		{"valid name", "general", nil},
		{"single char", "a", nil},
		{"127 chars", strings.Repeat("a", 127), nil},
		{"empty name", "", ferror.ErrorCategoryNameEmpty},
		{"128 chars", strings.Repeat("a", 128), ferror.ErrorCategoryNameTooLong},
		{"very long", strings.Repeat("a", 256), ferror.ErrorCategoryNameTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryType{}
			c.Name = tt.catName
			err := c.ValidateCategory()
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateCategory() expected error %v, got nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCategory() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCategoryAdd(t *testing.T) {
	setupTestCategoryDB(t)

	c := &CategoryType{}
	c.Name = "testing"
	c.Description = "Testing category"
	err := c.Add()
	if err != nil {
		t.Fatalf("Add() error: %v", err)
	}
	if c.Id == 0 {
		t.Error("Add() did not set category ID")
	}

	c2 := &CategoryType{}
	c2.Name = "testing"
	c2.Description = "Duplicate category"
	err = c2.Add()
	if err == nil {
		t.Error("Add() should fail for duplicate name")
	}
}

func TestCategoryAddEmptyName(t *testing.T) {
	setupTestCategoryDB(t)

	c := &CategoryType{}
	c.Name = ""
	c.Description = "No name"
	err := c.Add()
	if err == nil {
		t.Error("Add() should fail with empty name")
	}
}

func TestCategoryAddTooLongName(t *testing.T) {
	setupTestCategoryDB(t)

	c := &CategoryType{}
	c.Name = strings.Repeat("a", 128)
	c.Description = "Too long"
	err := c.Add()
	if err == nil {
		t.Error("Add() should fail with name >= 128 chars")
	}
}

func TestCategoryIsEmpty(t *testing.T) {
	var nilCat *CategoryType
	if !nilCat.IsEmpty() {
		t.Error("IsEmpty() should return true for nil")
	}

	emptyCat := &CategoryType{}
	if !emptyCat.IsEmpty() {
		t.Error("IsEmpty() should return true for empty CategoryType")
	}

	nonEmptyCat := &CategoryType{}
	nonEmptyCat.Id = 1
	nonEmptyCat.Name = "test"
	if nonEmptyCat.IsEmpty() {
		t.Error("IsEmpty() should return false for non-empty CategoryType")
	}

	partialCat := &CategoryType{}
	partialCat.Name = "only name"
	if partialCat.IsEmpty() {
		t.Error("IsEmpty() should return false for partially filled CategoryType")
	}
}

func TestCategorySelectCategoryById(t *testing.T) {
	setupTestCategoryDB(t)

	c := &CategoryType{}
	c.Name = "selecttest"
	c.Description = "Select test category"
	c.InsertCategory()

	c2 := &CategoryType{}
	c2.Id = c.Id
	err := c2.SelectCategoryById()
	if err != nil {
		t.Fatalf("SelectCategoryById() error: %v", err)
	}
	if c2.Name != "selecttest" {
		t.Errorf("SelectCategoryById() Name = %q, want %q", c2.Name, "selecttest")
	}
	if c2.Description != "Select test category" {
		t.Errorf("SelectCategoryById() Description = %q, want %q", c2.Description, "Select test category")
	}
}

func TestCategoryGetAll(t *testing.T) {
	setupTestCategoryDB(t)

	c1 := &CategoryType{}
	c1.Name = "cat1"
	c1.Description = "Category 1"
	c1.InsertCategory()

	c2 := &CategoryType{}
	c2.Name = "cat2"
	c2.Description = "Category 2"
	c2.InsertCategory()

	var categories CategoriesType
	err := categories.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("GetAll() returned %d categories, want 2", len(categories))
	}
}

func TestCategoryGetAllEmpty(t *testing.T) {
	setupTestCategoryDB(t)

	var categories CategoriesType
	err := categories.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}
	if len(categories) != 0 {
		t.Errorf("GetAll() returned %d categories, want 0", len(categories))
	}
}

func TestCategoryGetCategoriesByPostId(t *testing.T) {
	setupTestCategoryDB(t)

	user := UserType{}
	user.Username = "catpostuser"
	user.Email = "catpost@test.com"
	user.Hash, _ = utils.HashPassword("pass123456!")
	user.Add()

	cat1 := &CategoryType{}
	cat1.Name = "pcat1"
	cat1.Description = "PCategory 1"
	cat1.InsertCategory()

	cat2 := &CategoryType{}
	cat2.Name = "pcat2"
	cat2.Description = "PCategory 2"
	cat2.InsertCategory()

	post := &PostType{}
	post.Title = "Cat Post"
	post.Body = "Cat Body"
	post.UserId = user.Id
	post.Categories = CategoriesType{
		*cat1,
	}
	post.Add()

	db.InsertPostCategory(db.PostCategoryRow{PostId: post.Id, CategoryId: cat2.Id})

	categories, err := db.SelectCategoriesByPostId(post.Id)
	if err != nil {
		t.Fatalf("GetCategoriesByPostId() error: %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("GetCategoriesByPostId() returned %d categories, want 2", len(categories))
	}
}
