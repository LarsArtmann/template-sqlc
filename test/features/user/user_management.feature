Feature: User Management
  As a system administrator
  I want to manage user accounts
  So that users can access the system with proper authorization

  Background:
    Given a clean user system
    And the event publisher is cleared

  Scenario: Create a user with valid data
    When I create a user with valid data
    Then the user should be created successfully
    And a user created event should be published

  Scenario: Create user with existing email
    Given a user with email "test@example.com" and username "testuser1"
    When I create a user with email "test@example.com"
    Then I should receive a "user already exists" error

  Scenario: Create user with existing username
    Given a user with email "test1@example.com" and username "testuser"
    When I create a user with username "testuser"
    Then I should receive a "user already exists" error

  Scenario: Create user with invalid email
    When I create a user with email "invalid-email"
    Then I should receive a validation error

  Scenario: Create user with invalid username
    When I create a user with username "ab"
    Then I should receive a validation error

  Scenario: Authenticate user with valid credentials
    Given I have valid user credentials
    When I attempt to authenticate with these credentials
    Then the authentication should succeed
    And the session should be created
    And a user login event should be published

  Scenario: Authenticate user with invalid credentials
    Given I have invalid user credentials
    When I attempt to authenticate with these credentials
    Then the authentication should fail
    And a user login failed event should be published

  Scenario: Authenticate with inactive user
    Given an inactive user account
    And I have valid user credentials for this account
    When I attempt to authenticate with these credentials
    Then the authentication should fail

  Scenario: Authenticate with suspended user
    Given a suspended user account
    And I have valid user credentials for this account
    When I attempt to authenticate with these credentials
    Then the authentication should fail

  Scenario: Update user profile
    Given an active user account
    When I update the user profile
    Then the user profile should be updated
    And a user updated event should be published

  Scenario: Change user role
    Given an active user account
    When I change the user role to "admin"
    Then the user role should be changed to "admin"
    And a role changed event should be published

  Scenario: Verify user account
    Given an active user account
    When I verify the user account
    Then the user account should be verified
    And a user verified event should be published

  Scenario: Deactivate user account
    Given an active user account
    When I deactivate the user account
    Then the user account should be deactivated
    And a user updated event should be published

  Scenario: Get user statistics
    Given multiple user accounts with different statuses
    And a user account with status "active"
    And a user account with status "inactive"
    And a user account with status "suspended"
    When I get user statistics
    Then the statistics should include counts for each status

  Scenario: Create user with metadata and tags
    When I create a user with valid data
    Then the user should have the specified metadata
    And the user should have the specified tags

  Scenario: Create admin user
    When I create a user with valid data
    And the user role is "admin"
    Then the user should have admin privileges
    And a user created event should be published

  Scenario: Create moderator user
    When I create a user with valid data
    And the user role is "moderator"
    Then the user should have moderator privileges
    And a user created event should be published

  Scenario: Create user with pending status
    When I create a user with valid data
    And the user status is "pending"
    Then the user account should be in pending state
    And the user should not be able to authenticate

  Scenario: Create user with suspended status
    When I create a user with valid data
    And the user status is "suspended"
    Then the user account should be suspended
    And the user should not be able to authenticate

  Scenario: User session expiration
    Given I have valid user credentials
    When I attempt to authenticate with these credentials
    Then the authentication should succeed
    And the session should be created
    When the session expires
    Then the session should no longer be valid

  Scenario: Multiple active sessions
    Given I have valid user credentials
    When I attempt to authenticate from multiple devices
    Then multiple sessions should be created
    And all sessions should be active