package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"

	"github.com/graphql-go/graphql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Book struct {
	ID         uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	UserID     uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	Title      string    `db:"title" json:"title" validate:"required,lte=255"`
	Author     string    `db:"author" json:"author" validate:"required,lte=255"`
	BookStatus int       `db:"book_status" json:"book_status" validate:"required,len=1"`
}

func Query(queryString string) *graphql.Result {
	fields := graphql.Fields{
		"test": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "test successful", nil
			},
		},
	}

	objectConfig := graphql.ObjectConfig{Name: "haexr_schema_object", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(objectConfig)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err == nil {
		query := queryString
		params := graphql.Params{Schema: schema, RequestString: query}
		response := graphql.Do(params)
		if !response.HasErrors() {
			return response
		}
	}
	return nil
}

func main() {

	dsn := "host=localhost user=postgres password=2000 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Printf(db.Name())
	}

	db.AutoMigrate(&Book{})
	db.AutoMigrate(&User{})

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/query", func(c *fiber.Ctx) error {
		return c.JSON(Query("{test}"))
	})

	log.Fatal(app.Listen(":4000"))
}
