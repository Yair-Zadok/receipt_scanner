This is a SQLite database with various getters and setters public functions to handle batches of user receipt data. 

The database is architectured such that a user's email is central to all their respecive data, this was chosen as an email is a unique identifiers, thus removing the need for tracking integer IDs.



Functions:
Section: Setup

func Setup_db(db *sql.DB) -> error 
    - This function takes a pointer to an sql.DB and creates the necessary tables for storing user accounts and batches of user receipts.

    - Created tables have the following struture:
        Table: users
        Primary Key: email TEXT
        Fields: None


        Table: user_data
        Primary Key: INT entryId
        Fields: email TEXT, batch INTEGER, subtotal TEXT, total TEXT, tax TEXT, tips TEXT, date TEXT,
        supplier TEXT, account TEXT, encoded_image TEXT
        Foreign Key: email REFERENCES users email



Section: Readers

- Checks if a user with a certain email exists in the database
func Exists_email(db *sql.DB, email string) -> (bool, error)


- Retrieves a list of all user emails in the database
func Get_user_emails(db *sql.DB) -> ([]string, error)


- Retrieves a number id associated with the highest batch per a user email
func Get_highest_batch(db *sql.DB, user_email string) -> (int, error)


- Gets user receipt data corresponding to their last inputted batch
func Get_last_batch_entries(db *sql.DB, user_email string) -> ([]Receipt, error)


- Gets all receipt data per a user email
func Get_all_entries(db *sql.DB, user_email string) -> ([]Receipt, error)
 

- Gets all entryId's in a user's last batch
func Get_last_batch_ids(db *sql.DB, user_email string) -> ([]int, error) 



Section: Writers

- Adds a user
func Set_user(db *sql.DB, email string) -> error

- Adds a list of reciepts as a batch to a certain user
func Set_entry(db *sql.DB, email string, receipts []Receipt) -> error

- Updates the last batch of entries with a new list of data
func Update_entry(db *sql.DB, email string, receipts []Receipt) -> error










Copyright 2024, Yair Zadok, All rights reserved.

