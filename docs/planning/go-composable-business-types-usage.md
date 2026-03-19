# go-composable-business-types Integration Plan

## Executive Summary

This document outlines how to integrate the `go-composable-business-types` library into the template-sqlc project. The library provides strongly-typed, composable base values that improve type safety, prevent bugs, and enhance code clarity.

## Current State Analysis

### Existing Entity Types

The project currently defines custom type aliases for entity IDs:

```go
// Current implementation in internal/domain/entities/user.go
type UserID int64

func (id UserID) Int64() int64   { return int64(id) }
func (id UserID) String() string { return fmt.Sprintf("user:%d", id) }
```

```go
// Current implementation in internal/domain/entities/session.go
type SessionID int64

func (id SessionID) Int64() int64   { return int64(id) }
func (id SessionID) String() string { return fmt.Sprintf("session:%d", id) }
```

### Current Limitations

1. **No compile-time type safety between IDs**: `UserID` and `SessionID` are both `int64`, so they can be accidentally mixed
2. **Manual implementation**: Each ID type requires custom methods (`Int64()`, `String()`)
3. **No built-in serialization**: JSON, SQL, and binary marshaling must be implemented manually
4. **No zero-value handling**: No standard way to check if an ID is unset
5. **No comparison utilities**: Sorting and equality checks require manual implementation

## Proposed Integration Strategy

### Phase 1: ID Package Integration (High Priority)

The `id` package provides branded, type-safe identifiers that prevent mixing different entity IDs at compile time.

#### Benefits

- **Compile-time safety**: Cannot accidentally pass a UserID where a SessionID is expected
- **Zero boilerplate**: No need to implement `String()`, comparison, or serialization methods
- **Full serialization support**: JSON, SQL, binary, and text marshaling out of the box
- **Zero value utilities**: `IsZero()`, `Or()` methods for safe handling
- **Generic design**: Works with any underlying type (string, int64, UUID, etc.)

#### Implementation

```go
package entities

import (
    "github.com/larsartmann/go-composable-business-types/id"
)

// Brand types (unexported structs provide compile-time isolation)
type userBrand struct{}
type sessionBrand struct{}

// ID type aliases using branded IDs
type UserID = id.ID[userBrand, int64]
type SessionID = id.ID[sessionBrand, int64]

// Constructors
func NewUserID(value int64) UserID {
    return id.NewID[userBrand, int64](value)
}

func NewSessionID(value int64) SessionID {
    return id.NewID[sessionBrand, int64](value)
}
```

#### Usage Examples

```go
// Type-safe creation
userID := NewUserID(123)
sessionID := NewSessionID(456)

// These would be caught at compile time:
// GetUser(sessionID)  // ERROR: cannot use SessionID as UserID

// Zero value handling
if userID.IsZero() {
    // Handle missing ID
}

// Default value with Or
defaultID := NewUserID(0)
actualID := userID.Or(defaultID)

// Serialization (automatic)
data, _ := json.Marshal(userID)        // "123"
var restored UserID
json.Unmarshal(data, &restored)        // Works seamlessly

// Database operations (automatic)
row.Scan(&userID)                      // SQL scan
_, err := db.Exec("INSERT ...", userID) // SQL value
```

### Phase 2: NanoId Integration (Recommended for New IDs)

The `nanoid` package provides URL-safe, cryptographically secure random identifiers.

#### Benefits

- **URL-safe**: Uses `A-Za-z0-9_-` alphabet
- **Unique**: Cryptographically secure random generation
- **FIPS-140 compatible**: High-performance implementation
- **Configurable length**: Default 21 characters (126 bits entropy)
- **Better than UUIDs**: Shorter, URL-safe, no special characters

#### Implementation for String-Based IDs

```go
package entities

import (
    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

// For entities that use string IDs (e.g., external references)
type externalUserBrand struct{}
type ExternalUserID = id.ID[externalUserBrand, nanoid.NanoId]

func NewExternalUserID() ExternalUserID {
    return id.NewID[externalUserBrand](nanoid.NewNanoId())
}

func ParseExternalUserID(s string) (ExternalUserID, error) {
    n, err := nanoid.ParseNanoId(s)
    if err != nil {
        return ExternalUserID{}, err
    }
    return id.NewID[externalUserBrand](n), nil
}
```

