package message

import (
	"database/sql"
	"net/http"
	"io/ioutil"	
	"fmt"
	"encoding/json"
	"strconv"
)

type GetMessage struct {
    Id int32
	Message string
	SenderId int32
    Name string
    Icon string
    Datetime string
}

type PostMessage struct {
    Sender int
    Message string
    Datetime string
}


func New(rows *sql.Rows) *GetMessage {
    if rows.Next() == true {
        res := &GetMessage{}
        err := rows.Scan(
            &res.Id,
			&res.Message,
			&res.SenderId,
            &res.Name,
            &res.Icon,
            &res.Datetime,
        )
        if err != nil {
            fmt.Println(err)
            return nil
        }
        return res
    }
    return nil
}

func MessageAction(db *sql.DB, r *http.Request) []byte {
	switch r.Method {
	case "GET":
		query := "SELECT m.id, m.message, u.id, u.name, u.icon, m.datetime FROM messages as m INNER JOIN users as u ON m.sender_id = u.id ORDER BY m.id"
		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}
		var result []GetMessage
		for {
			rower := New(rows)
			if rower == nil {
				break
			}
			result = append(result, *rower)
		}

		decodedResult, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		return decodedResult
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var row PostMessage
		err = json.Unmarshal(body, &row)
		if err != nil {
			panic(err)
		}
		query := "INSERT INTO messages VALUES ("  + strconv.Itoa(row.Sender) + ", '" + row.Message + "', '" + row.Datetime + "')"
		_, err = db.Query(query)
		if err != nil {
			panic(err)
		}

		decodedRow, _ := json.Marshal(row)
		return decodedRow
	}
	return nil
}