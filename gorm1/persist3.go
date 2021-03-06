package main

import (
	"fmt"
	"github.com/geobe/go4web/gorm1/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "user=oosy dbname=gorm3 password=oosy2016 sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&model.Trip{}, &model.Person{})

	kirk := model.SomePersons[0]
	lara := model.SomePersons[1]
	kirk.Trips = append(kirk.Trips, model.SomeTrips...)
	lara.Trips = append(lara.Trips, model.SomeTrips[1])

	db.Save(&kirk)
	db.Save(&lara)

	// query
	var kiki, riki, rara, lada, tripRider model.Person
	db.First(&kiki, kirk.ID)
	db.First(&riki, kirk.ID).Related(&riki.Trips)
	db.Preload("Trips").First(&lada, lara.ID)
	db.First(&rara, lara.ID)
	// von Person -> Trip
	db.Model(&rara).Related(&rara.Trips)
	// von Trip -> Person
	db.Model(&lara.Trips[0]).Related(&tripRider)

	fmt.Println(kirk)
	fmt.Println(kiki)
	fmt.Println(riki)

	fmt.Println(lara)
	fmt.Println(lada)
	fmt.Println(rara)
	fmt.Println(tripRider)

	db.Delete(model.Person{})
	db.Delete(model.Trip{})

}
