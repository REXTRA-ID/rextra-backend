# Rextra Authentication System Documentation

This document provides a comprehensive overview of the authentication and authorization system for the Rextra backend.

## 1. Overview

The authentication system is built around JSON Web Tokens (JWT) with a refresh token mechanism. It supports standard email/password registration, social login via Google, and role-based access control (RBAC). Key features include email verification for new accounts and a secure password reset flow.

- **Access Tokens**: Short-lived JWTs (24 Hours) used to authenticate API requests.
- **Refresh Tokens**: Long-lived tokens (30 days) stored in the database, used to obtain new access tokens without requiring users to log in again.
- **Providers**: Supports local (email/password) and Google authentication.
- **RBAC**: Implemented via middleware to restrict access to certain endpoints based on user roles (`USER`, `ADMIN`, `EXPERT`).

---

## 2. Authentication Middleware

Authentication is enforced and managed by middleware functions.

### `Authenticate()`

This middleware protects endpoints that require a logged-in user.

- **Action**: It inspects the `Authorization` header for a `Bearer <token>`.
- **Validation**: It verifies the JWT's signature and expiration.
- **Context**: If the token is valid, it extracts the user's `user_id`, `email`, and `role` from the token payload and adds them to the Gin context for use in subsequent handlers.
- **Rejection**: If the token is missing, malformed, or invalid, it aborts the request with an authorization error.

### `OnlyAllow(roles ...string)`

This middleware provides role-based access control and must be used *after* the `Authenticate()` middleware.

- **Action**: It checks the `role` value set in the Gin context by `Authenticate()`.
- **Validation**: It ensures the user's role is present in the list of allowed roles passed to the function.
- **Rejection**: If the user's role is not permitted, it aborts the request with a forbidden error.

---

## 3. Core Entities

### User (`users` table)

Represents a user in the system.

| Field | Type | Description |
| :--- | :--- | :--- |
| `ID` | `uuid` | Primary Key. |
| `ProfileImageUrl` | `string` | URL for the user's profile picture. |
| `Fullname` | `string` | The user's full name. |
| `Email` | `string` | The user's unique email address. |
| `Password` | `string` | The hashed user password. |
| `IsVerified` | `bool` | Flag indicating if the user has verified their email. Defaults to `false`. |
| `PhoneNumber` | `string` | The user's phone number. |
| `Role` | `Role` | User role (`USER`, `ADMIN`, `EXPERT`). Defaults to `USER`. |

### Session (`sessions` table)

Stores refresh tokens for active user sessions.

| Field | Type | Description |
| :--- | :--- | :--- |
| `ID` | `uuid` | Primary Key. |
| `UserID` | `string` | Foreign key to the `users` table. |
| `Token` | `string` | The unique refresh token string. |
| `ExpiresAt` | `time.Time` | The expiration date of the refresh token. |
| `IsActive` | `bool` | Flag indicating if the session is active. Defaults to `true`. |
| `AuthProvider` | `string` | The provider used for login (e.g., "local", "google.com"). |
| `DeviceInfo` | `jsonb` | Optional information about the user's device. |

---

## 4. API Endpoints & Flows

All authentication routes are prefixed with `/api/v1/auth`.

### Register User

- **Endpoint**: `POST /api/v1/auth/register`
- **Authentication**: None
- **Description**: Creates a new standard user account.
- **Flow**:
  1. Checks if a user with the provided email already exists.
  2. Hashes the user's password.
  3. Creates a new user record with `IsVerified: false` and `Role: USER`.
  4. Generates a short-lived verification token and sends it in an email to the user.
- **Request Body**:
  ```json
  {
    "fullname": "John Doe",
    "email": "john.doe@example.com",
    "password": "a-strong-password",
    "phone_number": "1234567890"
  }
  ```
- **Success Response (200)**:
  ```json
  {
    "id": "user-uuid",
    "fullname": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "1234567890",
    "role": "USER"
  }
  ```

### Register Admin

