package migration

import (
	"github.com/amirazad1/creditor/config"
	"github.com/amirazad1/creditor/domain/models"
	"github.com/amirazad1/creditor/infra/persistance/database"
	"github.com/amirazad1/creditor/pkg/logging"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

func Up1() {
	database := database.GetDb()

	createTables(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	// Base Information
	tables = addNewTable(database, models.City{}, tables)
	tables = addNewTable(database, models.Building{}, tables)
	tables = addNewTable(database, models.Project{}, tables)
	tables = addNewTable(database, models.CostTypeParent{}, tables)
	tables = addNewTable(database, models.CostType{}, tables)
	tables = addNewTable(database, models.CostSubType{}, tables)
	tables = addNewTable(database, models.Year{}, tables)
	tables = addNewTable(database, models.CostDescription{}, tables)
	tables = addNewTable(database, models.Plan{}, tables)
	tables = addNewTable(database, models.PlanProject{}, tables)
	tables = addNewTable(database, models.Addition{}, tables)

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
