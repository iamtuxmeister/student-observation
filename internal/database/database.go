package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	DB() *gorm.DB
}

type service struct {
	db *gorm.DB
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

// Authentication Midleware
type User struct {
	gorm.Model
	Id          uint `gorm:"primaryKey;"`
	Username    string
	Password    string
	Salt        string
	First       string
	Last        string
	Email       string
	Permissions []Permission `gorm:"many2many:user_permissions;"`
	Groups      []Group      `gorm:"many2many:user_groups;"`
	IsReviewer  bool
}

type Group struct {
	gorm.Model
	Id          uint `gorm:"primaryKey;"`
	Name        string
	Permissions []Permission `gorm:"many2many:group_permissions;"`
}

type Permission struct {
	gorm.Model
	Id   uint `gorm:"primaryKey;"`
	Name string
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Permission{},
		&User{},
		&Group{},
	)
	if err != nil {
		log.Fatal(err)
	}

	//create := Permission{Id: 1, Name: "create_record"}
	//db.Create(&create)
	//read := Permission{Id: 2, Name: "read_record"}
	//db.Create(&read)
	//update := Permission{Id: 3, Name: "update_record"}
	//db.Create(&update)
	//del := Permission{Id: 4, Name: "delete_record"}
	//db.Create(&del)
	//admin := Group{Id: 1, Name: "admin", Permissions: []Permission{create, read, update, del}}
	//db.Create(&admin)
	//user := User{Id: 1, Username: "admin", Password: "Ma51LLpr", First: "Admin", Last: "BigGuy", Email: "yourmom@admin.net", Groups: []Group{admin}, IsReviewer: true}
	//db.Create(&user)

}
func (s *service) DB() *gorm.DB {
	return s.db
}

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", username, password, host, port, dbname)
	db, err := gorm.Open(sqlite.Open("mydb.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	Migrate(db)
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Health() map[string]string {
	return map[string]string{
		"status": "ok",
	}
}
func (s *service) Close() error {
	db, _ := s.db.DB()
	return db.Close()
}
