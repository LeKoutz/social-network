package forum

import "os"

func MockDbData() {
	if err := Init(); err != nil {
		(&Error{}).Consume(err).LogError()
		os.Exit(1)
	}
	InitCategories()
	InitUsers()
	InitPosts()
	InitComments()
	// InitPostsCategories()
	InitLikes()
	InitDislikes()
}

func InitUsers() {
	err := registerUserOnDB(User{
		Email:    "user@example.mock",
		Username: "justAUser",
		// password : lol
		Hash: "$2a$10$dHn/OlFTaoVxs7ClR23Q7e7qDHUvaJRJI1kdTP.Rfga31UzJ8.BMG",
	})
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
}

func InitPosts() {
	_, err := addPost(Post{
		Title:  "lol",
		Body:   "mpla",
		UserId: 1,
		Category: Category{
			Id: 1,
		},
	})
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
	_, err = addPost(Post{
		Title:  "kek",
		Body:   "alpm",
		UserId: 1,
		Category: Category{
			Id: 2,
		},
	})
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
}

func InitComments() {
	_, err := addComment(Comment{
		PostId: 1,
		UserId: 1,
		Body: "lol",
	})
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
}
// func InitPostsCategories() {
//
// 	if err != nil {
// 		(&Error{}).Consume(err).LogError()
// 	}
// }

func InitLikes() {
	err := addLikeToPost(1, 1)
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
}

func InitDislikes() {
	err := addDislikeToPost(1, 2)
	if err != nil {
		(&Error{}).Consume(err).LogError()
	}
}
