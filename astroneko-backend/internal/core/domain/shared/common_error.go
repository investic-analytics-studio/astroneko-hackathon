package shared

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"
)

var (
	// Firebase & Auth related errors
	ErrFirebaseNotInitialized     = errors.New("firebase client not initialized")
	ErrFirebaseWebAPIKeyMissing   = errors.New("firebase web API key not configured")
	ErrFirebaseCredentialInvalid  = errors.New("firebase credential invalid")
	ErrFirebaseTokenInvalid       = errors.New("invalid firebase token")
	ErrFirebaseTokenExpired       = errors.New("firebase token expired")
	ErrFirebaseTokenDecodeFailed  = errors.New("failed to decode firebase token")
	ErrFirebaseUserNotFound       = errors.New("firebase user not found")
	ErrFirebaseUserCreationFailed = errors.New("failed to create firebase user")
	ErrFirebaseUserDeletionFailed = errors.New("failed to delete firebase user")
	ErrFirebaseUIDMismatch        = errors.New("firebase UID mismatch")
	ErrFirebaseAuthFailed         = errors.New("firebase authentication failed")
	ErrFirebaseTokenRefreshFailed = errors.New("failed to refresh firebase token")
	ErrFirebaseTokenVerifyFailed  = errors.New("failed to verify firebase token")

	// Auth related errors
	ErrInvalidRefreshToken          = errors.New("invalid refresh token")
	ErrInvalidAccessToken           = errors.New("invalid access token")
	ErrInvalidTokenType             = errors.New("invalid token type")
	ErrUnexpectedSigningMethod      = errors.New("unexpected signing method")
	ErrPasswordMismatch             = errors.New("passwords do not match")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrPasswordExpired              = errors.New("password expired")
	ErrLatestPasswordPreviouslyUsed = errors.New("password previously used")
	ErrEmailNotFoundInToken         = errors.New("email not found in token")
	ErrGoogleTokenInvalid           = errors.New("invalid google token")
	ErrGoogleTokenEmailMissing      = errors.New("email not found in google token")

	// User related errors
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserAlreadyExists             = errors.New("user already exists")
	ErrUserCreationFailed            = errors.New("failed to create user")
	ErrUserUpdateFailed              = errors.New("failed to update user")
	ErrUserDeletionFailed            = errors.New("failed to delete user")
	ErrUserEmailAlreadyExists        = errors.New("user with this email already exists")
	ErrInvalidReferralCode           = errors.New("invalid referral code")
	ErrReferralAlreadyActivated      = errors.New("referral already activated")
	ErrWaitingListUserAlreadyExists  = errors.New("user already exists in waiting list")
	ErrWaitingListUserCreationFailed = errors.New("failed to add user to waiting list")

	// General errors
	ErrInvalidRequest           = errors.New("invalid request")
	ErrInternalServerError      = errors.New("internal server error")
	ErrDatabaseConnectionFailed = errors.New("database connection failed")
)

type CustomError struct {
	HTTPStatus int
	Code       string
	Module     string
	Message    string
	Details    string
}

