package user

import (
	"net/http"
	"io/ioutil"
	"database/sql"
	"encoding/json"
	"fmt"
)

type LoginInfo struct {
	Username string
	Password string
}

type loginRes struct {
	Id int32
	Name string
  Icon string
}

func Login(db *sql.DB, r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)

	var loginInfo LoginInfo
	err = json.Unmarshal(body, &loginInfo)
	if err != nil {
		fmt.Println(err)
	}

	query := "SELECT id, name, icon FROM Users WHERE username = '" + loginInfo.Username + "' AND password = '" + loginInfo.Password + "';"
	row, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	res := &loginRes{}	
	for {
		if row.Next() == true {
			err := row.Scan(
				&res.Id,
				&res.Name,
				&res.Icon,
			)
			if err != nil {
				panic(err)
				return nil
			}
			break
		}
		return nil
	}
	fmt.Println(res)
	decodedRes, _ := json.Marshal(res)
	fmt.Println(decodedRes)
	return decodedRes
}