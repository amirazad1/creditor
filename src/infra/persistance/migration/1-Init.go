package migration

import (
	"github.com/amirazad1/creditor/config"
	constants "github.com/amirazad1/creditor/constant"
	"github.com/amirazad1/creditor/domain/model"
	"github.com/amirazad1/creditor/infra/persistance/database"
	"github.com/amirazad1/creditor/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

func Up1() {
	database := database.GetDb()

	createTables(database)
	createDefaultUserInformation(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	// Base Information
	tables = addNewTable(database, model.City{}, tables)
	tables = addNewTable(database, model.Building{}, tables)
	tables = addNewTable(database, model.Project{}, tables)
	tables = addNewTable(database, model.CostTypeParent{}, tables)
	tables = addNewTable(database, model.CostType{}, tables)
	tables = addNewTable(database, model.CostSubType{}, tables)
	tables = addNewTable(database, model.Year{}, tables)
	tables = addNewTable(database, model.CostDescription{}, tables)
	tables = addNewTable(database, model.Plan{}, tables)
	tables = addNewTable(database, model.PlanProject{}, tables)
	tables = addNewTable(database, model.Addition{}, tables)

	// User
	tables = addNewTable(database, model.User{}, tables)
	tables = addNewTable(database, model.Role{}, tables)
	tables = addNewTable(database, model.UserRole{}, tables)


	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Error(logging.Postgres, logging.Migration, err.Error(), nil)
	}
	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}


func createDefaultUserInformation(database *gorm.DB) {

	adminRole := model.Role{Name: constants.AdminRoleName}
	createRoleIfNotExists(database, &adminRole)

	defaultRole := model.Role{Name: constants.DefaultRoleName}
	createRoleIfNotExists(database, &defaultRole)

	u := model.User{Username: constants.DefaultUserName, FirstName: "Test", LastName: "Test",
		MobileNumber: "09111112222", Email: "admin@admin.com"}
	pass := "12345678"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	createAdminUserIfNotExists(database, &u, adminRole.Id)

}

func createRoleIfNotExists(database *gorm.DB, r *model.Role) {
	exists := 0
	database.
		Model(&model.Role{}).
		Select("1").
		Where("name = ?", r.Name).
		First(&exists)
	if exists == 0 {
		database.Create(r)
	}
}

func createAdminUserIfNotExists(database *gorm.DB, u *model.User, roleId int) {
	exists := 0
	database.
		Model(&model.User{}).
		Select("1").
		Where("username = ?", u.Username).
		First(&exists)
	if exists == 0 {
		database.Create(u)
		ur := model.UserRole{UserId: u.Id, RoleId: roleId}
		database.Create(&ur)
	}
}
