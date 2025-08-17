package store
import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"gorm.io/gorm"
	"sync"
	"log"
	"os"
	"time"
)

var (
	db *gorm.DB
	once sync.Once
)

var newLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 		   200 * time.Millisecond,
		LogLevel:             logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:             true,
})


func NewDB(dsn string) *gorm.DB {
	once.Do(func(){
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Fatalf("failed to connect to the database: %v", err)
		}
		log.Println("Connected to the database successfully")
	})
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	return db
}