var errorRegistry = map[string]CustomError{
	// Firebase & Auth errors (1xxx)
	"ERR_1001": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1001",
		Module:     "firebase",
		Message:    "Firebase client not initialized",
		Details:    "Firebase client is not properly configured"},
	"ERR_1002": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1002",
		Module:     "firebase",
		Message:    "Firebase web API key missing",
		Details:    "Firebase web API key is not configured"},
	"ERR_1003": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1003",
		Module:     "firebase",
		Message:    "Invalid Firebase token",
		Details:    "Firebase token is invalid or malformed"},
	"ERR_1004": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1004",
		Module:     "firebase",
		Message:    "Firebase token expired",
		Details:    "Firebase token has expired"},
	"ERR_1005": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1005",
		Module:     "firebase",
		Message:    "Failed to decode Firebase token",
		Details:    "Error decoding Firebase token"},
	"ERR_1006": {
		HTTPStatus: http.StatusNotFound,
		Code:       "ERR_1006",
		Module:     "firebase",
		Message:    "Firebase user not found",
		Details:    "User not found in Firebase"},
	"ERR_1007": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1007",
		Module:     "firebase",
		Message:    "Failed to create Firebase user",
		Details:    "Error creating user in Firebase"},
	"ERR_1008": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1008",
		Module:     "firebase",
		Message:    "Failed to delete Firebase user",
		Details:    "Error deleting user from Firebase"},
	"ERR_1009": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1009",
		Module:     "firebase",
		Message:    "Firebase UID mismatch",
		Details:    "Firebase UID does not match database record"},
	"ERR_1010": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1010",
		Module:     "firebase",
		Message:    "Firebase authentication failed",
		Details:    "Firebase authentication request failed"},
	"ERR_1011": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1011",
		Module:     "firebase",
		Message:    "Failed to refresh Firebase token",
		Details:    "Error refreshing Firebase token"},
	"ERR_1012": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1012",
		Module:     "firebase",
		Message:    "Failed to verify Firebase token",
		Details:    "Error verifying Firebase token"},
	"ERR_1013": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1013",
		Module:     "auth",
		Message:    "Invalid refresh token",
		Details:    "Refresh token is invalid or expired"},
	"ERR_1014": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1014",
		Module:     "auth",
		Message:    "Invalid access token",
		Details:    "Access token is invalid or expired"},
	"ERR_1015": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1015",
		Module:     "auth",
		Message:    "Invalid token type",
		Details:    "Token type is not supported"},
	"ERR_1016": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1016",
		Module:     "auth",
		Message:    "Password mismatch",
		Details:    "Provided password does not match"},
	"ERR_1017": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1017",
		Module:     "auth",
		Message:    "Invalid password",
		Details:    "Password does not meet requirements"},
	"ERR_1018": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1018",
		Module:     "auth",
		Message:    "Password expired",
		Details:    "Password has expired and needs to be changed"},
	"ERR_1019": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1019",
		Module:     "auth",
		Message:    "Password previously used",
		Details:    "New password must be different from previous passwords"},
	"ERR_1020": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1020",
		Module:     "auth",
		Message:    "Email not found in token",
		Details:    "Email claim not found in authentication token"},
	"ERR_1021": {
		HTTPStatus: http.StatusUnauthorized,
		Code:       "ERR_1021",
		Module:     "auth",
		Message:    "Invalid Google token",
		Details:    "Google authentication token is invalid"},
	"ERR_1022": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1022",
		Module:     "auth",
		Message:    "Email not found in Google token",
		Details:    "Email claim not found in Google token"},
	"ERR_1023": {
		HTTPStatus: http.StatusNotFound,
		Code:       "ERR_1023",
		Module:     "user",
		Message:    "User not found",
		Details:    "User not found in database"},
	"ERR_1024": {
		HTTPStatus: http.StatusConflict,
		Code:       "ERR_1024",
		Module:     "user",
		Message:    "User already exists",
		Details:    "User with this email already exists"},
	"ERR_1025": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1025",
		Module:     "user",
		Message:    "Failed to create user",
		Details:    "Error creating user in database"},
	"ERR_1026": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1026",
		Module:     "user",
		Message:    "Failed to update user",
		Details:    "Error updating user in database"},
	"ERR_1027": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1027",
		Module:     "user",
		Message:    "Failed to delete user",
		Details:    "Error deleting user from database"},
	"ERR_1028": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1028",
		Module:     "user",
		Message:    "User email already exists",
		Details:    "User with this email address already exists"},
	"ERR_1029": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1029",
		Module:     "general",
		Message:    "Invalid request",
		Details:    "Request validation failed"},
	"ERR_1030": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1030",
		Module:     "general",
		Message:    "Internal server error",
		Details:    "An unexpected error occurred"},
	"ERR_1031": {
		HTTPStatus: http.StatusServiceUnavailable,
		Code:       "ERR_1031",
		Module:     "general",
		Message:    "Database connection failed",
		Details:    "Unable to connect to database"},
	"ERR_1032": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1032",
		Module:     "user",
		Message:    "Invalid referral code",
		Details:    "The provided referral code is invalid"},
	"ERR_1033": {
		HTTPStatus: http.StatusBadRequest,
		Code:       "ERR_1033",
		Module:     "user",
		Message:    "Referral already activated",
		Details:    "Referral has already been activated for this user"},
	"ERR_1034": {
		HTTPStatus: http.StatusConflict,
		Code:       "ERR_1034",
		Module:     "waiting_list",
		Message:    "User already exists in waiting list",
		Details:    "User with this email already exists in the waiting list"},
	"ERR_1035": {
		HTTPStatus: http.StatusInternalServerError,
		Code:       "ERR_1035",
		Module:     "waiting_list",
		Message:    "Failed to add user to waiting list",
		Details:    "Error adding user to waiting list in database"},
}

func NewErrorResponse(code string, detailOverride ...string) (int, ResponseBody) {
	err := errorRegistry[code]

	if len(detailOverride) > 0 {
		err.Details = detailOverride[0]
	}

	return err.HTTPStatus, ResponseBody{
		Status: Status{
			Code:    err.Code,
			Message: []string{err.Message, err.Details},
		},
	}
}

func WrapDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			if strings.Contains(pgErr.ConstraintName, "idx_auth_users_email") {
				return fmt.Errorf("email already exists")
			}
			if strings.Contains(pgErr.ConstraintName, "idx_auth_roles_name") {
				return fmt.Errorf("role name already exists")
			}
			return fmt.Errorf("duplicate value violates unique constraint: %s", pgErr.ConstraintName)
		default:
			return fmt.Errorf("database error: %s", pgErr.Message)
		}
	}

	return err
}
