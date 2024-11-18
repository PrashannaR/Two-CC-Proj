package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PrashannaR/two-cc-project/models"
	"github.com/PrashannaR/two-cc-project/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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
	bookModel := models.Books{}

	id := c.Params("id")

	if id == ""{
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not delete book"})
		return err.Error
	}
	
	c.Status(http.StatusOK).JSON(&fiber.Map{"message":"book deleted successfully"})

	return nil
}

//MARK: GetBookById
func (r *Repository) GetBookById(c *fiber.Ctx) error{
	bookModel := models.Books{}
	id := c.Params("id")

	if id == ""{
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message":"id cannot be empty"})
		return nil
	}

	err := r.DB.First(&bookModel, id)

	if err.Error != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not find book"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book found successfully", "data":bookModel})

	return nil
}

//MARK: Get All Books
func (r *Repository) GetBooks(c *fiber.Ctx) error{
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"an error occurred while fetching the books"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"books fetched successfully", "data":bookModels})
	return nil
}

//MARK: Main
func main(){
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:		os.Getenv("DB_HOST"),
		Port:		os.Getenv("DB_PORT"),
		Password:	os.Getenv("DB_PASSWORD"),
		User:		os.Getenv("DB_USER"),
		SSLMode:	os.Getenv("DB_SSLMODE"),
		DBName:		os.Getenv("DB_DBNAME"),
	}
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("Error migrating database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}