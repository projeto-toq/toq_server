package config

import (
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *config) VerifyDatabase() {
	//check id the database version is the same as the application version
	if getDBVersion(c) != globalmodel.AppVersion {
		//if not, update the database
		updateDB(c)
	}
	//check if the database is new
	if isEmpty(c) {
		//populate the database
		populate(c)
	}
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

func isEmpty(c *config) bool {
	_, err := c.userService.GetUsers(c.context)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return true
		}
		slog.Error("error getting users", "error", err)
		panic(err)
	}

	return false
}

func populate(c *config) {
	slog.Info("Populating database")
	createBaseRoles(c)
	createRootUser(c)
}

func createBaseRoles(c *config) {
	c.userService.CreateBaseRole(c.context, usermodel.RoleRoot, "Root")
	c.userService.CreateBaseRole(c.context, usermodel.RoleOwner, "Proprietário")
	c.userService.CreateBaseRole(c.context, usermodel.RoleRealtor, "Corretor")
	c.userService.CreateBaseRole(c.context, usermodel.RoleAgency, "Imobiliária")
}

func createRootUser(c *config) {
	root := usermodel.NewUser()
	root.SetID(0)
	root.SetFullName("TOQ Root")
	root.SetNickName("Root")
	root.SetNationalID("52642435000133")
	born, err := time.Parse("2006-01-02", "2023-10-01")
	if err != nil {
		slog.Error("error parsing date on Root User creation", "error", err)
		panic(err)
	}
	root.SetBornAt(born)
	root.SetPhoneNumber("+551152413731")
	root.SetEmail("toq@toq.app.br")
	root.SetZipCode("06472-001")
	root.SetStreet("Av. Copacabana")
	root.SetNumber("268")
	root.SetComplement("Sala 2305")
	root.SetNeighborhood("Alphaville")
	root.SetCity("Barueri")
	root.SetState("SP")
	root.SetPassword("Senh@123")
	root.SetLastActivityAt(time.Now().UTC())
	root.SetDeleted(false)
	err = c.userService.CreateRoot(c.context, root)
	if err != nil {
		slog.Error("error creating root user", "error", err)
		panic(err)
	}
}
