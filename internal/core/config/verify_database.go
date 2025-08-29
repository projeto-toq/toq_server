package config

import (
	"log/slog"
	// "time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	// usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (c *config) VerifyDatabase() {
	//check if the database version is the same as the application version
	if getDBVersion(c) != globalmodel.AppVersion {
		//if not, update the database
		updateDB(c)
	}

	// //check if should populate database based on env parameter
	// if c.env.DATABASE.Populate {
	// 	slog.Info("Populating database (forced by configuration)")
	// 	populate(c)
	// }
}

func getDBVersion(c *config) string {
	config, err := c.globalService.GetConfiguration(c.context)
	if err != nil {
		slog.Error("error getting configuration", "error", err)
		panic(err)
	}
	return config["version"]
}

func updateDB(c *config) {
	panic("Database version is different from application version. Please update the database.")
}

// func populate(c *config) {
// 	slog.Info("Starting database population")

// 	// Criar roles base primeiro
// 	createBaseRoles(c)
// 	slog.Info("Base roles created successfully")

// 	// Criar usuário root
// 	createRootUser(c)
// 	slog.Info("Root user created successfully")

// 	slog.Info("Database population completed")
// }

// func createBaseRoles(c *config) {
// 	roles := []struct {
// 		role usermodel.UserRole
// 		name string
// 	}{
// 		{usermodel.RoleRoot, "Root"},
// 		{usermodel.RoleOwner, "Proprietário"},
// 		{usermodel.RoleRealtor, "Corretor"},
// 		{usermodel.RoleAgency, "Imobiliária"},
// 	}

// 	for _, r := range roles {
// 		err := c.userService.CreateBaseRole(c.context, r.role, r.name)
// 		if err != nil {
// 			slog.Error("error creating base role", "role", r.name, "error", err)
// 			// Continue mesmo com erro, pois pode ser que o role já exista
// 			slog.Warn("continuing despite base role creation error", "role", r.name)
// 		} else {
// 			slog.Info("base role created successfully", "role", r.name)
// 		}
// 	}
// }

// func createRootUser(c *config) {
// 	slog.Info("Creating root user")

// 	root := usermodel.NewUser()
// 	root.SetID(0)
// 	root.SetFullName("TOQ Root")
// 	root.SetNickName("Root")
// 	root.SetNationalID("52642435000133")
// 	born, err := time.Parse("2006-01-02", "2023-10-01")
// 	if err != nil {
// 		slog.Error("error parsing date on Root User creation", "error", err)
// 		panic(err)
// 	}
// 	root.SetBornAt(born)
// 	root.SetPhoneNumber("+551152413731")
// 	root.SetEmail("toq@toq.app.br")
// 	root.SetZipCode("06472-001")
// 	root.SetStreet("Av. Copacabana")
// 	root.SetNumber("268")
// 	root.SetComplement("Sala 2305")
// 	root.SetNeighborhood("Alphaville")
// 	root.SetCity("Barueri")
// 	root.SetState("SP")
// 	root.SetPassword("Senh@123")
// 	root.SetLastActivityAt(time.Now().UTC())
// 	root.SetDeleted(false)

// 	err = c.userService.CreateRoot(c.context, root)
// 	if err != nil {
// 		slog.Error("error creating root user", "error", err)
// 		panic(err)
// 	}

// 	slog.Info("root user created successfully", "email", "toq@toq.app.br")
// }
