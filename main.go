package main

import (
	"compute-service-app/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	initDB()
	initialMigration()
}

type Config struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
}

func initDB() {
	config := Config{
		DB_USERNAME: util.GetConfig("DB_USERNAME"),
		DB_PASSWORD: util.GetConfig("DB_PASSWORD"),
		DB_PORT:     util.GetConfig("DB_PORT"),
		DB_HOST:     util.GetConfig("DB_HOST"),
		DB_NAME:     util.GetConfig("DB_NAME"),
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_USERNAME,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

type User struct {
	ID       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func initialMigration() {
	DB.AutoMigrate(&User{})
}

func GetUsersController(c echo.Context) error {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messages": "success get all users",
		"users":    users,
	})
}

func GetUserController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var user User
	if err := DB.First(&user, "id = ?", id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messages": "success get user",
		"user":     user,
	})
}

func DeleteUserController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var user User
	if err := DB.Delete(&user, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messages": fmt.Sprintf("success delete user with id = %d", id),
	})

}

func UpdateUserController(c echo.Context) error {
	var userEdit User
	var user User
	c.Bind(&userEdit)
	id, _ := strconv.Atoi(c.Param("id"))

	if err := DB.First(&user, "id = ?", id).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user.Email = userEdit.Email
	user.Password = userEdit.Password
	user.Name = userEdit.Name

	if err := DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messages": "success update user",
		"user":     user,
	})
}

func CreateUserController(c echo.Context) error {
	var user User
	c.Bind(&user)

	if err := DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"messages": "success create user",
		"user":     user,
	})
}

func main() {
	e := echo.New()
	e.GET("/users", GetUsersController)
	e.GET("/users/:id", GetUserController)
	e.POST("/users", CreateUserController)
	e.PUT("/users/:id", UpdateUserController)
	e.DELETE("/users/:id", DeleteUserController)

	port := "8000"

	e.Logger.Fatal(e.Start(":" + port))
}
