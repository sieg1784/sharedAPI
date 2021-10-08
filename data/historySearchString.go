package data

import "fmt"

type HistorySearchString struct {
	SearchString string
	Count        int
}

func InsertHistorySearchString(str string, count int) (bool, error) {
	stmt := "INSERT INTO history_search_string (search_string, count) " +
		" VALUES ($1, $2) "
	_, err := Db.Exec(stmt, str, count)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func UpdateHistorySearchString(historySearchString HistorySearchString) (bool, error) {
	stmt := "UPDATE history_search_string SET count = $1 WHERE search_string = $2"
	_, err := Db.Exec(stmt, historySearchString.Count, historySearchString.SearchString)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func QueryHistorySearchString(str string) (HistorySearchString, error) {
	historySearchString := HistorySearchString{}
	stmt := "SELECT * FROM history_search_string WHERE search_string = $1"
	row := Db.QueryRow(stmt, str)
	err := row.Scan(&historySearchString.SearchString, &historySearchString.Count)
	return historySearchString, err
}

func QueryRecommendationString(str string, threshold int) ([]string, error) {

	var resultStringList []string
	stmt :=
		"SELECT * FROM ( " +
			" SELECT book_name AS recommendation_string FROM book WHERE book_name LIKE $1 || '%' " +
			" UNION " +
			" SELECT search_string AS recommendation_string FROM history_search_string WHERE search_string LIKE $1 || '%' AND count >= $2 " +
			" ) a "
	rows, err := Db.Query(stmt, str, threshold)
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
