package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Product struct{
	ID int64
	Name string
	Price float64
}

func main() {
	// Celar all product record
	var err = ClearProducts()
	if err != nil {
		log.Fatal(err)
	}

	// do Insert
	fmt.Println("")
	IDs, err := InsertProducts("Test Name", 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Succes insert with Id", IDs)

	// Get All Products
	fmt.Println("")
	products, err := GetProducts()
	if err != nil {
		log.Fatal("Gett All product err :", err)
	}
	for _, product := range products {
		fmt.Printf("ID: %d, Name: %s, Price %.2f\n", product.ID, product.Name, product.Price)
	}

	// Get Product By ID 1
	fmt.Println("")
	product, err := GetProduct(1)
	if err != nil {
		log.Fatal("Get product By ID", 1, " error :", err)
	}
	fmt.Printf("Product found: ID: %d, Name: %s, Price: %.2f\n", product.ID, product.Name, product.Price)

	// Delete product by ID
	fmt.Println("")
	err = DeleteProduct(1)
	if err !=nil {
		log.Fatal(err)
	}

}

// function to connect database
func Conn() (*sql.DB, error) {
	// Load env
	err := LoadVar()
	if err != nil {
		log.Fatal(err)
	}
	// get var
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// dsn := "myuser:myuser@tcp(127.0.0.1:3306)/training"
	dsn := user+":"+pass+"@tcp("+host+":"+port+")/"+name
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	// ping connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	println("Succes connected")
	return db, nil
}

// Load Variable
func LoadVar() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error load .env file")
	}
	return nil
}

// function insert to product
func InsertProducts(name string, price float64) (int64, error) {
	db, err := Conn();
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO products (name, price) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	
	res, err := stmt.Exec(name,price)
	if err != nil {
		return 0, err
	} 

	resultID, _ := res.LastInsertId()

	return resultID, nil
}

// Function get all products
func GetProducts() ([]Product, error) {
	// Connect to DB
	db, err := Conn();
	if err != nil {
		log.Fatal(err)
	}
	
	// Query 
	query, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	var products []Product
	for query.Next() {
		var p Product
		err := query.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// Get Product By Id
func GetProduct(id int64) (*Product, error){
	db, err := Conn()
	if err != nil {
		log.Fatal(err)
	}

	var product Product
	err = db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No ID in list products", id)
		}
		return nil, err
	}
	defer db.Close()
	return &product, nil	
}

// Delete product by ID
func DeleteProduct(id int64)(error){
	db, err := Conn()
	if err != nil {
		log.Fatal(err)
	}

	// prepare statement
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	row, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if row == 0 {
		return fmt.Errorf("No product found by id %d", id)
	}
	fmt.Printf("Deleted product by Id: ", id)
	return nil
}

// Clear Products
func ClearProducts() error {
	db, err := Conn()
	if err != nil {
		log.Fatal(err)
		return err
	}

	stmt, err := db.Prepare("DELETE FROM products")
	if err != nil {
		log.Fatal(err)
		return err
	}
	res, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Succes clear products.\n", row)
	return nil
}