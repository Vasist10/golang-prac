package main 

import (
	"github.com/gofiber/fiber"
	"simple-crm/database"
	"simple-crm/lead"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

func setupRoutes(app *fiber.App){
	app.Get("/api/v1/lead" , lead.GetLeads)
	app.Get("/api/v1/lead/:id",lead.GetLead)
	app.Post("api/v1/lead",lead.NewLead)
	app.Delete("/api/v1/lead/:id",lead.DeleteLead)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")
	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Database connection established")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database migrated successfully")
}
func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(3000)
	defer database.DBConn.Close()
}