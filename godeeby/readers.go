// Copyright 2024, Yair Zadok, All rights reserved.

package godeeby

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// Checks if a user with a certain email exists in the database
func Exists_email(db *sql.DB, email string) (bool, error) {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return false, err
    }

    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM user_data WHERE email = ?", email).Scan(&count)

    if (err != nil ) { fmt.Println(err); return false, err }
    return count > 0, nil
}

// Retrieves a list of all user emails in the database
func Get_user_emails(db *sql.DB) ([]string, error) {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return []string{}, err
    }

    emails := []string{}

    rows, err := db.Query("SELECT email FROM users")
    if (err != nil ) { fmt.Println(err); return []string{}, err }
    defer rows.Close()

    var email string
    for rows.Next() {
        rows.Scan(&email)
        emails = append(emails, email)
    }
    
    return emails, err
}

// Retrieves a number id associated with the highest batch per a user email
func Get_highest_batch(db *sql.DB, user_email string) (int, error) {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return 0, err
    }

    rows, err := db.Query(`SELECT MAX(batch) FROM user_data WHERE email = ?`, user_email)
    defer rows.Close() 

    if (err != nil ) { fmt.Println(err); return -999, err }

    var batch int  
    
    for rows.Next() {
        rows.Scan(&batch)
        return batch, nil
    }

    return 0, nil 
}

// Gets user receipt data corresponding to their last inputted batch
func Get_last_batch_entries(db *sql.DB, user_email string) ([]Receipt, error) {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return []Receipt{}, err
    }

    receipts := []Receipt{}

    rows, err := db.Query(`SELECT subtotal, total, tax, tips, date, 
    supplier, account, encoded_image FROM user_data WHERE email=?
    AND batch = (SELECT MAX(batch) FROM user_data WHERE email = ?)`, user_email, user_email)
    if (err != nil ) { fmt.Println(err); return []Receipt{}, err }
    defer rows.Close()

    var subtotal string
    var total string
    var tax string
    var tips string
    var date string
    var supplier string
    var account string
    var encoded_image string
    
    for rows.Next() {
        rows.Scan(&subtotal, &total, &tax, &tips, &date, &supplier, &account, &encoded_image)
        receipts = append(receipts, Receipt{subtotal, total, tax, tips, date, supplier, account, encoded_image})
    } 

    return receipts, nil
}

// Gets all receipt data per a user email
func Get_all_entries(db *sql.DB, user_email string) ([]Receipt, error) {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return []Receipt{}, err
    }

    receipts := []Receipt{}

    rows, err := db.Query(`SELECT subtotal, total, tax, tips, date, 
    supplier, account, encoded_image FROM user_data WHERE email=?`, user_email)
    defer rows.Close()

    if (err != nil ) { fmt.Println(err); return []Receipt{}, err }

    var subtotal string
    var total string
    var tax string
    var tips string
    var date string
    var supplier string
    var account string
    var encoded_image string

    for rows.Next() {
        rows.Scan(&subtotal, &total, &tax, &tips, &date, &supplier, &account, &encoded_image)
        receipts = append(receipts, Receipt{subtotal, total, tax, tips, date, supplier, account, encoded_image})
    } 

    return receipts, nil
}

// Gets all entryId's in a user's last batch
func Get_last_batch_ids(db *sql.DB, user_email string) ([]int, error) {
    err := db.Ping()
    if err != nil {
        fmt.Println("DB connection close:", err)
        return []int{}, err
    }

    ids := []int{}

    rows, err := db.Query(`SELECT entryId FROM user_data WHERE email=?
    AND batch = (SELECT MAX(batch) FROM user_data WHERE email = ?)`, user_email, user_email)
    if (err != nil ) { fmt.Println(err); return []int{}, err }
    defer rows.Close()

    var id int

    for rows.Next() {
        rows.Scan(&id)
        ids = append(ids, id)
    } 

    return ids, nil
}


