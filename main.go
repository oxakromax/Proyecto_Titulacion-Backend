package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/oxakromax/ProyectoTitulo-Monitoreo/Utils"
	"github.com/oxakromax/Proyecto_Titulacion-Backend/SQL/Structs"
	"os"
	"strings"
)

var (
	PostgreUser = os.Getenv("Postgre-User")
	PostgrePass = os.Getenv("Postgre-Pass")
	PostgreHost = os.Getenv("Postgre-Host")
	PostgrePort = os.Getenv("Postgre-Port")
	PostgreDB   = os.Getenv("Postgre-DB")
)

func PostProcesses(c *fiber.Ctx) error {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
	defer db.Close()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var process Structs.ProcesosBDD
	err = json.Unmarshal(c.Body(), &process)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	insertSentence := `INSERT INTO "procesos" ("nombre","folderid","organizacion_id","warning_tolerance","error_tolerance","fatal_tolerance") VALUES ($1,$2,$3,$4,$5,$6)`
	_, err = db.Exec(insertSentence, process.Nombre, process.Folderid, process.OrganizacionId, process.WarningTolerance, process.ErrorTolerance, process.FatalTolerance)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendString("OK")
}

func GetProcesses(c *fiber.Ctx) error {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
	defer db.Close()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if c.Query("id") != "" {
		var process Structs.ProcesosBDD
		selectSentence := `Select "id","nombre","folderid","organizacion_id","warning_tolerance","error_tolerance","fatal_tolerance" from "procesos" where "id"=$1`
		err = db.QueryRow(selectSentence, c.Query("id")).Scan(&process.ID, &process.Nombre, &process.Folderid, &process.OrganizacionId, &process.WarningTolerance, &process.ErrorTolerance, &process.FatalTolerance)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(process)
	}
	ProcessArray := new(Structs.ProcessBDDArray)
	selectSentence := `Select "id","nombre","folderid","organizacion_id","warning_tolerance","error_tolerance","fatal_tolerance" from "procesos"`
	rows, err := db.Query(selectSentence)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer rows.Close()
	for rows.Next() {
		var process Structs.ProcesosBDD
		err = rows.Scan(&process.ID, &process.Nombre, &process.Folderid, &process.OrganizacionId, &process.WarningTolerance, &process.ErrorTolerance, &process.FatalTolerance)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		ProcessArray.Processes = append(ProcessArray.Processes, process)
	}
	return c.JSON(ProcessArray)
}

