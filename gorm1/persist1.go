package main

import (
	"github.com/jinzhu/gorm"
	//	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// Demonstriert Benutzung der gorm Bibliothek,  basiert auf dem Quickstart von Jinzhu
func main() {
	//	db, err := gorm.Open("sqlite3", "c:/usr/sqlitedata/test.db")
	db, err := gorm.Open("postgres", "user=oosy dbname=gorm1 password=oosy2016 sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product, pr2 Product
	var products []Product
	db.First(&product) // find product with smallest id
	fmt.Printf("Erstes Insert um %v\n", product.CreatedAt)
	db.Last(&product, "code = ?", "L1212") // find product with code l1212
	db.Find(&products)
	fmt.Printf("%d Produkte\n", len(products))
	db.First(&pr2) // find product with highest id
	fmt.Printf("Letztes Insert um %v\n", pr2.CreatedAt)

	// Update product's price to 2345
	product.Price = 2345
	db.Save(&product)
	//	db.Model(&pr2).Update("Price", pr2.Price + 1230)
	//	fmt.Printf("Produkt: %v\n", pr2)

	conflictingChanges(db)
	// Delete - delete product
	//db.Delete(&product)
}

// finde heraus, wie die Datenbank auf konkurrierende Änderungen reagiert
func conflictingChanges(db *gorm.DB) {
	var p1, p2, p3 Product

	// sicherstellen, dass beide Variablen mit gleichen Daten
	// aus der Datenbank gelesen werden
	// beginne Transaktion: Alle
	// Datenbankoperationen nur mit tx!
	tx := db.Begin()
	// lies Objekt aus DB in p1
	tx.First(&p1)
	tx.First(&p2)
	// ändere p1
	p1.Price = 1111
	p1.Code = "P4711"
	// sichere p1 zurück in die DB
	db.Save(&p1)
	// beende Transaktion
	tx.Commit()
	//time.Sleep(1 * time.Millisecond)

	p1.Price = 1111
	p1.Code = "P4711"

	p2.Price = 2222
	p2.Code = "42-HAL"

	ar := db.Save(&p1).RowsAffected
	if ar > 0 {
		fmt.Printf("Saved to DB, RowsAffected = %d\n", ar)
	} else {
		fmt.Printf("Not saved to DB, RowsAffected = %d\n", ar)
	}

	// überschreibt p1
	//db.Save(&p2)

	// überschreibt nicht -> optimistic locking
	// feststellen, ob Daten in DB geändert wurden:

	// Zeitstempel auf fremden update überprüfen
	if db.Where("updated_at = ?", p2.UpdatedAt).
		// nur speichern, wenn keine Änderung
		Save(&p2).
		// Anzahl der geänderten Zeilen in DB testen
		RowsAffected == 0 {
		// geänderte Daten integrieren ...
		fmt.Printf("Not saved to DB\n")
	} else {
		// alles OK, kein Update-Konflikt
		fmt.Printf("Saved to DB\n")
	}

	db.First(&p3)
	fmt.Printf("Produkt in DB: %s, %d\n", p3.Code, p3.Price)
}
