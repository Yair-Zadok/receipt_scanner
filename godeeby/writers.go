// Copyright 2024, Yair Zadok, All rights reserved.

package godeeby

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// Adds a user
func Set_user(db *sql.DB, email string) error {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection closed:", err)
        return err
    }

    exists, err := Exists_email(db, email)
    if (err != nil ) { fmt.Println(err); return err }

    if (!exists) {
        statement, err := db.Prepare("INSERT INTO users (email) VALUES (?)")
        if (err != nil ) { fmt.Println(err); return err }

        _, err = statement.Exec(email)
        if (err != nil ) { fmt.Println(err); return err }
    }
    return nil
}

// Adds a list of reciepts as a batch to a certain user
func Set_entry(db *sql.DB, email string, receipts []Receipt) error {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection closed:", err)
        return err
    }

    statement, err := db.Prepare(`INSERT INTO user_data (email, batch, subtotal, total, tax, tips, date, 
    supplier, account, encoded_image)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
    if (err != nil ) { fmt.Println(err); return err }
    defer statement.Close()

    batch, err := Get_highest_batch(db, email) 
    if (err != nil ) { fmt.Println(err); return err }

    batch += 1

    for _, receipt := range receipts {
        _, err = statement.Exec(email, batch, receipt.Subtotal, receipt.Total, receipt.Tax, receipt.Tips, 
        receipt.Date, receipt.Supplier, receipt.Account, receipt.Encoded_image)
        if (err != nil ) { fmt.Println(err); return err }
    }

    return nil
}

// Updates the last batch of entries with a new list of data
func Update_entry(db *sql.DB, email string, receipts []Receipt) error {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection closed:", err)
        return err
    }

    statement, err := db.Prepare(`UPDATE user_data 
    SET subtotal=?, total=?, tax=?, tips=?, date=?, supplier=?, account=? 
    WHERE id=? AND email=?`)
    if (err != nil ) { fmt.Println(err); return err }
    defer statement.Close()

    ids, err := Get_last_batch_ids(db, email) 
    if (err != nil ) { fmt.Println(err); return err }

    for i, receipt := range receipts { 
        _, err = statement.Exec(receipt.Subtotal, receipt.Total, receipt.Tax, receipt.Tips, 
        receipt.Date, receipt.Supplier, receipt.Account, ids[i], email)
        if (err != nil ) { fmt.Println(err); return err }
    }
    
    return nil
}