### Phase 3: BoundedString Integration (Medium Priority)

The `bounded` package provides strings with validated length constraints.

#### Current Implementation

```go
// Current manual validation
type Username string

func NewUsername(username string) (Username, error) {
    if len(username) < 3 || len(username) > 50 {
        return "", ErrInvalidUsername
    }
    return Username(strings.TrimSpace(username)), nil
}
```

#### Proposed Implementation

```go
package entities

import (
    "github.com/larsartmann/go-composable-business-types/bounded"
)

// Factory functions for domain-specific bounded strings
var NewUsername = bounded.TrimmedBoundedStringOf(3, 50)
var NewFirstName = bounded.TrimmedBoundedStringOf(1, 100)
var NewLastName = bounded.TrimmedBoundedStringOf(1, 100)

// Usage
username, err := NewUsername("john_doe")
if err != nil {
    return err
}
// username is now a bounded.BoundedString with guaranteed constraints
```

#### Benefits

- **Centralized validation**: Constraints defined once, enforced everywhere
- **Trimmed input**: Automatic whitespace handling
- **Clear error messages**: Validation errors are descriptive
- **Composability**: Can be wrapped in higher-level domain types

### Phase 4: Email Type Integration (Medium Priority)

Replace custom email validation with the library's `Email` type.

#### Current Implementation

```go
// Current: undefined function isValidEmail (bug!)
type Email string

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return "", ErrInvalidEmail
    }
    return Email(strings.ToLower(strings.TrimSpace(email))), nil
}
```

#### Proposed Implementation

```go
package entities

import (
    "github.com/larsartmann/go-composable-business-types/bounded"
)

// Use the library's Email type (or bounded string with validation)
type Email = bounded.BoundedString

func NewEmail(email string) (Email, error) {
    // Validate format first
    if !isValidEmailFormat(email) {
        return Email{}, ErrInvalidEmail
    }
    // Use bounded string for length validation
    return bounded.TrimmedBoundedString(1, 254, email)
}
```

### Phase 5: DataPoint Integration (Future Consideration)

The `datapoint` package provides self-contained data units with complete audit trails.

#### Use Case

For entities that require full audit history (who, when, why):

```go
package entities

import (
    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/actor"
)

// DataPoint wraps any entity with full audit trail
type AuditedUser struct {
    datapoint.DataPoint[UserState]
}

type UserState struct {
    Email    string
    Username string
    Status   string
}

// Creation with audit trail
user := datapoint.NewDataPoint(
    UserState{Email: "john@example.com", Username: "john"},
    actor.UserActor(userID, "john@example.com"),
).WithReason("user registration").WithTrigger(enums.TriggerManual)
```

#### When to Use

- **Event sourcing**: Building audit/lineage graphs
- **Compliance requirements**: Full traceability without external systems
- **Temporal queries**: "What was the state at time X?"

## Implementation Plan

### Step 1: Add Dependency

```bash
# In the appropriate module directory
go get github.com/larsartmann/go-composable-business-types/id
go get github.com/larsartmann/go-composable-business-types/nanoid
go get github.com/larsartmann/go-composable-business-types/bounded
```

### Step 2: Refactor UserID

1. Replace `type UserID int64` with branded ID
2. Update all references to use `NewUserID()` constructor
3. Verify JSON serialization works (SQLC, API responses)
4. Run tests to ensure no regressions

### Step 3: Refactor SessionID

1. Replace `type SessionID int64` with branded ID
2. Update `UserSession.userID` field type
3. Verify foreign key relationships work
4. Run integration tests

### Step 4: Add BoundedString for Username/Name Fields

1. Replace manual validation with `bounded` package
2. Update constructors to use factory functions
3. Verify error handling still works

### Step 5: Fix Email Validation

1. Implement proper email validation or use library type
2. Fix the `isValidEmail` undefined error

## Compatibility Considerations

### SQLC Integration

SQLC generates code from SQL queries. The branded ID types must work with:

1. **Scanning from database**: `ID` implements `sql.Scanner`
2. **Value for database**: `ID` implements `driver.Valuer`
3. **JSON marshaling**: For API responses

```sql
-- SQLC will generate code that scans int64 into UserID
SELECT id, email, username FROM users WHERE id = ?;
```

