package controllers

import (
	"context"
	"github/GitKing-web/goSimpleBlogApp/config"
	"github/GitKing-web/goSimpleBlogApp/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func HandleSignup(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
		})
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	collection := config.Database.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already registered",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	user.Password = string(hashedPassword)

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func HandleLogin(c *fiber.Ctx) error {
	var loginData models.User
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
		})
	}

	validate := validator.New()
	err := validate.Struct(loginData)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	collection := config.Database.Collection("users")
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

func CreatePost(c *fiber.Ctx) error {
	var post models.Posts
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
		})
	}

	validate := validator.New()
	err := validate.Struct(post)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	collection := config.Database.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post.PostId = primitive.NewObjectID()
	post.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	_, err = collection.InsertOne(ctx, post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post created successfully",
		"post":    post,
	})
}

func GetPosts(c *fiber.Ctx) error {
	collection := config.Database.Collection("posts")
	var posts []models.Posts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving posts",
		})
	}
	err = post.All(ctx, &posts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving posts",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Posts retrieved successfully",
		"posts":   posts,
	})
}

func GetPost(c *fiber.Ctx) error {
	var post models.Posts
	id := c.Params("id")
	collection := config.Database.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}
	err = collection.FindOne(ctx, bson.M{"postId": objId}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post retrieved successfully",
		"post":    post,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := config.Database.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid post ID",
		})
	}
	var updatePost models.Posts
	if err := c.BodyParser(&updatePost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body",
		})
	}
	updatePost.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	update := bson.M{
		"$set": bson.M{
			"title":     updatePost.Title,
			"content":   updatePost.Content,
			"updatedAt": updatePost.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"postId": objId}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating post",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
	})
}

func DeletePost(c *fiber.Ctx) error {
	collection := config.Database.Collection("posts")
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post Id",
		})
	}

	_, err = collection.DeleteOne(ctx, bson.M{"postId": objId})
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error deleting post",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "post deleted successfully",
	})
}
