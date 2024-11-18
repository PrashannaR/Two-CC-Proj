package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/mod/sumdb/storage"
	"gorm.io/gorm"
)

//MARK: Repository
type Repository struct{
	DB *gorm.DB
}

//MARK: Routes
func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

//MARK: Book struct
type Book struct{
	Author			string		`json:"author"`
	Title			string		`json:"title"`
	Publisher		string		`json:"publisher"`
}

//MARK: CreateBook
func (r *Repository) CreateBook(c *fiber.Ctx) error{
	book := Book{}

	err := c.BodyParser(&book)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"an error occurred while creating the book"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"book created successfully", "data":book})

	return nil

}

//MARK: DeleteBook
func (r *Repository) DeleteBook (c *fiber.Ctx) error{

}

//MARK: GetBookById
func (r *Repository) GetBookById(c *fiber.Ctx) error{

}

//MARK: Get All Books
func (r *Repository) GetBooks(c *fiber.Ctx) error{

}

//MARK: Main
func main(){
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}