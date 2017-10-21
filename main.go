package main

import (
    "log"
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"

    "github.com/gorilla/websocket"

    message "./message"
    ws "./websocket"
)

func main() {
    // db接続
    db, err := sql.Open("postgres", "user=twochat_client dbname=twochat sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    hub := ws.Hub{}
    hub.Connection = map[ws.Client]*websocket.Conn{}

    http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set( "Access-Control-Allow-Credentials", "true" )
        w.Header().Set( "Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization" )
        w.Write(message.MessageAction(db, r))
    })

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        // hubへclientを追加
        var upgrader = websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
            CheckOrigin: func(r *http.Request) bool {
                return true
            },
        }
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Print(err)
        }
        c := ws.Client{conn}
        fmt.Println(c)
        fmt.Println(conn)
        hub.Connection[c] = conn
        fmt.Println(hub.Connection)
        ws.Connection(&hub)

        ws.Cast(db, w, r, &hub, c)
    })

    fmt.Println("listen on port:8080")
    http.ListenAndServe(":8080", nil)

}