Generated code will work if `UserID` is underlying `int64`.

### Database Schema

No database schema changes required. The `id` package stores the underlying value (`int64`, `string`, etc.) in the database.

### API Compatibility

JSON serialization produces the same output:

```json
// Before and After: Same output
{
  "id": 123,
  "email": "user@example.com"
}
```

## Testing Strategy

1. **Unit tests**: Verify each ID type can be created, compared, and serialized
2. **Integration tests**: Verify SQLC-generated code works with new types
3. **API tests**: Verify JSON serialization produces expected output
4. **Compile-time tests**: Attempt to mix IDs (should fail to compile)

## Migration Checklist

- [ ] Add `go-composable-business-types` dependency
- [ ] Create brand types for each entity
- [ ] Replace `UserID` with branded ID
- [ ] Replace `SessionID` with branded ID
- [ ] Update repository interfaces
- [ ] Update service layer
- [ ] Update mappers
- [ ] Fix `isValidEmail` function
- [ ] Add bounded string for username/name fields
- [ ] Run full test suite
- [ ] Verify SQLC generation still works
- [ ] Update documentation

## Risk Assessment

| Risk                        | Likelihood | Impact | Mitigation                                 |
| --------------------------- | ---------- | ------ | ------------------------------------------ |
| SQLC compatibility issues   | Low        | High   | Test with generated code before merge      |
| JSON serialization changes  | Low        | Medium | Verify API responses match expected format |
| Breaking changes in library | Low        | Medium | Pin to specific version                    |
| Team learning curve         | Medium     | Low    | Document patterns and provide examples     |

## Recommendation

**Proceed with Phase 1 (ID package) and Phase 2 (NanoId)**. These provide immediate value with minimal risk:

1. Compile-time type safety prevents bugs
2. Zero boilerplate reduces code maintenance
3. Full serialization support eliminates manual implementation
4. Backward compatible with existing database schema

**Defer Phase 5 (DataPoint)** until there's a clear requirement for audit trails and temporal queries.

## Example: Complete Refactored entities/user.go

```go
package entities

import (
    "time"

    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/larsartmann/go-composable-business-types/bounded"
)

// Brand types for compile-time isolation
type userBrand struct{}

// UserID is a branded, type-safe identifier
type UserID = id.ID[userBrand, int64]

// NewUserID creates a new UserID from an int64
func NewUserID(value int64) UserID {
    return id.NewID[userBrand, int64](value)
}

// Email is a validated email address
type Email = bounded.BoundedString

// NewEmail creates a validated email
func NewEmail(email string) (Email, error) {
    if !isValidEmailFormat(email) {
        return Email{}, ErrInvalidEmail
    }
    return bounded.TrimmedBoundedString(1, 254, email)
}

// Username is a validated username
type Username = bounded.BoundedString

// NewUsername creates a validated username
var NewUsername = bounded.TrimmedBoundedStringOf(3, 50)

// FirstName is a validated first name
type FirstName = bounded.BoundedString

// NewFirstName creates a validated first name
var NewFirstName = bounded.TrimmedBoundedStringOf(1, 100)

// LastName is a validated last name
type LastName = bounded.BoundedString

// NewLastName creates a validated last name
var NewLastName = bounded.TrimmedBoundedStringOf(1, 100)

// User entity (unchanged structure, updated types)
type User struct {
    id          UserID
    email       Email
    username    Username
    firstName   FirstName
    lastName    LastName
    status      UserStatus
    role        UserRole
    isVerified  bool
    createdAt   time.Time
    updatedAt   time.Time
}

// Methods (simplified due to library support)
func (u *User) ID() UserID       { return u.id }
func (u *User) Email() Email      { return u.email }
func (u *User) Username() Username { return u.username }

// ID methods now provided by library:
// - userID.String()     // Returns "123"
// - userID.IsZero()     // Returns false
// - userID.Equal(other) // Type-safe comparison
// - json.Marshal(userID) // Returns "123"
```

## Conclusion

The `go-composable-business-types` library aligns perfectly with the project's goal of type-safe, maintainable code. The integration provides immediate value through compile-time safety and reduced boilerplate, while remaining compatible with the existing SQLC-generated code and database schema.
