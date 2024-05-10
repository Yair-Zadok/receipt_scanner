// Copyright 2024, Yair Zadok, All rights reserved.

package godeeby 

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// Makes the tables necessary for storing user accounts and receipt batches
func Setup_db(db *sql.DB) error {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return err
    }

    // TABLE: users
    statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (email TEXT NOT NULL PRIMARY KEY UNIQUE)")
    if (err != nil ) { fmt.Println(err); return err }
    
    _, err = statement.Exec()
    if (err != nil ) { fmt.Println(err); return err }
    
    // TABLE: user_data
    statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS user_data 
    (entryId INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT NOT NULL, batch INTEGER NOT NULL, subtotal TEXT NOT NULL,
    total TEXT NOT NULL, tax TEXT NOT NULL, tips TEXT NOT NULL, date TEXT NOT NULL, supplier TEXT NOT NULL, 
    account TEXT NOT NULL, encoded_image TEXT NOT NULL, FOREIGN KEY (email) REFERENCES users(email))`)
    if (err != nil ) { fmt.Println(err); return err  }
    
    _, err = statement.Exec()
    if (err != nil ) { fmt.Println(err); return err }
    
    return nil
}