func main() {
	var err error
	Utils.QueryAuth.ClientId = os.Getenv("APP-ID")
	Utils.QueryAuth.ClientSecret = os.Getenv("APP-Secret")
	Utils.QueryAuth.Scope = os.Getenv("APP-Scope")
	Utils.UipathOrg.Id = os.Getenv("Orchestrator-ID")
	if err != nil {
		return
	}
	app := fiber.New()
	app.Get("/Processes", GetProcesses)
	app.Post("/Processes", PostProcesses)
	app.Get("/Clientes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if ctx.Query("id") != "" {
			var client Structs.Clientes
			selectSentence := `Select id, nombre, apellido, email, organizacion_id from "clientes" where "id"=$1`
			err = db.QueryRow(selectSentence, ctx.Query("id")).Scan(&client.ID, &client.Nombre, &client.Apellido, &client.Email, &client.OrganizacionID)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return ctx.JSON(client)
		}
		ClientArray := new(Structs.ClientesArray)
		selectSentence := `Select id, nombre, apellido, email, organizacion_id from "clientes"`
		rows, err := db.Query(selectSentence)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var client Structs.Clientes
			err = rows.Scan(&client.ID, &client.Nombre, &client.Apellido, &client.Email, &client.OrganizacionID)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			ClientArray.Clientes = append(ClientArray.Clientes, client)
		}
		return ctx.JSON(ClientArray)
	})
	app.Post("/Clientes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var client Structs.Clientes
		err = json.Unmarshal(ctx.Body(), &client)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		insertSentence := `INSERT INTO "clientes" ("nombre","apellido","email","organizacion_id") VALUES ($1,$2,$3,$4)`
		_, err = db.Exec(insertSentence, client.Nombre, client.Apellido, client.Email, client.OrganizacionID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("OK")
	})
	app.Get("/Incidentes_Detalle", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if ctx.Query("id") != "" {
			var incidenteDetalle Structs.Incidentes_Detalle
			selectSentence := `Select id, incidente_id, detalle, fecha_inicio, fecha_fin from "incidentes_detalle" where "id"=$1`
			err = db.QueryRow(selectSentence, ctx.Query("id")).Scan(&incidenteDetalle.ID, &incidenteDetalle.IncidenteID, &incidenteDetalle.Detalle, &incidenteDetalle.Fecha_Inicio, &incidenteDetalle.Fecha_Fin)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return ctx.JSON(incidenteDetalle)
		}
		if ctx.Query("incidente_id") != "" {
			var incidenteDetalle Structs.Incidentes_Detalle
			selectSentence := `Select id, incidente_id, detalle, fecha_inicio, fecha_fin from "incidentes_detalle" where "incidente_id"=$1`
			err = db.QueryRow(selectSentence, ctx.Query("incidente_id")).Scan(&incidenteDetalle.ID, &incidenteDetalle.IncidenteID, &incidenteDetalle.Detalle, &incidenteDetalle.Fecha_Inicio, &incidenteDetalle.Fecha_Fin)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return ctx.JSON(incidenteDetalle)
		}
		Incidentes_DetalleArray := new(Structs.Incidentes_DetalleArray)
		selectSentence := `Select id, incidente_id, detalle, fecha_inicio, fecha_fin from "incidentes_detalle"`
		rows, err := db.Query(selectSentence)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var incidenteDetalle Structs.Incidentes_Detalle
			err = rows.Scan(&incidenteDetalle.ID, &incidenteDetalle.IncidenteID, &incidenteDetalle.Detalle, &incidenteDetalle.Fecha_Inicio, &incidenteDetalle.Fecha_Fin)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			Incidentes_DetalleArray.Incidentes_Detalle = append(Incidentes_DetalleArray.Incidentes_Detalle, incidenteDetalle)
		}
		return ctx.JSON(Incidentes_DetalleArray)
	})
	app.Post("/Incidentes_Detalle", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var incidenteDetalle Structs.Incidentes_Detalle
		err = json.Unmarshal(ctx.Body(), &incidenteDetalle)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		insertSentence := `INSERT INTO "incidentes_detalle" ("incidente_id","detalle","fecha_inicio","fecha_fin") VALUES ($1,$2,$3,$4)`
		_, err = db.Exec(insertSentence, incidenteDetalle.IncidenteID, incidenteDetalle.Detalle, incidenteDetalle.Fecha_Inicio, incidenteDetalle.Fecha_Fin)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("OK")
	})
	app.Post("/Procesos_Clientes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var procesoCliente Structs.Procesos_Clientes
		err = ctx.BodyParser(&procesoCliente)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		insertSentence := `INSERT INTO "procesos_clientes" ("proceso_id","cliente_id") VALUES ($1,$2)`
		_, err = db.Exec(insertSentence, procesoCliente.ProcesoID, procesoCliente.ClienteID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("OK")
	})
	app.Delete("/Procesos_Clientes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var procesoCliente Structs.Procesos_Clientes
		err = ctx.BodyParser(&procesoCliente)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		deleteSentence := `DELETE FROM "procesos_clientes" WHERE "proceso_id"=$1 AND "cliente_id"=$2`
		_, err = db.Exec(deleteSentence, procesoCliente.ProcesoID, procesoCliente.ClienteID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("OK")
	})
	app.Post("/Procesos_Usuarios", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var procesoUsuario Structs.Procesos_Usuarios
		err = ctx.BodyParser(&procesoUsuario)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		insertSentence := `INSERT INTO "procesos_usuarios" ("proceso_id","usuario_id") VALUES ($1,$2)`
		_, err = db.Exec(insertSentence, procesoUsuario.ProcesoID, procesoUsuario.UsuarioID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("OK")
	})
	app.Delete("/Procesos_Usuarios", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var procesoUsuario Structs.Procesos_Usuarios
		err = ctx.BodyParser(&procesoUsuario)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		deleteSentence := `DELETE FROM "procesos_usuarios" WHERE "proceso_id"=$1 AND "usuario_id"=$2`
		_, err = db.Exec(deleteSentence, procesoUsuario.ProcesoID, procesoUsuario.UsuarioID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("OK")
	})
	app.Get("/Incidentes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if ctx.Query("id") != "null" && ctx.Query("id") != "" {
			var incidente Structs.Incidentes_Procesos
			selectSentence := `Select id, proceso_id, incidente, tipo, estado from "incidentes_procesos" where "id"=$1`
			err = db.QueryRow(selectSentence, ctx.Query("id")).Scan(&incidente.ID, &incidente.ProcesoID, &incidente.Incidente, &incidente.Tipo, &incidente.Estado)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			var Detalles []Structs.Incidentes_Detalle
			selectSentence = `Select id, incidente_id, detalle, fecha_inicio, fecha_fin from "incidentes_detalle" where "incidente_id"=$1`
			rows, err := db.Query(selectSentence, ctx.Query("id"))
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			defer rows.Close()
			for rows.Next() {
				var detalle Structs.Incidentes_Detalle
				err = rows.Scan(&detalle.ID, &detalle.IncidenteID, &detalle.Detalle, &detalle.Fecha_Inicio, &detalle.Fecha_Fin)
				if err != nil {
					return ctx.Status(500).SendString(err.Error())
				}
				Detalles = append(Detalles, detalle)
			}
			incidente.Detalles = Detalles
			marshal, err := json.Marshal(incidente)
			if err != nil {
				return err
			}
			// Add name of process
			name := ""
			selectSentence = `Select nombre from "procesos" where "id"=$1`
			err = db.QueryRow(selectSentence, incidente.ProcesoID).Scan(&name)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			// Solo al final
			marshal = []byte(strings.Replace(string(marshal), "]}", "],\"nombre_proceso\":\""+name+"\"}", 1))
			return ctx.SendString(string(marshal))
			//return ctx.JSON(incidente)
		}
		if ctx.Query("id_usuario") != "" && ctx.Query("id_usuario") != "null" {
			IncidentesArray := new(Structs.Incidentes_ProcesosArray)
			selectSentence := `Select ip.id, ip.proceso_id, incidente, tipo, estado from incidentes_procesos as ip join procesos_usuarios as pu on ip.proceso_id=pu.proceso_id where pu.usuario_id=$1`
			rows, err := db.Query(selectSentence, ctx.Query("id_usuario"))
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			defer rows.Close()
			for rows.Next() {
				var incidente Structs.Incidentes_Procesos
				err = rows.Scan(&incidente.ID, &incidente.ProcesoID, &incidente.Incidente, &incidente.Tipo, &incidente.Estado)
				if err != nil {
					return ctx.Status(500).SendString(err.Error())
				}
				var Detalles []Structs.Incidentes_Detalle
				selectSentence = `Select id, incidente_id, detalle, fecha_inicio, fecha_fin from "incidentes_detalle" where "incidente_id"=$1`
				rows2, err := db.Query(selectSentence, incidente.ID)
				if err != nil {
					return ctx.Status(500).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				defer rows2.Close()
				for rows2.Next() {
					var detalle Structs.Incidentes_Detalle
					err = rows2.Scan(&detalle.ID, &detalle.IncidenteID, &detalle.Detalle, &detalle.Fecha_Inicio, &detalle.Fecha_Fin)
					if err != nil {
						return ctx.Status(500).SendString(err.Error())
					}
					Detalles = append(Detalles, detalle)
				}
				incidente.Detalles = Detalles
				IncidentesArray.Incidentes_Procesos = append(IncidentesArray.Incidentes_Procesos, incidente)
			}
			return ctx.JSON(IncidentesArray)
		}
		Incidentes_ProcesosArray := new(Structs.Incidentes_ProcesosArray)
		selectSentence := `Select id, proceso_id, incidente, tipo, estado from "incidentes_procesos"`
		rows, err := db.Query(selectSentence)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var incidente Structs.Incidentes_Procesos
			err = rows.Scan(&incidente.ID, &incidente.ProcesoID, &incidente.Incidente, &incidente.Tipo, &incidente.Estado)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			var Detalles []Structs.Incidentes_Detalle
			selectSentence = `Select id, incidente_id, detalle, fecha_inicio, fecha_fin from "incidentes_detalle" where "incidente_id"=$1`
			rows2, err := db.Query(selectSentence, incidente.ID)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			defer rows2.Close()
			for rows2.Next() {
				var detalle Structs.Incidentes_Detalle
				err = rows2.Scan(&detalle.ID, &detalle.IncidenteID, &detalle.Detalle, &detalle.Fecha_Inicio, &detalle.Fecha_Fin)
				if err != nil {
					return ctx.Status(500).SendString(err.Error())
				}
				Detalles = append(Detalles, detalle)
			}
			incidente.Detalles = Detalles
			Incidentes_ProcesosArray.Incidentes_Procesos = append(Incidentes_ProcesosArray.Incidentes_Procesos, incidente)
		}
		return ctx.JSON(Incidentes_ProcesosArray)
	})
	app.Post("/Incidentes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var incidente Structs.Incidentes_Procesos
		err = ctx.BodyParser(&incidente)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if incidente.Detalles == nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": "No hay detalles, es de caracter obligatorio",
			})
		}
		insertSentence := `INSERT INTO "incidentes_procesos" ("proceso_id","incidente","tipo","estado") VALUES ($1,$2,$3,$4)`
		_, err = db.Exec(insertSentence, incidente.ProcesoID, incidente.Incidente, incidente.Tipo, incidente.Estado)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		idIncidente := 0
		selectSentence := `Select id from "incidentes_procesos" where "proceso_id"=$1 and "incidente"=$2 and "tipo"=$3 and "estado"=$4`
		err = db.QueryRow(selectSentence, incidente.ProcesoID, incidente.Incidente, incidente.Tipo, incidente.Estado).Scan(&idIncidente)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		for _, detalle := range incidente.Detalles {
			insertSentence = `INSERT INTO "incidentes_detalle" ("incidente_id","detalle","fecha_inicio","fecha_fin") VALUES ($1,$2,$3,$4)`
			_, err = db.Exec(insertSentence, idIncidente, detalle.Detalle, detalle.Fecha_Inicio, detalle.Fecha_Fin)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
		var emails []string
		ProcessName := ""
		selectSentence = `Select email,p.nombre from "usuarios" as users join procesos_usuarios pu on users.id = pu.usuario_id join procesos p on p.id = pu.proceso_id where p.id = $1`
		rows, err := db.Query(selectSentence, incidente.ProcesoID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var email string
			err = rows.Scan(&email, &ProcessName)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			emails = append(emails, email)
		}
		selectSentence = `select email from clientes join procesos_clientes pc on clientes.id = pc.cliente_id where pc.proceso_id = $1`
		rows, err = db.Query(selectSentence, incidente.ProcesoID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var email string
			err = rows.Scan(&email)
			if err != nil {
				return ctx.Status(500).SendString(err.Error())
			}
			emails = append(emails, email)
		}
		// -- 1: incidente, 2: mejora, 3: cambio, 4: solicitud, 5: requerimiento
		var tipo string
		switch incidente.Tipo {
		case 1:
			tipo = "Incidente"
		case 2:
			tipo = "Mejora"
		case 3:
			tipo = "Cambio"
		case 4:
			tipo = "Solicitud"
		case 5:
			tipo = "Requerimiento"
		}
		msgbody := fmt.Sprintf("Se ha registrado un nuevo %s en el proceso %s. Por favor ingrese a la plataforma para ver los detalles.", tipo, ProcessName)
		subject := fmt.Sprintf("Nuevo %s en el proceso %s, ID#%d", tipo, ProcessName, idIncidente)
		err = SendMail(msgbody, subject, emails)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(fiber.Map{
			"message": "Incidente creado",
			"id":      idIncidente,
		})
	})
	app.Post("/AuthUsuario", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var authUsuario Structs.AuthUsuario
		err = ctx.BodyParser(&authUsuario)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var usuario Structs.Usuario
		selectSentence := `SELECT "id","nombre","apellido","email" FROM "usuarios" WHERE "email"=$1 AND "password"=$2`
		err = db.QueryRow(selectSentence, authUsuario.Email, authUsuario.Password).Scan(&usuario.ID, &usuario.Nombre, &usuario.Apellido, &usuario.Email)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		selectSentence = `SELECT "id","nombre" FROM "roles" WHERE "id" IN (SELECT "rol_id" FROM "usuarios_roles" WHERE "usuario_id"=$1)`
		rows, err := db.Query(selectSentence, usuario.ID)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var rol struct {
				ID     int    `json:"ID"`
				Nombre string `json:"Nombre"`
			}
			err = rows.Scan(&rol.ID, &rol.Nombre)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			usuario.Roles = append(usuario.Roles, rol)
		}
		return ctx.Status(200).JSON(usuario)
	})
	app.Get("/Usuarios", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var usuarios []Structs.Usuario
		selectSentence := `SELECT "id","nombre","apellido","email" FROM "usuarios"`
		rows, err := db.Query(selectSentence)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer rows.Close()
		for rows.Next() {
			var usuario Structs.Usuario
			err = rows.Scan(&usuario.ID, &usuario.Nombre, &usuario.Apellido, &usuario.Email)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			selectSentence = `SELECT "id","nombre" FROM "roles" WHERE "id" IN (SELECT "rol_id" FROM "usuarios_roles" WHERE "usuario_id"=$1)`
			rows2, err := db.Query(selectSentence, usuario.ID)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			defer rows2.Close()
			for rows2.Next() {
				var rol struct {
					ID     int    `json:"ID"`
					Nombre string `json:"Nombre"`
				}
				err = rows2.Scan(&rol.ID, &rol.Nombre)
				if err != nil {
					return ctx.Status(500).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				usuario.Roles = append(usuario.Roles, rol)
			}
			usuarios = append(usuarios, usuario)
		}
		return ctx.Status(200).JSON(usuarios)
	})
	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
