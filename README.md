# MongoDB Base Model Library

Go library ที่ให้ base model สำหรับ MongoDB collections พร้อมกับ common fields และ methods สำหรับการจัดการ metadata

## Features

- **BaseCollection struct** พร้อม common fields:
  - `_id` (ObjectID)
  - `created_at` (timestamp)
  - `updated_at` (timestamp)
  - `deleted_at` (timestamp สำหรับ soft delete)

- **Built-in methods**:
  - `SetInsertMeta()` - ตั้งค่า metadata สำหรับการ insert
  - `SetUpdateMeta()` - ตั้งค่า metadata สำหรับการ update
  - `SetDeleteMeta()` - ตั้งค่า metadata สำหรับการ soft delete
  - `IsDeleted()` - ตรวจสอบว่าถูก soft delete หรือไม่
  - Getter methods สำหรับ fields ต่างๆ

## Installation

```bash
go get github.com/thitipong-pu/mongo-basemodel
```

## Usage

### 1. Basic Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/thitipong-pu/mongo-basemodel"
)

// สร้าง struct ของคุณโดย embed BaseCollection
type User struct {
    basemodel.BaseCollection `bson:",inline"`
    Name  string `json:"name" bson:"name"`
    Email string `json:"email" bson:"email"`
}

func main() {
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // เมื่อ insert ใหม่
    user.SetInsertMeta()
    fmt.Printf("User ID: %s\n", user.GetID())
    fmt.Printf("Created At: %v\n", user.GetCreatedAt())
    
    // เมื่อ update
    user.Name = "John Smith"
    user.SetUpdateMeta()
    fmt.Printf("Updated At: %v\n", user.GetUpdatedAt())
    
    // เมื่อ soft delete
    user.SetDeleteMeta()
    fmt.Printf("Deleted: %v\n", user.IsDeleted())
}
```

### 2. Integration with MongoDB Driver

```go
package main

import (
    "context"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    
    "github.com/thitipong-pu/mongo-basemodel"
)

type Product struct {
    basemodel.BaseCollection `bson:",inline"`
    Name        string  `json:"name" bson:"name"`
    Price       float64 `json:"price" bson:"price"`
    Description string  `json:"description" bson:"description"`
}

func main() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.TODO())
    
    db := client.Database("testdb")
    collection := db.Collection("products")
    
    // Create new product
    product := &Product{
        Name:        "Laptop",
        Price:       25000.00,
        Description: "Gaming Laptop",
    }
    
    // Set insert metadata
    product.SetInsertMeta()
    
    // Insert to MongoDB
    result, err := collection.InsertOne(context.TODO(), product)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Inserted product with ID: %v", result.InsertedID)
    
    // Update product
    product.Price = 23000.00
    product.SetUpdateMeta()
    
    filter := bson.M{"_id": product.Oid}
    update := bson.M{"$set": bson.M{
        "price":      product.Price,
        "updated_at": product.UpdatedAt,
    }}
    
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Product updated successfully")
    
    // Soft delete
    product.SetDeleteMeta()
    
    deleteUpdate := bson.M{"$set": bson.M{
        "deleted_at": product.DeletedAt,
    }}
    
    _, err = collection.UpdateOne(context.TODO(), filter, deleteUpdate)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Product soft deleted")
}
```

### 3. Query with Soft Delete Filter

```go
// ค้นหาเฉพาะ records ที่ไม่ถูก soft delete
func FindActiveProducts(collection *mongo.Collection) ([]*Product, error) {
    filter := bson.M{"deleted_at": bson.M{"$exists": false}}
    
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())
    
    var products []*Product
    for cursor.Next(context.TODO()) {
        var product Product
        if err := cursor.Decode(&product); err != nil {
            return nil, err
        }
        products = append(products, &product)
    }
    
    return products, nil
}

// ค้นหาเฉพาะ records ที่ถูก soft delete
func FindDeletedProducts(collection *mongo.Collection) ([]*Product, error) {
    filter := bson.M{"deleted_at": bson.M{"$exists": true}}
    
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())
    
    var products []*Product
    for cursor.Next(context.TODO()) {
        var product Product
        if err := cursor.Decode(&product); err != nil {
            return nil, err
        }
        products = append(products, &product)
    }
    
    return products, nil
}
```

### 4. Repository Pattern Example

```go
type ProductRepository struct {
    collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
    return &ProductRepository{
        collection: db.Collection("products"),
    }
}

func (r *ProductRepository) Create(product *Product) error {
    product.SetInsertMeta()
    _, err := r.collection.InsertOne(context.TODO(), product)
    return err
}

func (r *ProductRepository) Update(id string, product *Product) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    
    product.SetUpdateMeta()
    filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
    update := bson.M{"$set": product}
    
    _, err = r.collection.UpdateOne(context.TODO(), filter, update)
    return err
}

func (r *ProductRepository) SoftDelete(id string) error {
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

func (r *ProductRepository) FindByID(id string) (*Product, error) {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    
    filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
    
    var product Product
    err = r.collection.FindOne(context.TODO(), filter).Decode(&product)
    if err != nil {
        return nil, err
    }
    
    return &product, nil
}
```

## API Reference

### BaseCollection Fields

- `Oid` - MongoDB ObjectID
- `CreatedAt` - วันที่สร้าง record
- `UpdatedAt` - วันที่อัพเดท record ล่าสุด (nullable)
- `DeletedAt` - วันที่ soft delete (nullable)

### Methods

#### SetInsertMeta()
ตั้งค่า metadata สำหรับการ insert ใหม่:
- สร้าง ObjectID ใหม่
- ตั้งค่า CreatedAt เป็นเวลาปัจจุบัน

#### SetUpdateMeta()
ตั้งค่า metadata สำหรับการ update:
- ตั้งค่า UpdatedAt เป็นเวลาปัจจุบัน

#### SetDeleteMeta()
ตั้งค่า metadata สำหรับการ soft delete:
- ตั้งค่า DeletedAt เป็นเวลาปัจจุบัน

#### IsDeleted() bool
ตรวจสอบว่า record ถูก soft delete หรือไม่

#### GetID() string
ส่งคืน ObjectID เป็น string

#### GetCreatedAt() time.Time
ส่งคืน timestamp การสร้าง

#### GetUpdatedAt() *time.Time
ส่งคืน timestamp การอัพเดทล่าสุด

#### GetDeletedAt() *time.Time
ส่งคืน timestamp การ soft delete

## Best Practices

1. **ใช้ inline embedding** เมื่อสร้าง struct ของคุณ
2. **เรียก SetInsertMeta()** ก่อนการ insert
3. **เรียก SetUpdateMeta()** ก่อนการ update
4. **ใช้ soft delete** แทนการลบจริงด้วย SetDeleteMeta()
5. **เพิ่ม filter สำหรับ deleted_at** ในการ query เพื่อแยก active records

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Credits

This project was developed by:
- **Mr. Thitipong Punprom**
- **Mr. Kittipong Kraisit**

## License

This project is licensed under the MIT License - see the LICENSE file for details. 