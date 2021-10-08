package data

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type BookExtContent struct {
	ExtType    null.String `json:"extType"`
	ExtContent null.String `json:"extContent"`
}

func QueryBookAudioUrl(bookId string) (string, error) {
	var url string
	row := Db.QueryRow("select video_url from book where book_id = $1", bookId)
	err := row.Scan(&url)
	if err != nil {
		fmt.Println(err)
		return url, err
	}

	return url, err
}

func QueryBookExtContent(bookId string) ([]BookExtContent, error) {
	var bookExtContentList []BookExtContent
	rows, err := Db.Query("select ext_type, ext_content from book_ext_content where book_id = $1", bookId)
	if err != nil {
		fmt.Println(err)
		return bookExtContentList, err
	}
	defer rows.Close()
	for rows.Next() {
		var bookExtContent BookExtContent
		err := rows.Scan(&bookExtContent.ExtType, &bookExtContent.ExtContent)
		if err != nil {
			fmt.Println(err)
			return bookExtContentList, err
		}

		bookExtContentList = append(bookExtContentList, bookExtContent)
	}

	return bookExtContentList, err
}
