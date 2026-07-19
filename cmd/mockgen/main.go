package main

import (
	"fmt"
	"os"
	"time"

	"forum/src/models"
	"forum/src/utils"
)

const demoPassword = "Demo1234!"

type demoUser struct {
	Username  string
	FirstName string
	LastName  string
	Age       int64
	Gender    string
	PostCount int
}

var demoUsers = []demoUser{
	{Username: "demo_alex", FirstName: "Alex", LastName: "Morgan", Age: 24, Gender: "other", PostCount: 3},
	{Username: "demo_bella", FirstName: "Bella", LastName: "Stone", Age: 27, Gender: "female", PostCount: 2},
	{Username: "demo_charlie", FirstName: "Charlie", LastName: "Young", Age: 31, Gender: "male", PostCount: 3},
	{Username: "demo_diana", FirstName: "Diana", LastName: "Green", Age: 22, Gender: "female", PostCount: 2},
	{Username: "demo_ethan", FirstName: "Ethan", LastName: "Cole", Age: 29, Gender: "male", PostCount: 3},
	{Username: "demo_freya", FirstName: "Freya", LastName: "Hall", Age: 26, Gender: "female", PostCount: 2},
	{Username: "demo_george", FirstName: "George", LastName: "King", Age: 35, Gender: "male", PostCount: 3},
	{Username: "demo_helen", FirstName: "Helen", LastName: "Price", Age: 28, Gender: "female", PostCount: 2},
	{Username: "demo_iris", FirstName: "Iris", LastName: "Baker", Age: 23, Gender: "other", PostCount: 3},
	{Username: "demo_jason", FirstName: "Jason", LastName: "Scott", Age: 33, Gender: "male", PostCount: 2},
	{Username: "demo_kira", FirstName: "Kira", LastName: "Adams", Age: 25, Gender: "female", PostCount: 3},
	{Username: "demo_leon", FirstName: "Leon", LastName: "Evans", Age: 30, Gender: "male", PostCount: 2},
	{Username: "demo_maria", FirstName: "Maria", LastName: "Lopez", Age: 21, Gender: "female", PostCount: 3},
	{Username: "demo_nikos", FirstName: "Nikos", LastName: "Pappas", Age: 32, Gender: "male", PostCount: 2},
	{Username: "demo_zoe", FirstName: "Zoe", LastName: "Reed", Age: 27, Gender: "other", PostCount: 3},
}

type demoConversation struct {
	OtherUsername  string
	MessageCount   int
	LastMessageAgo time.Duration
}

var demoConversations = []demoConversation{
	{OtherUsername: "demo_diana", MessageCount: 4, LastMessageAgo: 5 * time.Minute},
	{OtherUsername: "demo_bella", MessageCount: 5, LastMessageAgo: 20 * time.Minute},
	{OtherUsername: "demo_ethan", MessageCount: 12, LastMessageAgo: time.Hour},
	{OtherUsername: "demo_charlie", MessageCount: 4, LastMessageAgo: 3 * time.Hour},
}

// mockGen generates mock data for the database.
func mockGen(dbPath string) error {
	if err := models.InitDB(dbPath); err != nil {
		return err
	}

	categories := []models.Category{
		{Name: "General", Description: "General discussions"},
		{Name: "Tech", Description: "Technology related discussions"},
		{Name: "Random", Description: "Random topics"},
	}
	for _, category := range categories {
		if !category.DoesCategoryExist() {
			if err := models.AddCategory(category); err != nil {
				return err
			}
		}
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		return err
	}
	if len(categories) == 0 {
		return fmt.Errorf("cannot create demo posts without categories")
	}
	if err := seedDemoUsers(); err != nil {
		return err
	}
	usersByUsername, err := getUsersByUsername()
	if err != nil {
		return err
	}
	if err := seedDemoPosts(usersByUsername, categories); err != nil {
		return err
	}
	return seedDemoMessages(usersByUsername)
}