- **Endpoint**: `POST /api/v1/auth/admin`
- **Authentication**: None
- **Description**: Creates a new admin user. Requires a secret token for authorization.
- **Request Body**:
  ```json
  {
    "fullname": "Admin User",
    "email": "admin@example.com",
    "password": "a-strong-password",
    "phone_number": "0987654321",
    "token": "your-secret-admin-token"
  }
  ```
- **Success Response (200)**: Similar to user registration, but with `"role": "ADMIN"`.

### Login

- **Endpoint**: `POST /api/v1/auth/login`
- **Authentication**: None
- **Description**: Authenticates a user with email and password.
- **Flow**:
  1. Validates user credentials.
  2. Checks if the user's account is verified (`IsVerified: true`).
  3. Generates and returns a short-lived JWT `access_token` and a long-lived `refresh_token`.
  4. Creates a new record in the `sessions` table to store the refresh token.
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "a-strong-password"
  }
  ```
- **Success Response (200)**:
  ```json
  {
    "access_token": "jwt-access-token",
    "refresh_token": "database-refresh-token",
    "role": "USER"
  }
  ```

### Login with Google

- **Endpoint**: `POST /api/v1/auth/google`
- **Authentication**: None
- **Description**: Authenticates a user or creates a new account using a Google ID token.
- **Flow**:
  1. Verifies the `id_token` with Firebase Auth.
  2. If no user exists with the email from the token, a new user is created with `IsVerified: true`.
  3. Generates and returns `access_token` and `refresh_token` as in the standard login flow.
- **Request Body**:
  ```json
  {
    "id_token": "google-id-token-from-frontend"
  }
  ```
- **Success Response (200)**: Same as standard login.

### Verify Account

- **Endpoint**: `GET /api/v1/auth/verify`
- **Authentication**: None
- **Description**: Verifies a user's email address using a token sent during registration.
- **Query Parameter**: `token` (The verification token from the email link).
- **Flow**:
  1. Validates the token.
  2. Updates the corresponding user's `IsVerified` status to `true`.
- **Success Response (200)**: A success message.

### Resend Verification Email

- **Endpoint**: `POST /api/v1/auth/send-email`
- **Authentication**: None
- **Description**: Resends the verification email to a user.
- **Request Body**:
  ```json
  {
    "email": "unverified.user@example.com"
  }
  ```
- **Success Response (200)**: A success message.

### Forget Password

- **Endpoint**: `POST /api/v1/auth/forget`
- **Authentication**: None
- **Description**: Initiates the password reset process for a verified user.
- **Flow**:
  1. Finds the user by email.
  2. Generates a short-lived token and sends a password reset link to the user's email.
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com"
  }
  ```
- **Success Response (200)**: A success message.

### Change Password

- **Endpoint**: `POST /api/v1/auth/change`
- **Authentication**: None
- **Description**: Sets a new password for a user using a token from the "Forget Password" flow.
- **Query Parameter**: `token` (The password reset token from the email link).
- **Request Body**:
  ```json
  {
    "new_password": "my-new-strong-password"
  }
  ```
- **Success Response (200)**: A success message.

### Get User Profile

- **Endpoint**: `GET /api/v1/auth/me`
- **Authentication**: Required (`Authenticate()` middleware)
- **Description**: Retrieves the profile information of the currently authenticated user.
- **Success Response (200)**:
  ```json
  {
    "personal_info": {
      "id": "user-uuid",
      "fullname": "John Doe",
      "email": "john.doe@example.com",
      "phone_number": "1234567890",
      "role": "USER"
    }
  }
  ```

### Logout

- **Endpoint**: `DELETE /api/v1/auth/logout`
- **Authentication**: Required (`Authenticate()` middleware)
- **Description**: Logs out the user by invalidating their refresh token.
- **Flow**:
  1. Finds the session in the database using the `refresh_token`.
  2. Marks the session as inactive and then deletes it.
- **Request Body**:
  ```json
  {
    "refresh_token": "database-refresh-token"
  }
  ```
- **Success Response (200)**: A success message.
