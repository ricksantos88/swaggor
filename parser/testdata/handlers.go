package testdata

import "net/http"

// ListUsers returns paginated users.
//
// @Route    GET /api/v1/users
// @Summary  List Users
// @Desc     Returns all users with pagination
// @Tags     users,admin
// @Query    page  "Page number"  optional
// @Query    limit "Results per page" required
// @Response 200  UserResponse  "Successful"
// @Response 401  ErrorResponse "Unauthorized"
// @Auth     bearer
// @Cache    60s
// @For      nethttp
func ListUsers(w http.ResponseWriter, r *http.Request) {}

// GetUser returns a single user by ID.
//
// @Route   GET /api/v1/users/{id}
// @Summary Get User
// @Path    id "User ID"
// @Response 200 UserResponse "OK"
// @Response 404 ErrorResponse "Not Found"
// @Auth    bearer
func GetUser(w http.ResponseWriter, r *http.Request) {}

// CreateUser creates a new user.
//
// @Route   POST /api/v1/users
// @Summary Create User
// @Body    "User payload" required
// @Response 201 UserResponse "Created"
// @Response 400 ErrorResponse "Bad Request"
// @Auth    bearer
// @For     fiber
func CreateUser(w http.ResponseWriter, r *http.Request) {}

// NotAnnotated is a handler without annotations — should be ignored.
func NotAnnotated(w http.ResponseWriter, r *http.Request) {}

// OnlyComment has a comment but no @Route tag — should be ignored.
//
// This is just a regular comment.
func OnlyComment(w http.ResponseWriter, r *http.Request) {}
