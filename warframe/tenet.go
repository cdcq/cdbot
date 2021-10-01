package warframe

import (
	"cdbot/helpers"
	"database/sql"
)

func TenetResponse() string {
	funcName := "tenet.go: TenetResponse"

	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		helpers.AddLog(funcName, "open database", err)
		return ""
	}
	prep, err := db.Prepare("SELECT CONTENT FROM WF_MISC WHERE NAME=?")
	if err != nil {
		helpers.AddLog(funcName, "prepare query", err)
		return ""
	}
	rows, err := prep.Query("Tenet")
	if err != nil {
		helpers.AddLog(funcName, "execute query", err)
		return ""
	}
	var res string
	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			helpers.AddLog(funcName, "scan row", err)
			return ""
		}
	}
	return res
}

func TenetUpdate(content string) error {
	funcName := "tenet.go: TenetUpdate"

	db, err := sql.Open("sqlite3", "../database.db")
	if err != nil {
		helpers.AddLog(funcName, "open database", err)
		return err
	}
	prep, err := db.Prepare("UPDATE WF_MISC SET CONTENT=? WHERE NAME=?")
	if err != nil {
		helpers.AddLog(funcName, "prepare query", err)
		return err
	}
	_, err = prep.Exec(content, "Tenet")
	if err != nil {
		helpers.AddLog(funcName, "prepare query", err)
		return err
	}
	return nil
}
