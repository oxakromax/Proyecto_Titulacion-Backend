package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/oxakromax/ProyectoTitulo-Monitoreo/Utils"
	"os"
)

var (
	PostgreUser = os.Getenv("Postgre-User")
	PostgrePass = os.Getenv("Postgre-Pass")
	PostgreHost = os.Getenv("Postgre-Host")
	PostgrePort = os.Getenv("Postgre-Port")
	PostgreDB   = os.Getenv("Postgre-DB")
	db          *sql.DB
)

func main() {
	var err error
	Utils.QueryAuth.ClientId = os.Getenv("APP-ID")
	Utils.QueryAuth.ClientSecret = os.Getenv("APP-Secret")
	Utils.QueryAuth.Scope = os.Getenv("APP-Scope")
	Utils.UipathOrg.Id = os.Getenv("Orchestrator-ID")
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return
	}
	fmt.Print("Connected to database")
	selectSentence := `Select "id","nombre" from "clientes"`
	rows, err := db.Query(selectSentence)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return
		}
		fmt.Println(id, name)
	}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
