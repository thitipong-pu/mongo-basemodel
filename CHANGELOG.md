# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-05-30

### Added
- Initial release of MongoDB BaseModel library
- `BaseCollection` struct with common MongoDB fields:
  - `_id` (ObjectID)
  - `created_at` (timestamp)
  - `updated_at` (timestamp, nullable)
  - `deleted_at` (timestamp, nullable for soft delete)
- Core methods:
  - `SetInsertMeta()` - Sets ObjectID and creation timestamp
  - `SetUpdateMeta()` - Sets update timestamp
  - `SetDeleteMeta()` - Sets soft delete timestamp
  - `IsDeleted()` - Checks if record is soft deleted
- Getter methods:
  - `GetID()` - Returns ObjectID as string
  - `GetCreatedAt()` - Returns creation timestamp
  - `GetUpdatedAt()` - Returns update timestamp
  - `GetDeletedAt()` - Returns deletion timestamp
- Comprehensive test coverage with unit tests and benchmarks
- Example implementations:
  - Basic usage example
  - MongoDB integration example with full CRUD operations
- Repository pattern example with User model
- Complete documentation with usage examples
- MIT License

### Features
- Easy integration with existing Go MongoDB projects
- Automatic metadata management for insert, update, and delete operations
- Support for soft delete functionality
- Type-safe ObjectID handling
- Clean and simple API design
- Struct embedding support for easy inheritance
- JSON and BSON tag support for serialization

### Documentation
- Comprehensive README with usage examples
- Multiple example projects demonstrating different use cases
- Best practices guide for MongoDB operations
- Repository pattern implementation example 