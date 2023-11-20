package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbG *gorm.DB

func SetDB(gormDB *gorm.DB) {
	dbG = gormDB
}

type Result struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:Key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}

type BookInfo struct {
	ID                int    `gorm:"column:idBC"`
	BookLibraryCode   string `gorm:"column:book_library_code"`
	Author            string `gorm:"column:author"`
	Title             string `gorm:"column:title"`
	Publisher         string `gorm:"column:publisher"`
	YearOfPublication int    `gorm:"column:year_of_publication"`
	BookType          string `gorm:"column:book_type"`
	Edition           string `gorm:"column:edition"`
	NumberOfPages     int    `gorm:"column:number_of_pages"`
	Topic             string `gorm:"column:topic"`
	Price             int    `gorm:"column:price"`
}

var showTablesCmd = &cobra.Command{
	Use:   "show-tables",
	Short: "Show all tables in the database",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer db.Close()

		tables, err := getTables(db)
		if err != nil {
			fmt.Println("Error getting tables:", err)
			os.Exit(1)
		}

		fmt.Println("Tables in the database:")
		for _, table := range tables {
			fmt.Println(table)
		}
	},
}

var showStructureCmd = &cobra.Command{
	Use:   "show-structure",
	Short: "Show the structure of a specific table",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: dbCLI show-structure <table_name>")
			os.Exit(1)
		}

		tableName := args[0]

		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer db.Close()

		err = showTableStructure(db, tableName)
		if err != nil {
			fmt.Println("Error showing table structure:", err)
			os.Exit(1)
		}
	},
}

var callHeadCountCMD = &cobra.Command{
	Use:   "call-HeadCount",
	Short: "Call the HeadCount function in the database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: dbCLI call-HeadCount <num>")
			return
		}

		num, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid num. Please provide a valid integer.")
			return
		}

		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			return
		}
		defer db.Close()

		result, err := callHeadCount(dbG, num)
		if err != nil {
			fmt.Println("Error calling HeadCount function:", err)
			return
		}

		fmt.Println(result)
	},
}

var calculateTotalPriceCmd = &cobra.Command{
	Use:   "calculate-total-price",
	Short: "Calculate the sum of prices in the bookcatalog table",
	Run: func(cmd *cobra.Command, args []string) {
		total, err := calculateTotalPrice()
		if err != nil {
			fmt.Println("Error calculating total price:", err)
			return
		}

		fmt.Printf("Total price of all books: $%.2f\n", total)
	},
}

var getBooksAfterYearCmd = &cobra.Command{
	Use:   "get-books-after-year",
	Short: "Get book names and years of publication where the year is higher than the provided year",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		year, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid year argument. Please provide a valid number.")
			return
		}
		books, err := getBooksAfterYear(year)
		if err != nil {
			fmt.Println("Error getting books after the specified year:", err)
			return
		}
		fmt.Printf("Books published after %d:\n", year)
		for _, book := range books {
			fmt.Printf("Name: %s, Year of Publication: %d\n", book.Title, book.YearOfPublication)
		}
	},
}

var showANDCmd = &cobra.Command{
	Use:   "show-books-Mpages-Mprice",
	Short: "Show books based on criteria",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Usage: dbCLI show-books <max_pages> <max_price>")
			os.Exit(1)
		}

		maxPages, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error parsing max_pages:", err)
			os.Exit(1)
		}

		maxPrice, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error parsing max_price:", err)
			os.Exit(1)
		}

		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer db.Close()

		err = showBooksPRICEandPAGES(db, maxPages, maxPrice)
		if err != nil {
			fmt.Println("Error showing books:", err)
			os.Exit(1)
		}
	},
}

var sortBooksCmd = &cobra.Command{
	Use:   "sort-books",
	Short: "Sort and show books based on price",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer db.Close()

		err = sortBooks(db)
		if err != nil {
			fmt.Println("Error sorting and showing books:", err)
			os.Exit(1)
		}
	},
}

var changePublisherCmd = &cobra.Command{
	Use:   "change-publisher <book_library_code> <new_publisher>",
	Short: "Change the publisher of a book",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer db.Close()

		bookLibraryCode := args[0]
		newPublisher := args[1]

		err = changePublisher(db, bookLibraryCode, newPublisher)
		if err != nil {
			fmt.Println("Error changing publisher:", err)
			os.Exit(1)
		}
	},
}

var showValuesCmd = &cobra.Command{
	Use:   "show-values",
	Short: "Show the values of a specific table",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: dbCLI show-values <table_name>")
			os.Exit(1)
		}

		tableName := args[0]

		db, err := initDB()
		if err != nil {
			fmt.Println("Error initializing database:", err)
			os.Exit(1)
		}
		defer db.Close()

		err = showTableValues(dbG, tableName)
		if err != nil {
			fmt.Println("Error showing table values:", err)
			os.Exit(1)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "dbCLI",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(showTablesCmd)
	rootCmd.AddCommand(showStructureCmd)
	rootCmd.AddCommand(callHeadCountCMD)
	rootCmd.AddCommand(calculateTotalPriceCmd)
	rootCmd.AddCommand(getBooksAfterYearCmd)
	rootCmd.AddCommand(showANDCmd)
	rootCmd.AddCommand(sortBooksCmd)
	rootCmd.AddCommand(changePublisherCmd)
	rootCmd.AddCommand(showValuesCmd)

}
