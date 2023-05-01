package main

import (
	"github.com/RuhullahReza/SecondHand/db"
)


func main() {
	
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	app := Inject(DB,cld)

	app.Run(":3000")
}


