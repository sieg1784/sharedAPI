package data

import (
	"bookAPI/apiservice/utils"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type ResultBook struct {
	BookId      string               `json:"bookId"`
	BookName    string               `json:"bookName"`
	ImageUrl    null.String          `json:"imageUrl"`
	Category    null.String          `json:"category"`
	Description null.String          `json:"description"`
	Coworker    null.String          `json:"coworker"`
	OnStockDate utils.TimeStampInt64 `json:"onStockDate"`
	TimeLength  null.Int             `json:"timeLength"`
	ReadingTime null.Int             `json:"readingTime"`
}

type ResultSet struct {
	BookId      string
	BookName    string
	ImageUrl    null.String
	Category    null.String
	Description null.String
	Coworker    null.String
	OnStockDate null.Time
	TimeLength  null.Int
	ReadingTime null.Int
}

func QueryNowLearning(userId int) ([]ResultBook, error) {

	var nowLearningList []ResultBook
	stmt :=
		"SELECT " +
			"b.book_id, b.book_name, b.image_url, b.category, b.description, b.onstock_date, b.time_length, a.reading_time " +
			"FROM user_reading_book a JOIN book b on a.book_id = b.book_id " +
			"WHERE a.user_id = $1 ORDER BY reading_date DESC limit 6"
	rows, err := Db.Query(stmt, userId)
	if err != nil {
		fmt.Println(err)
		return nowLearningList, err
	}
	defer rows.Close()
	for rows.Next() {
		var resultSet ResultSet
		err := rows.Scan(&resultSet.BookId, &resultSet.BookName, &resultSet.ImageUrl, &resultSet.Category, &resultSet.Description, &resultSet.OnStockDate, &resultSet.TimeLength, &resultSet.ReadingTime)
		if err != nil {
			fmt.Println(err)
			return nowLearningList, err
		}

		resultbook := ResultBook{}
		resultSet2ResultBook(&resultbook, resultSet)
		nowLearningList = append(nowLearningList, resultbook)

	}

	return nowLearningList, err
}

func QueryRecommendation(userId int, page int, pageSize int) ([]ResultBook, error) {
	return QueryBookList(userId, page, pageSize)
}

func QueryBookList(userId int, page int, pageSize int) ([]ResultBook, error) {
	rowStart := ((page - 1) * pageSize) + 1
	rowEnd := rowStart + pageSize - 1
	var recommendationList []ResultBook
	stmt :=
		"SELECT " +
			"	a.book_id, a.book_name, a.image_url, a.category, a.description, a.onstock_date, a.time_length, a.reading_time " +
			"FROM ( " +
			"	SELECT " +
			"		a.book_id, a.book_name, a.image_url, a.category, a.description, a.onstock_date, a.time_length, b.reading_time, " +
			" 		row_number() OVER (ORDER BY a.onstock_date DESC) as order_id " +
			"	FROM " +
			"		book a LEFT JOIN user_reading_book b ON a.book_id = b.book_id AND b.user_id = $1 " +
			") a " +
			"WHERE a.order_id BETWEEN $2 AND $3"
	rows, err := Db.Query(stmt, userId, rowStart, rowEnd)
	if err != nil {
		fmt.Println(err)
		return recommendationList, err
	}

	defer rows.Close()
	for rows.Next() {
		var resultSet ResultSet

		err := rows.Scan(&resultSet.BookId, &resultSet.BookName, &resultSet.ImageUrl, &resultSet.Category, &resultSet.Description, &resultSet.OnStockDate, &resultSet.TimeLength, &resultSet.ReadingTime)
		if err != nil {
			fmt.Println(err)
			return recommendationList, err
		}
		resultbook := ResultBook{}
		resultSet2ResultBook(&resultbook, resultSet)
		recommendationList = append(recommendationList, resultbook)
	}

	return recommendationList, err
}

func QueryCowork(userId int, page int, pageSize int) ([]ResultBook, error) {
	rowStart := ((page - 1) * pageSize) + 1
	rowEnd := rowStart + pageSize - 1
	var recommendationList []ResultBook
	stmt :=
		"SELECT " +
			" a.book_id, a.book_name, a.coworker, a.image_url, a.category, a.description, a.onstock_date, a.time_length, a.reading_time " +
			" FROM ( " +
			"SELECT " +
			" b.book_id, b.book_name, a.coworker, b.image_url, b.category, b.description, b.onstock_date, b.time_length, c.reading_time,  " +
			" row_number() OVER (ORDER BY b.onstock_date DESC) as order_id " +
			" FROM " +
			" cowork_book a " +
			" JOIN book b ON a.book_id = b.book_id " +
			" LEFT JOIN user_reading_book c ON a.book_id = c.book_id and c.user_id = $1 " +
			" ) a " +
			"WHERE a.order_id BETWEEN $2 AND $3"
	rows, err := Db.Query(stmt, userId, rowStart, rowEnd)
	if err != nil {
		fmt.Println(err)
		return recommendationList, err
	}

	defer rows.Close()
	for rows.Next() {
		var resultSet ResultSet

		err := rows.Scan(&resultSet.BookId, &resultSet.BookName, &resultSet.Coworker, &resultSet.ImageUrl, &resultSet.Category, &resultSet.Description, &resultSet.OnStockDate, &resultSet.TimeLength, &resultSet.ReadingTime)
		if err != nil {
			fmt.Println(err)
			return recommendationList, err
		}
		resultbook := ResultBook{}
		resultSet2ResultBook(&resultbook, resultSet)
		recommendationList = append(recommendationList, resultbook)
	}

	return recommendationList, err
}

func resultSet2ResultBook(resultBook *ResultBook, resultSet ResultSet) {
	resultBook.BookId = resultSet.BookId
	resultBook.BookName = resultSet.BookName
	resultBook.ImageUrl = resultSet.ImageUrl
	resultBook.Category = resultSet.Category
	resultBook.Description = resultSet.Description
	resultBook.Coworker = resultSet.Coworker
	resultBook.OnStockDate = utils.Time2MilliTimestamp(resultSet.OnStockDate)
	resultBook.TimeLength = resultSet.TimeLength
	resultBook.ReadingTime = resultSet.ReadingTime
}

func QueryPlayerData(userId int) ([]ResultBook, error) {
	var recommendationList []ResultBook
	stmt :=
		"SELECT " +
			" b.book_id, b.book_name, a.coworker, b.image_url, b.category, b.description, b.onstock_date, b.time_length, c.reading_time " +
			" FROM " +
			" cowork_book a " +
			" JOIN book b ON a.book_id = b.book_id " +
			" LEFT JOIN user_reading_book c ON a.book_id = c.book_id and c.user_id = $1"
	rows, err := Db.Query(stmt, userId)
	if err != nil {
		fmt.Println(err)
		return recommendationList, err
	}

	defer rows.Close()
	for rows.Next() {
		var resultSet ResultSet

		err := rows.Scan(&resultSet.BookId, &resultSet.BookName, &resultSet.Coworker, &resultSet.ImageUrl, &resultSet.Category, &resultSet.Description, &resultSet.OnStockDate, &resultSet.TimeLength, &resultSet.ReadingTime)
		if err != nil {
			fmt.Println(err)
			return recommendationList, err
		}
		resultbook := ResultBook{}
		resultSet2ResultBook(&resultbook, resultSet)
		recommendationList = append(recommendationList, resultbook)
	}

	return recommendationList, err
}

func QueryBookByInputString(userId int, searchString string) ([]ResultBook, error) {

	var recommendationList []ResultBook
	stmt :=
		"SELECT " +
			"	a.book_id, a.book_name, a.image_url, a.category, a.description, a.onstock_date, a.time_length, a.reading_time " +
			"FROM ( " +
			"	SELECT " +
			"		a.book_id, a.book_name, a.book_engname, a.author_name, a.author_engname, a.image_url, a.category, a.description, a.onstock_date, a.time_length, b.reading_time, " +
			" 		row_number() OVER (ORDER BY a.onstock_date DESC) as order_id " +
			"	FROM " +
			"		book a LEFT JOIN user_reading_book b ON a.book_id = b.book_id AND b.user_id = $1 " +
			") a " +
			" WHERE " +
			"	a.book_name like '%' || $2 || '%' " +
			" 	or a.book_engname like '%' || $3 || '%' " +
			" 	or a.author_name like '%' || $4 || '%' " +
			" 	or a.author_engname like '%' || $5 || '%' " +
			" 	or a.description like '%' || $6 || '%' " +
			"	or a.category like '%' || $7 || '%'"
	rows, err := Db.Query(stmt, userId, searchString, searchString, searchString, searchString, searchString, searchString)
	if err != nil {
		fmt.Println(err)
		return recommendationList, err
	}

	defer rows.Close()
	for rows.Next() {
		var resultSet ResultSet

		err := rows.Scan(&resultSet.BookId, &resultSet.BookName, &resultSet.ImageUrl, &resultSet.Category, &resultSet.Description, &resultSet.OnStockDate, &resultSet.TimeLength, &resultSet.ReadingTime)
		if err != nil {
			fmt.Println(err)
			return recommendationList, err
		}
		resultbook := ResultBook{}
		resultSet2ResultBook(&resultbook, resultSet)
		recommendationList = append(recommendationList, resultbook)
	}

	return recommendationList, err
}

func QueryRecommendationBookName(str string) ([]string, error) {
	var resultStringList []string
	stmt := "SELECT book_name FROM book WHERE book_name LIKE $1 || '%'"
	rows, err := Db.Query(stmt, str)
	if err != nil {
		fmt.Println(err)
		return resultStringList, err
	}
	for rows.Next() {
		resultString := ""
		err := rows.Scan(&resultString)
		if err != nil {
			fmt.Println(err)
			return resultStringList, err
		}

		resultStringList = append(resultStringList, resultString)
	}
	return resultStringList, err
}
