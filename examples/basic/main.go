package main

import (
	"fmt"
	"log"

	basemodel "github.com/yourusername/mongo-basemodel"
)

// User struct that embeds BaseCollection
type User struct {
	basemodel.BaseCollection `bson:",inline"`
	Name                     string `json:"name" bson:"name"`
	Email                    string `json:"email" bson:"email"`
	Age                      int    `json:"age" bson:"age"`
}

// Product struct that embeds BaseCollection
type Product struct {
	basemodel.BaseCollection `bson:",inline"`
	Name                     string  `json:"name" bson:"name"`
	Price                    float64 `json:"price" bson:"price"`
	Description              string  `json:"description" bson:"description"`
	Category                 string  `json:"category" bson:"category"`
}

func main() {
	fmt.Println("=== MongoDB BaseModel Example ===\n")

	// Example 1: Creating a new user
	fmt.Println("1. Creating a new user:")
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	// Set metadata for insert
	user.SetInsertMeta()
	fmt.Printf("   User ID: %s\n", user.GetID())
	fmt.Printf("   Name: %s\n", user.Name)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Created At: %v\n", user.GetCreatedAt())
	fmt.Printf("   Is Deleted: %v\n\n", user.IsDeleted())

	// Example 2: Updating the user
	fmt.Println("2. Updating the user:")
	user.Name = "John Smith"
	user.Age = 31
	user.SetUpdateMeta()
	fmt.Printf("   Updated Name: %s\n", user.Name)
	fmt.Printf("   Updated Age: %d\n", user.Age)
	fmt.Printf("   Updated At: %v\n\n", user.GetUpdatedAt())

	// Example 3: Soft deleting the user
	fmt.Println("3. Soft deleting the user:")
	user.SetDeleteMeta()
	fmt.Printf("   Is Deleted: %v\n", user.IsDeleted())
	fmt.Printf("   Deleted At: %v\n\n", user.GetDeletedAt())

	// Example 4: Creating a product
	fmt.Println("4. Creating a new product:")
	product := &Product{
		Name:        "Gaming Laptop",
		Price:       45000.00,
		Description: "High-performance gaming laptop with RTX 4080",
		Category:    "Electronics",
	}

	product.SetInsertMeta()
	fmt.Printf("   Product ID: %s\n", product.GetID())
	fmt.Printf("   Name: %s\n", product.Name)
	fmt.Printf("   Price: %.2f THB\n", product.Price)
	fmt.Printf("   Created At: %v\n", product.GetCreatedAt())

	// Example 5: Multiple updates on product
	fmt.Println("\n5. Updating product multiple times:")

	// First update
	product.Price = 42000.00
	product.SetUpdateMeta()
	fmt.Printf("   First update - New price: %.2f THB at %v\n",
		product.Price, product.GetUpdatedAt())

	// Second update
	product.Description = "High-performance gaming laptop with RTX 4080 - On Sale!"
	product.SetUpdateMeta()
	fmt.Printf("   Second update - Updated description at %v\n", product.GetUpdatedAt())

	// Example 6: Demonstrating all getter methods
	fmt.Println("\n6. All getter methods:")
	fmt.Printf("   ID: %s\n", product.GetID())
	fmt.Printf("   Created: %v\n", product.GetCreatedAt())
	fmt.Printf("   Updated: %v\n", product.GetUpdatedAt())
	fmt.Printf("   Deleted: %v\n", product.GetDeletedAt())
	fmt.Printf("   Is Deleted: %v\n", product.IsDeleted())

	fmt.Println("\n=== Example completed successfully! ===")
}

// Helper function to demonstrate error handling
func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
