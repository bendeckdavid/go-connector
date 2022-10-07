package conn

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Start && Migrate Database
func StartDB(path string, modules ...[]any) {
	var err error

	DB, err = gorm.Open(sqlite.Open(path+".db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var migrations []any

	for _, md := range modules {
		migrations = append(migrations, md...)
	}

	// Migrate the schema
	DB.AutoMigrate(migrations...)
}
