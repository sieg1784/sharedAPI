package data

import (
	"log"
	"time"
)

type UserReadingBook struct {
	UserId      int       `json:"userId"`
	BookId      string    `json:"bookId"`
	ReadingDate time.Time `json:"readingDate"`
	ReadingTime int       `json:"readingTime"`
}

func InsertUserReadingBook(userReadingBook UserReadingBook) (bool, error) {
	stmt := "INSERT INTO user_reading_book (user_id, book_id, reading_date, reading_time)  " +
		" VALUES ($1, $2, $3, $4) "
	_, err := Db.Exec(stmt, userReadingBook.UserId, userReadingBook.BookId, userReadingBook.ReadingDate, userReadingBook.ReadingTime)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func UpdateUserReadingBook(userReadingBook UserReadingBook) (bool, error) {
	stmt := "UPDATE user_reading_book SET reading_date = $1, reading_time = $2 WHERE user_id = $3 AND book_id = $4 "
	_, err := Db.Exec(stmt, userReadingBook.ReadingDate, userReadingBook.ReadingTime, userReadingBook.UserId, userReadingBook.BookId)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func QueryUserReadingBook(userId int, bookId string) (UserReadingBook, error) {
	userReadingBook := UserReadingBook{}
	stmt :=
		"SELECT * FROM user_reading_book WHERE user_id = $1 AND book_id = $2 "

	row := Db.QueryRow(stmt, userId, bookId)
	err := row.Scan(&userReadingBook.UserId, &userReadingBook.BookId, &userReadingBook.ReadingDate, &userReadingBook.ReadingTime)
	if err != nil {
		log.Println(err)

	}
	return userReadingBook, err
}
