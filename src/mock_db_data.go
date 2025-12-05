package forum

import "log"

func MockDbData() {
	if err := Init(); err != nil {
		log.Fatal(err.Error())
	}
	InitCategories()
	InitUsers()
	InitPosts()
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
		log.Print(err.Error())
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
		log.Print(err.Error())
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
		log.Print(err.Error())
	}
}

// func InitPostsCategories() {
//
// 	if err != nil {
// 		log.Print(err.Error())
// 	}
// }

func InitLikes() {
	err := addLikeToPost(1, 1)
	if err != nil {
		log.Print(err.Error())
	}
}

func InitDislikes() {
	err := addDislikeToPost(1, 2)
	if err != nil {
		log.Print(err.Error())
	}
}
