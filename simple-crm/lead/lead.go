package lead

import (
	"github.com/jinzhu/gorm"
	"github.com/gofiber/fiber"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"simple-crm/database"
)
type Lead struct {
	gorm.Model
	Name  		string  	`json:"name"`
	Company 	string 		`json:"company"`
	Email 		string 		`json:"email"`
	Phone 		int 		`json:"phone"`
}

func GetLeads(c *fiber.Ctx){
	db := database.DBConn
	var leads []Lead
	db.Find(&leads)
	if len(leads) == 0 {
		c.Status(404).SendString("No leads found")
		return
	}
	c.JSON(leads)
}
func GetLead(c *fiber.Ctx){
	id := c.Params("id")
	db := database.DBConn
	var lead Lead
	db.Find(&lead, id)
	c.JSON(lead)
}
func NewLead(c *fiber.Ctx) {
	db := database.DBConn
	lead := new(Lead)
	if err:= c.BodyParser(lead); err != nil {
		c.Status(503).Send(err)
		return
	}
	db.Create(&lead)
	c.JSON(lead)
}
func DeleteLead(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn
	var lead Lead
	db.First(&lead, id)
	if lead.Name == "" {
		c.Status(404).SendString("Lead not found")
		return
	}
	db.Delete(&lead)
	c.SendString("Lead successfully deleted")
}