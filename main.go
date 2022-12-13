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
)

var (
	PostgreUser = os.Getenv("Postgre-User")
	PostgrePass = os.Getenv("Postgre-Pass")
	PostgreHost = os.Getenv("Postgre-Host")
	PostgrePort = os.Getenv("Postgre-Port")
	PostgreDB   = os.Getenv("Postgre-DB")
)

// Postgre SQL Structure:
// create table organizaciones
//(
//    id     serial,
//    nombre VARCHAR(50) NOT NULL,
//    UipathName VARCHAR(50) NOT NULL, -- Cada organización tiene un nombre en UiPath
//    TenantName VARCHAR(50) NOT NULL, -- Cada organización tiene 1 tenant en UiPath
//    PRIMARY KEY (id)
//);
//
//CREATE TABLE Procesos
//(
//    id     serial,
//    nombre VARCHAR(50) NOT NULL,
//    FolderID int NOT NULL, -- Cada proceso pertenece a una carpeta, solo se guarda el ID de esta, ya que el resto es irrelevante
//    organizacion_id INTEGER NOT NULL,
//    warning_tolerance int NOT NULL default 10, -- Cantidad de warnings permitidos
//    error_tolerance int NOT NULL default 0, -- Cantidad de errores permitidos
//    fatal_tolerance int NOT NULL default 0, -- Cantidad de fatales permitidos
//    PRIMARY KEY (id),
//    FOREIGN KEY (organizacion_id) REFERENCES organizaciones(id)
//);
//
//create table usuarios
//(
//    id       serial,
//    nombre   VARCHAR(50) NOT NULL,
//    apellido VARCHAR(50) NOT NULL,
//    email    VARCHAR(50) NOT NULL,
//    password VARCHAR(50) NOT NULL,
//    PRIMARY KEY (id)
//);
//
//create table roles
//(
//    id     serial,
//    nombre VARCHAR(50) NOT NULL, -- Administrador, Desarrrollador
//    PRIMARY KEY (id)
//);
//
//create table clientes
//(
//    id              serial,
//    nombre          VARCHAR(50) NOT NULL,
//    apellido        VARCHAR(50) NOT NULL,
//    email           VARCHAR(50) NOT NULL,
//    organizacion_id INT         NOT NULL REFERENCES organizaciones (id),
//    PRIMARY KEY (id)
//);
//
//create table usuarios_roles
//(
//    id         serial, -- 1 Desarrollador, 2 Administrador
//    usuario_id int NOT NULL,
//    rol_id     int NOT NULL,
//    PRIMARY KEY (id),
//    FOREIGN KEY (usuario_id) REFERENCES usuarios (id),
//    FOREIGN KEY (rol_id) REFERENCES roles (id)
//);
//
//create table procesos_usuarios
//(
//    -- Cada usuario puede tener acceso a varios procesos, y cada proceso puede tener varios usuarios
//    -- Esto sirve para que varios desarrolladores puedan participar de la resolución de un proceso
//    -- o por ejemplo si un desarrollador no está disponible, otro pueda tomar su lugar y compartir la información
//    id         serial,
//    proceso_id int NOT NULL,
//    usuario_id int NOT NULL,
//    PRIMARY KEY (id),
//    FOREIGN KEY (proceso_id) REFERENCES procesos (id),
//    FOREIGN KEY (usuario_id) REFERENCES usuarios (id)
//);
//
//create table procesos_clientes
//(
//    -- Cada cliente puede tener acceso a varios procesos, y cada proceso puede tener varios clientes
//    -- Esto sirve para que se sepa donde comunicar los errores de un proceso, y el estado de los mismos
//    id         serial,
//    proceso_id int NOT NULL,
//    cliente_id int NOT NULL,
//    PRIMARY KEY (id),
//    FOREIGN KEY (proceso_id) REFERENCES procesos (id),
//    FOREIGN KEY (cliente_id) REFERENCES clientes (id)
//);
//
//create table incidentes_procesos
//(
//    id         serial,
//    proceso_id int NOT NULL,
//    incidente  text,
//    tipo       int not null default 1, -- 1: incidente, 2: mejora, 3: cambio, 4: solicitud, 5: requerimiento
//    estado     int not null default 1, -- 1: abierto, 2: en proceso, 3: cerrado
//    PRIMARY KEY (id),
//    FOREIGN KEY (proceso_id) REFERENCES procesos (id)
//);
//
//create table incidentes_detalle
//(
//    -- Los incidentes pueden tener varios detalles, por ejemplo, si un incidente es un error, se puede agregar el stacktrace
//    -- y dejar un historial de los cambios que se le hicieron, o cuanto tiempo se le dedicó en total
//    id           serial,
//    incidente_id int NOT NULL,
//    detalle      text,
//    fecha_inicio timestamp,
//    fecha_fin    timestamp,
//    PRIMARY KEY (id),
//    FOREIGN KEY (incidente_id) REFERENCES incidentes_procesos (id)
//);

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
	app.Get("/Incidentes", func(ctx *fiber.Ctx) error {
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", PostgreUser, PostgrePass, PostgreHost, PostgrePort, PostgreDB))
		defer db.Close()
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if ctx.Query("id") != "" {
			var incidente Structs.Incidentes_Procesos
			selectSentence := `Select id, proceso_id, incidente, tipo, estado from "incidentes_procesos" where "id"=$1`
			err = db.QueryRow(selectSentence, ctx.Query("id")).Scan(&incidente.ID, &incidente.ProcesoID, &incidente.Incidente, &incidente.Tipo, &incidente.Estado)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return ctx.JSON(incidente)
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
		insertSentence := `INSERT INTO "incidentes_procesos" ("proceso_id","incidente","tipo","estado") VALUES ($1,$2,$3,$4)`
		_, err = db.Exec(insertSentence, incidente.ProcesoID, incidente.Incidente, incidente.Tipo, incidente.Estado)
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
		return ctx.JSON(usuario)
	})
	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
