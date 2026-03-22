package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"social-network/models"
	"social-network/pkg/helpers"
	"time"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	Db *sql.DB
}


func InitDB() (*sql.DB, error) {
	// prepare the driver
	db, er := sql.Open("sqlite3", "./db/rtf.db")
	if er != nil {
		return nil, er
	}

	// test connection
	er = db.Ping()
	if er != nil {
		return nil, er
	}

	//
	sql, er := os.ReadFile("./db/init.sql")
	if er != nil {
		return nil, er
	}
	_, er = db.Exec(string(sql))
	if er != nil {
		return nil, er
	}

	return db, nil
}


// try to insert the user into the data base, any invalid input will return an error with a specific message
func (r *Repo) InsertUserDB(user models.User) error {
	// check the user infos if correct
	err := helpers.AreUserInfosCorret(user)
	if err != nil {
		return err
	}

	// check the user existance in DB
	var exist int
	err = r.Db.QueryRow("SELECT 1 FROM users WHERE nickname=? OR email=?", user.Nickname, user.Email).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		return errors.New("SERVER ERROR")
	}
	if exist > 0 {
		return errors.New("user alrady exist")
	}

	//
	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userUUID, err := uuid.NewV4()
	if err != nil {
		return errors.New("SERVER ERROR")
	}
	user.ID = userUUID.String()

	//
	_, err = r.Db.Exec(
		"INSERT INTO users(id, nickname, birthday, gender, firstname, lastname, email, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Nickname, user.Birthday, user.Gender, user.Firstname, user.Lastname, user.Email, hashed,
	)
	if err != nil {
		return errors.New("SERVER ERROR")
	}

	return nil
}

// check existance of the user in the DB
func (r *Repo) IsUserExist(user *models.User) (string, error) {
	var id string
	var hashedPassword string

	if len(user.Email) > 60 || len(user.Nickname) > 60 || len(user.Password) > 60 {
		return "", errors.New("user not exist")
	}

	err := r.Db.QueryRow("SELECT id, password FROM users WHERE nickname=? OR email=?", user.Nickname, user.Nickname).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not exist")
		}
		return "", errors.New("SERVER ERROR")
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)) != nil {
		return "", errors.New("password invalid")
	}

	return id, nil
}

// set new session in case of user login
func (r *Repo) SetUserSession(w http.ResponseWriter, userID string) ([]interface{}, error) {
	sessionUUID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("SERVER ERROR")
	}
	sessionId := sessionUUID.String()
	now := time.Now()
	timeExpired := now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	timeNow := now.Format("2006-01-02 15:04:05")
	e := now.Add(24 * time.Hour)

	_, err = r.Db.Exec("UPDATE users SET session_id=?, session_created_at=?, session_expired_at=? WHERE id=?", sessionId, timeNow, timeExpired, userID)
	if err != nil {
		return nil, errors.New("SERVER ERROR")
	}
	return []interface{}{sessionId, e}, nil
}

// delete the session from the DB in case of logout
func (r *Repo) DisconnectUser(userID string) error {
	_, er := r.Db.Exec("UPDATE users SET session_id=NULL, session_created_at=NULL, session_expired_at=NULL WHERE id=?", userID)
	if er != nil {
		return errors.New("SERVER ERROR")
	}
	return nil
}

// this check session sended by the browser if it is included in the DB
func (r *Repo) CheckSessionExistance(req *http.Request) (models.User, error) {
	var user models.User

	// check in the browser
	cookie, err := req.Cookie("session_id")
	if err != nil || cookie == nil || cookie.Value == "" {
		return user, fmt.Errorf("Error-session")
	}

	// check in DB
	err = r.Db.QueryRow("SELECT id, nickname, birthday, gender, firstname, lastname, email, session_expired_at FROM users WHERE session_id = ?", cookie.Value).
		Scan(&user.ID, &user.Nickname, &user.Birthday, &user.Gender, &user.Firstname, &user.Lastname, &user.Email, &user.SessionExpired)
	if err != nil {
		return user, err
	}

	// check if the session already expired
	if user.SessionExpired != "" {
		sessionExpiredTime, err := time.Parse("2006-01-02 15:04:05", user.SessionExpired)
		if err != nil {
			return user, err
		}
		if time.Now().After(sessionExpiredTime) {
			return user, errors.New("session expired")
		}
	}

	return user, nil
}

// get user infos from DB
func (r *Repo) GetUserInfos(userID string) (models.User, error) {
	var user models.User

	err := r.Db.QueryRow(`SELECT id, nickname, birthday, gender, firstname, lastname, email 
		FROM users WHERE id=?`, userID).
		Scan(&user.ID, &user.Nickname, &user.Birthday, &user.Gender, &user.Firstname, &user.Lastname, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}
