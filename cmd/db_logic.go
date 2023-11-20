package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbUser     = "root"
	dbPassword = "MyNewPass"
	dbName     = "Lab1"
	dbHost     = "127.0.0.1"
	dbPort     = "3306"
)

var db *sql.DB

func initDB() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getTables(db *sql.DB) ([]string, error) {
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_type = 'BASE TABLE'"
	rows, err := db.Query(query, dbName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func showTableStructure(db *sql.DB, tableName string) error {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return err
	}

	var results []Result

	if err := gormDB.Raw("DESCRIBE " + tableName).Scan(&results).Error; err != nil {
		return err
	}

	fmt.Printf("Structure of table %s:\n", tableName)
	for _, result := range results {
		if result.Default == "" {
			result.Default = "NULL"
		}

		fmt.Printf("Field: %s, Type: %s, Null: %s, Key: %s, Default: %s, Extra: %s\n",
			result.Field, result.Type, result.Null, result.Key, result.Default, result.Extra)
	}

	return nil
}

func callHeadCount(db *gorm.DB, num int) (string, error) {
	var result string
	if err := db.Raw("SELECT HeadCount(?) as result", num).Scan(&result).Error; err != nil {
		return "", err
	}

	return result, nil
}

// In your db_logic.go file
func calculateTotalPrice() (float64, error) {
	var total float64

	// Execute the SQL query to get the total price
	err := dbG.Raw("SELECT SUM(price) FROM bookcatalog").Scan(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func getBooksAfterYear(year int) ([]BookInfo, error) {
	var books []BookInfo

	// Execute the SQL query to get book titles and years of publication after the specified year
	err := dbG.Raw("SELECT title, year_of_publication FROM bookcatalog WHERE year_of_publication > ?", year).Scan(&books).Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

func showBooksPRICEandPAGES(db *sql.DB, maxPages, maxPrice int) error {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return err
	}

	var books []BookInfo

	// Execute a raw SQL query using GORM
	if err := gormDB.Raw("SELECT title, number_of_pages, price FROM bookcatalog WHERE number_of_pages < ? AND price < ?", maxPages, maxPrice).Find(&books).Error; err != nil {
		return err
	}

	fmt.Printf("Books with number of pages less than %d and price less than %d:\n", maxPages, maxPrice)
	for _, book := range books {
		fmt.Printf("Title: %s, Pages: %d, Price: %d\n", book.Title, book.NumberOfPages, book.Price)
	}

	return nil
}

func sortBooks(db *sql.DB) error {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return err
	}

	var books []BookInfo

	// Execute a raw SQL query using GORM to sort books by price
	if err := gormDB.Raw("SELECT title, book_library_code, price FROM bookcatalog ORDER BY price").Find(&books).Error; err != nil {
		return err
	}

	fmt.Println("Books sorted by price:")
	for _, book := range books {
		fmt.Printf("Title: %s, Library Code: %s, Price: %d\n", book.Title, book.BookLibraryCode, book.Price)
	}

	return nil
}

func changePublisher(db *sql.DB, bookLibraryCode, newPublisher string) error {
	// Execute an UPDATE SQL statement to change the publisher
	_, err := db.Exec("UPDATE bookcatalog SET publisher_id = ? WHERE book_library_code = ?", newPublisher, bookLibraryCode)
	if err != nil {
		return err
	}

	fmt.Printf("Publisher for book with Library Code %s changed to %s\n", bookLibraryCode, newPublisher)
	return nil
}

func showTableValues(db *gorm.DB, tableName string) error {
	// Retrieve all values from the specified table
	var results []map[string]interface{}
	if err := db.Table(tableName).Find(&results).Error; err != nil {
		return err
	}

	fmt.Printf("Values of table %s:\n", tableName)

	// Sort the column names for consistent order
	var columns []string
	for key := range results[0] {
		columns = append(columns, key)
	}
	sort.Strings(columns)

	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)

	// Add rows to the table
	for _, row := range results {
		var values []string
		for _, col := range columns {
			values = append(values, fmt.Sprintf("%v", row[col]))
		}
		table.Append(values)
	}

	// Set table formatting options if needed
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Render the table
	table.Render()

	return nil
}
