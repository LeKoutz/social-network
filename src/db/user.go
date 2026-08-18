package db

import (
	"database/sql"
	"errors"
	"forum/src/ferror"
	"forum/src/utils"
)

type UserRowType struct {
	Id            int64
	Username      string
	Hash          string
	Email         string
	OAuthProvider string
	SessionId     string
	FirstName     string
	LastName      string
	Age           int64
	Gender        string
	LastMessageTimestamp	int64
}

// Returns ONLY the `User.Hash` field for comparison against the given password
func (user *UserRowType) SelectUserPasswordByIdentifier(identifier string) error {
	err := db.QueryRow(`SELECT hash FROM users WHERE email = ? OR username = ?`, identifier, identifier).Scan(&user.Hash)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (user *UserRowType) SelectUserByIdentifier(identifier string) error {
	err := db.QueryRow(`SELECT id, email, username FROM users WHERE email = ? OR username = ?`, identifier, identifier).Scan(&user.Id, &user.Email, &user.Username)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (user *UserRowType) SelectUserBySession() error {
	err := db.QueryRow(`SELECT id, email, username FROM users WHERE session_key = ?`, user.SessionId).Scan(&user.Id, &user.Email, &user.Username)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (user *UserRowType) SelectUserByOAuthProviderAndEmail() error {
	err := db.QueryRow(`SELECT id, email, username, oauth_provider FROM users WHERE oauth_provider = ? AND email = ?`, user.OAuthProvider, user.Email).Scan(&user.Id, &user.Email, &user.Username, &user.OAuthProvider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ferror.ErrorNoRows
		}
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (user *UserRowType) SelectUserById() error {
	err := db.QueryRow(`SELECT username FROM users WHERE id = ?`, user.Id).Scan(&user.Username)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (u *UserRowType) InsertUserWithHash() error {
	stmt, err := db.Prepare("INSERT INTO users (username, first_name, last_name, gender, age, email, hash) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	res, err := stmt.Exec(u.Username, u.FirstName, u.LastName, u.Gender, u.Age, u.Email, u.Hash)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	u.Id, err = res.LastInsertId()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (u *UserRowType) InsertUserWithOAuth() error {
	stmt, err := db.Prepare("INSERT INTO users (username, email, oauth_provider) VALUES (?, ?, ?)")
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	res, err := stmt.Exec(u.Username, u.Email, u.OAuthProvider)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	u.Id, err = res.LastInsertId()
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	return nil
}

func (u *UserRowType) SelectPosts() (PostRowsType, error) {
	var posts PostRowsType
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	WHERE user_id = ?`, u.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return PostRowsType{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post PostRowType
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return PostRowsType{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return PostRowsType{}, err
		}
		post.TimestampString = utils.ConvertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func (u *UserRowType) SelectLikedPosts() (PostRowsType, error) {
	var posts PostRowsType
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	JOIN reactions r ON posts.id = r.post_id
	WHERE r.user_id = ? AND r.value = 1
	`, u.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return PostRowsType{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post PostRowType
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return PostRowsType{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
			return PostRowsType{}, err
		}
		post.TimestampString = utils.ConvertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func (u *UserRowType) UpdateUserSession(session_key string) error {
	stmt, err := db.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	_, err = stmt.Exec(session_key, u.Id)
	if err != nil {
		if utils.GlobalConfig.Debug { err = errors.Join(utils.GetFunctionName(), err) }
		return err
	}
	u.SessionId = session_key
	return nil
}