func seedDemoUsers() error {
	hash, err := utils.HashPassword(demoPassword)
	if err != nil {
		return err
	}
	for _, demo := range demoUsers {
		if !models.IsUniqueUsername(demo.Username) {
			continue
		}
		user := models.User{
			Username:  demo.Username,
			FirstName: demo.FirstName,
			LastName:  demo.LastName,
			Age:       demo.Age,
			Gender:    demo.Gender,
			Email:     demo.Username + "@example.test",
			Hash:      hash,
		}
		if err := user.Add(); err != nil {
			return fmt.Errorf("add demo user %s: %w", demo.Username, err)
		}
	}
	return nil
}

func getUsersByUsername() (map[string]models.User, error) {
	users, err := models.GetAllUsers()
	if err != nil {
		return nil, err
	}
	usersByUsername := make(map[string]models.User, len(users))
	for _, user := range users {
		usersByUsername[user.Username] = user
	}
	return usersByUsername, nil
}

func seedDemoPosts(usersByUsername map[string]models.User, categories models.Categories) error {
	for userIndex, demo := range demoUsers {
		user, ok := usersByUsername[demo.Username]
		if !ok {
			return fmt.Errorf("demo user %s was not found", demo.Username)
		}
		existingPosts, err := user.GetPosts()
		if err != nil {
			return fmt.Errorf("get posts for %s: %w", demo.Username, err)
		}
		for postIndex := len(existingPosts); postIndex < demo.PostCount; postIndex++ {
			category := categories[(userIndex+postIndex)%len(categories)]
			post := models.Post{
				Title:      fmt.Sprintf("Demo post %d by %s", postIndex+1, demo.Username),
				Body:       fmt.Sprintf("This is demo content from %s for testing the forum feed and user ordering.", demo.Username),
				UserId:     user.Id,
				Categories: models.Categories{category},
			}
			if _, err := post.Add(); err != nil {
				return fmt.Errorf("add post for %s: %w", demo.Username, err)
			}
		}
	}
	return nil
}

func seedDemoMessages(usersByUsername map[string]models.User) error {
	alex, ok := usersByUsername["demo_alex"]
	if !ok {
		return fmt.Errorf("demo user demo_alex was not found")
	}
	now := time.Now()
	for _, conversation := range demoConversations {
		other, ok := usersByUsername[conversation.OtherUsername]
		if !ok {
			return fmt.Errorf("demo user %s was not found", conversation.OtherUsername)
		}
		existing, err := models.GetChatHistory(alex.Id, other.Id, 0)
		if err != nil {
			return fmt.Errorf("get chat history with %s: %w", conversation.OtherUsername, err)
		}
		if len(existing) > 0 {
			continue
		}
		lastMessageAt := now.Add(-conversation.LastMessageAgo)
		firstMessageAt := lastMessageAt.Add(-time.Duration(conversation.MessageCount-1) * time.Minute)
		for messageIndex := 0; messageIndex < conversation.MessageCount; messageIndex++ {
			sender := alex
			recipient := other
			if messageIndex%2 == 0 {
				sender, recipient = other, alex
			}
			message := models.ChatMessage{
				SenderId:        sender.Id,
				RecipientId:     recipient.Id,
				Body:            fmt.Sprintf("Demo message %d between %s and demo_alex", messageIndex+1, conversation.OtherUsername),
				TimestampString: fmt.Sprintf("%d", firstMessageAt.Add(time.Duration(messageIndex)*time.Minute).Unix()),
			}
			if _, err := message.Add(); err != nil {
				return fmt.Errorf("add chat message with %s: %w", conversation.OtherUsername, err)
			}
		}
	}
	return nil
}

func main() {
	dbPath := "./db.db"
	if len(os.Args) == 2 {
		dbPath = os.Args[1]
	}
	if err := mockGen(dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating demo data: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Demo data ready in %s\n", dbPath)
	fmt.Printf("Demo password for all demo_* users: %s\n", demoPassword)
}
