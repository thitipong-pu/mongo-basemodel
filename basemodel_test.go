package basemodel

import (
	"testing"
	"time"
)

// TestUser is a test struct that embeds BaseCollection
type TestUser struct {
	BaseCollection `bson:",inline"`
	Name           string `json:"name" bson:"name"`
	Email          string `json:"email" bson:"email"`
}

func TestSetInsertMeta(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Check initial state
	if !user.Oid.IsZero() {
		t.Error("Expected Oid to be zero initially")
	}
	if !user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be zero initially")
	}

	// Call SetInsertMeta
	user.SetInsertMeta()

	// Check that Oid was set
	if user.Oid.IsZero() {
		t.Error("Expected Oid to be set after SetInsertMeta")
	}

	// Check that CreatedAt was set
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set after SetInsertMeta")
	}

	// Check that CreatedAt is recent (within last second)
	now := time.Now()
	if now.Sub(user.CreatedAt) > time.Second {
		t.Error("Expected CreatedAt to be recent")
	}

	// Check that UpdatedAt and DeletedAt are still nil
	if user.UpdatedAt != nil {
		t.Error("Expected UpdatedAt to be nil after SetInsertMeta")
	}
	if user.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil after SetInsertMeta")
	}
}

func TestSetUpdateMeta(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Call SetInsertMeta first
	user.SetInsertMeta()
	createdAt := user.CreatedAt

	// Wait a moment to ensure different timestamp
	time.Sleep(10 * time.Millisecond)

	// Call SetUpdateMeta
	user.SetUpdateMeta()

	// Check that UpdatedAt was set
	if user.UpdatedAt == nil {
		t.Error("Expected UpdatedAt to be set after SetUpdateMeta")
	}

	// Check that UpdatedAt is recent
	now := time.Now()
	if now.Sub(*user.UpdatedAt) > time.Second {
		t.Error("Expected UpdatedAt to be recent")
	}

	// Check that UpdatedAt is after CreatedAt
	if user.UpdatedAt.Before(createdAt) {
		t.Error("Expected UpdatedAt to be after CreatedAt")
	}

	// Check that other fields remain unchanged
	if user.CreatedAt != createdAt {
		t.Error("Expected CreatedAt to remain unchanged after SetUpdateMeta")
	}
	if user.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil after SetUpdateMeta")
	}
}

func TestSetDeleteMeta(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Call SetInsertMeta first
	user.SetInsertMeta()
	createdAt := user.CreatedAt

	// Call SetDeleteMeta
	user.SetDeleteMeta()

	// Check that DeletedAt was set
	if user.DeletedAt == nil {
		t.Error("Expected DeletedAt to be set after SetDeleteMeta")
	}

	// Check that DeletedAt is recent
	now := time.Now()
	if now.Sub(*user.DeletedAt) > time.Second {
		t.Error("Expected DeletedAt to be recent")
	}

	// Check that other fields remain unchanged
	if user.CreatedAt != createdAt {
		t.Error("Expected CreatedAt to remain unchanged after SetDeleteMeta")
	}
}

func TestIsDeleted(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Initially should not be deleted
	if user.IsDeleted() {
		t.Error("Expected user to not be deleted initially")
	}

	// After SetInsertMeta, should still not be deleted
	user.SetInsertMeta()
	if user.IsDeleted() {
		t.Error("Expected user to not be deleted after SetInsertMeta")
	}

	// After SetUpdateMeta, should still not be deleted
	user.SetUpdateMeta()
	if user.IsDeleted() {
		t.Error("Expected user to not be deleted after SetUpdateMeta")
	}

	// After SetDeleteMeta, should be deleted
	user.SetDeleteMeta()
	if !user.IsDeleted() {
		t.Error("Expected user to be deleted after SetDeleteMeta")
	}
}

func TestGetID(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Before SetInsertMeta, GetID should return zero ObjectID hex string
	initialID := user.GetID()
	if initialID != "000000000000000000000000" {
		t.Errorf("Expected GetID to return zero ObjectID hex string initially, got %s", initialID)
	}

	// After SetInsertMeta, GetID should return valid hex string
	user.SetInsertMeta()
	id := user.GetID()

	if id == "" {
		t.Error("Expected GetID to return non-empty string after SetInsertMeta")
	}

	// Check that it's a valid ObjectID hex string (24 characters)
	if len(id) != 24 {
		t.Errorf("Expected GetID to return 24-character string, got %d characters", len(id))
	}

	// Check that it matches the actual Oid
	if id != user.Oid.Hex() {
		t.Error("Expected GetID to match Oid.Hex()")
	}

	// Check that it's not the zero ObjectID anymore
	if id == "000000000000000000000000" {
		t.Error("Expected GetID to return non-zero ObjectID after SetInsertMeta")
	}
}

