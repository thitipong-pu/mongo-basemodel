package basemodel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseCollection provides common fields and methods for MongoDB collections
type BaseCollection struct {
	Oid       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

// SetInsertMeta sets the metadata for insert operations
// It generates a new ObjectID and sets the CreatedAt timestamp
func (b *BaseCollection) SetInsertMeta() {
	now := time.Now()
	b.Oid = primitive.NewObjectID()
	b.CreatedAt = now
}

// SetUpdateMeta sets the metadata for update operations
// It sets the UpdatedAt timestamp
func (b *BaseCollection) SetUpdateMeta() {
	now := time.Now()
	b.UpdatedAt = &now
}

// SetDeleteMeta sets the metadata for soft delete operations
// It sets the DeletedAt timestamp for soft deletion
func (b *BaseCollection) SetDeleteMeta() {
	now := time.Now()
	b.DeletedAt = &now
}

// IsDeleted checks if the record is soft deleted
func (b *BaseCollection) IsDeleted() bool {
	return b.DeletedAt != nil
}

// GetID returns the ObjectID as a string
func (b *BaseCollection) GetID() string {
	return b.Oid.Hex()
}

// GetCreatedAt returns the creation timestamp
func (b *BaseCollection) GetCreatedAt() time.Time {
	return b.CreatedAt
}

// GetUpdatedAt returns the last update timestamp
func (b *BaseCollection) GetUpdatedAt() *time.Time {
	return b.UpdatedAt
}

// GetDeletedAt returns the deletion timestamp
func (b *BaseCollection) GetDeletedAt() *time.Time {
	return b.DeletedAt
}
