package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	basemodel "github.com/yourusername/mongo-basemodel"
)

// User struct with BaseCollection embedded
type User struct {
	basemodel.BaseCollection `bson:",inline"`
	Name                     string `json:"name" bson:"name"`
	Email                    string `json:"email" bson:"email"`
	Age                      int    `json:"age" bson:"age"`
	IsActive                 bool   `json:"is_active" bson:"is_active"`
}

// UserRepository handles database operations for users
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create inserts a new user
func (r *UserRepository) Create(user *User) error {
	user.SetInsertMeta()
	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}

// FindByID finds user by ID (excluding deleted)
func (r *UserRepository) FindByID(id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"deleted_at": bson.M{"$exists": false},
	}

	var user User
	err = r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindAll finds all active users (excluding deleted)
func (r *UserRepository) FindAll() ([]*User, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []*User
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// Update updates user information
func (r *UserRepository) Update(id string, user *User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user.SetUpdateMeta()
	filter := bson.M{
		"_id":        objID,
		"deleted_at": bson.M{"$exists": false},
	}

	update := bson.M{"$set": bson.M{
		"name":       user.Name,
		"email":      user.Email,
		"age":        user.Age,
		"is_active":  user.IsActive,
		"updated_at": user.UpdatedAt,
	}}

	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// SoftDelete performs soft delete on user
func (r *UserRepository) SoftDelete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"deleted_at": &now}}

	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// FindDeleted finds all soft-deleted users
func (r *UserRepository) FindDeleted() ([]*User, error) {
	filter := bson.M{"deleted_at": bson.M{"$exists": true}}

	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []*User
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func main() {
	fmt.Println("=== MongoDB BaseModel Integration Example ===\n")

	// Connect to MongoDB
	// Note: Change the connection string to match your MongoDB setup
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		log.Println("Please make sure MongoDB is running on localhost:27017")
		log.Println("You can also change the connection string in the code")
		return
	}
	defer client.Disconnect(context.TODO())

	// Test connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		log.Println("Please make sure MongoDB is running and accessible")
		return
	}

	fmt.Println("✓ Connected to MongoDB successfully!")

	db := client.Database("basemodel_example")
	userRepo := NewUserRepository(db)

	// Clean up collection for demo
	userRepo.collection.Drop(context.TODO())

	// Example 1: Create users
	fmt.Println("\n1. Creating users:")
	users := []*User{
		{
			Name:     "Alice Johnson",
			Email:    "alice@example.com",
			Age:      28,
			IsActive: true,
		},
		{
			Name:     "Bob Smith",
			Email:    "bob@example.com",
			Age:      35,
			IsActive: true,
		},
		{
			Name:     "Charlie Brown",
			Email:    "charlie@example.com",
			Age:      42,
			IsActive: false,
		},
	}

	var createdUserIDs []string
	for i, user := range users {
		err := userRepo.Create(user)
		if err != nil {
			log.Printf("Failed to create user %d: %v", i+1, err)
			continue
		}
		createdUserIDs = append(createdUserIDs, user.GetID())
		fmt.Printf("   ✓ Created user: %s (ID: %s)\n", user.Name, user.GetID())
	}

	// Example 2: Find all users
	fmt.Println("\n2. Finding all active users:")
	allUsers, err := userRepo.FindAll()
	if err != nil {
		log.Printf("Failed to find users: %v", err)
		return
	}

	for _, user := range allUsers {
		fmt.Printf("   - %s (%s) - Age: %d, Active: %v\n",
			user.Name, user.Email, user.Age, user.IsActive)
	}

	// Example 3: Find user by ID
	fmt.Println("\n3. Finding user by ID:")
	if len(createdUserIDs) > 0 {
		foundUser, err := userRepo.FindByID(createdUserIDs[0])
		if err != nil {
			log.Printf("Failed to find user: %v", err)
		} else {
			fmt.Printf("   Found: %s (Created: %v)\n",
				foundUser.Name, foundUser.GetCreatedAt())
		}
	}

	// Example 4: Update user
	fmt.Println("\n4. Updating user:")
	if len(createdUserIDs) > 0 {
		updateUser := &User{
			Name:     "Alice Johnson Updated",
			Email:    "alice.updated@example.com",
			Age:      29,
			IsActive: true,
		}

		err := userRepo.Update(createdUserIDs[0], updateUser)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
		} else {
			fmt.Printf("   ✓ Updated user successfully\n")

			// Verify update
			updatedUser, err := userRepo.FindByID(createdUserIDs[0])
			if err != nil {
				log.Printf("Failed to find updated user: %v", err)
			} else {
				fmt.Printf("   Updated info: %s (%s) - Updated at: %v\n",
					updatedUser.Name, updatedUser.Email, updatedUser.GetUpdatedAt())
			}
		}
	}

	// Example 5: Soft delete user
	fmt.Println("\n5. Soft deleting user:")
	if len(createdUserIDs) > 1 {
		err := userRepo.SoftDelete(createdUserIDs[1])
		if err != nil {
			log.Printf("Failed to soft delete user: %v", err)
		} else {
			fmt.Printf("   ✓ Soft deleted user successfully\n")
		}
	}

	// Example 6: Find active users after deletion
	fmt.Println("\n6. Active users after soft deletion:")
	activeUsers, err := userRepo.FindAll()
	if err != nil {
		log.Printf("Failed to find active users: %v", err)
	} else {
		fmt.Printf("   Active users count: %d\n", len(activeUsers))
		for _, user := range activeUsers {
			fmt.Printf("   - %s (%s)\n", user.Name, user.Email)
		}
	}

	// Example 7: Find deleted users
	fmt.Println("\n7. Soft deleted users:")
	deletedUsers, err := userRepo.FindDeleted()
	if err != nil {
		log.Printf("Failed to find deleted users: %v", err)
	} else {
		fmt.Printf("   Deleted users count: %d\n", len(deletedUsers))
		for _, user := range deletedUsers {
			fmt.Printf("   - %s (Deleted at: %v)\n", user.Name, user.GetDeletedAt())
		}
	}

	// Example 8: Advanced query with filters
	fmt.Println("\n8. Advanced query - Active users over 30:")
	filter := bson.M{
		"age":        bson.M{"$gt": 30},
		"is_active":  true,
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := userRepo.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to execute advanced query: %v", err)
	} else {
		defer cursor.Close(context.TODO())

		var filteredUsers []*User
		for cursor.Next(context.TODO()) {
			var user User
			if err := cursor.Decode(&user); err != nil {
				log.Printf("Failed to decode user: %v", err)
				continue
			}
			filteredUsers = append(filteredUsers, &user)
		}

		fmt.Printf("   Found %d active users over 30:\n", len(filteredUsers))
		for _, user := range filteredUsers {
			fmt.Printf("   - %s (%d years old)\n", user.Name, user.Age)
		}
	}

	fmt.Println("\n=== MongoDB integration example completed! ===")
}
