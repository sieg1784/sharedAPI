package data

import "fmt"

type BookCategory struct {
	Category    string `json:"category"`
	DisplayText string `json:"displayText"`
}

func QueryBookCategory() ([]BookCategory, error) {
	var bookCategoryList []BookCategory
	stmt := "SELECT category, display_text FROM book_category"
	rows, err := Db.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return bookCategoryList, err
	}
	defer rows.Close()
	for rows.Next() {
		var bookCategory BookCategory
		err := rows.Scan(&bookCategory.Category, &bookCategory.DisplayText)
		if err != nil {
			fmt.Println(err)
			return bookCategoryList, err
		}

		bookCategoryList = append(bookCategoryList, bookCategory)
	}

	return bookCategoryList, err
}