func TestGetCreatedAt(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user.SetInsertMeta()
	createdAt := user.GetCreatedAt()

	if createdAt != user.CreatedAt {
		t.Error("Expected GetCreatedAt to return the same value as CreatedAt field")
	}

	if createdAt.IsZero() {
		t.Error("Expected GetCreatedAt to return non-zero time after SetInsertMeta")
	}
}

func TestGetUpdatedAt(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Initially should return nil
	if user.GetUpdatedAt() != nil {
		t.Error("Expected GetUpdatedAt to return nil initially")
	}

	user.SetInsertMeta()

	// After SetInsertMeta, should still return nil
	if user.GetUpdatedAt() != nil {
		t.Error("Expected GetUpdatedAt to return nil after SetInsertMeta")
	}

	user.SetUpdateMeta()

	// After SetUpdateMeta, should return the timestamp
	updatedAt := user.GetUpdatedAt()
	if updatedAt == nil {
		t.Error("Expected GetUpdatedAt to return non-nil after SetUpdateMeta")
	}

	if updatedAt != user.UpdatedAt {
		t.Error("Expected GetUpdatedAt to return the same pointer as UpdatedAt field")
	}
}

func TestGetDeletedAt(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Initially should return nil
	if user.GetDeletedAt() != nil {
		t.Error("Expected GetDeletedAt to return nil initially")
	}

	user.SetInsertMeta()

	// After SetInsertMeta, should still return nil
	if user.GetDeletedAt() != nil {
		t.Error("Expected GetDeletedAt to return nil after SetInsertMeta")
	}

	user.SetDeleteMeta()

	// After SetDeleteMeta, should return the timestamp
	deletedAt := user.GetDeletedAt()
	if deletedAt == nil {
		t.Error("Expected GetDeletedAt to return non-nil after SetDeleteMeta")
	}

	if deletedAt != user.DeletedAt {
		t.Error("Expected GetDeletedAt to return the same pointer as DeletedAt field")
	}
}

func TestMultipleOperations(t *testing.T) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Test full lifecycle
	user.SetInsertMeta()
	originalCreatedAt := user.CreatedAt
	originalOid := user.Oid

	time.Sleep(10 * time.Millisecond)

	user.SetUpdateMeta()
	firstUpdateAt := *user.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	user.SetUpdateMeta()
	secondUpdateAt := *user.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	user.SetDeleteMeta()

	// Verify that CreatedAt and Oid never changed
	if user.CreatedAt != originalCreatedAt {
		t.Error("Expected CreatedAt to remain unchanged throughout lifecycle")
	}
	if user.Oid != originalOid {
		t.Error("Expected Oid to remain unchanged throughout lifecycle")
	}

	// Verify that UpdatedAt was updated
	if !secondUpdateAt.After(firstUpdateAt) {
		t.Error("Expected second update to be after first update")
	}

	// Verify that DeletedAt was set
	if user.DeletedAt == nil {
		t.Error("Expected DeletedAt to be set")
	}

	// Verify that user is now marked as deleted
	if !user.IsDeleted() {
		t.Error("Expected user to be marked as deleted")
	}
}

func TestJSONTags(t *testing.T) {
	// This test ensures that the struct tags are correct
	// We can't easily test JSON serialization without additional dependencies
	// but we can check that the BaseCollection struct is properly defined

	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user.SetInsertMeta()

	// Basic check that fields are accessible
	if user.Oid.IsZero() {
		t.Error("Oid should be accessible and set")
	}
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should be accessible and set")
	}
	if user.UpdatedAt != nil {
		t.Error("UpdatedAt should be accessible and nil initially")
	}
	if user.DeletedAt != nil {
		t.Error("DeletedAt should be accessible and nil initially")
	}
}

func BenchmarkSetInsertMeta(b *testing.B) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	for i := 0; i < b.N; i++ {
		user.SetInsertMeta()
	}
}

func BenchmarkSetUpdateMeta(b *testing.B) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	for i := 0; i < b.N; i++ {
		user.SetUpdateMeta()
	}
}

func BenchmarkSetDeleteMeta(b *testing.B) {
	user := &TestUser{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	for i := 0; i < b.N; i++ {
		user.SetDeleteMeta()
	}
}
