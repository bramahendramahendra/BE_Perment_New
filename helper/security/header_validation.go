package security

import (
	"fmt"
	"permen_api/helper"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// MaxHeaderSize defines the maximum allowed size for security-sensitive headers
	MaxHeaderSize = 5024 // 1KB limit per header

	// MaxJWTTokenSize defines the maximum allowed size for a JWT bearer token.
	// A well-formed JWT (header.payload.signature) is well under 2 KB.
	// This limit prevents resource exhaustion in cryptographic verification.
	MaxJWTTokenSize = 2048

	// JWT segment limits to cap parser/verification work.
	maxJWTHeaderSegmentSize    = 512
	maxJWTPayloadSegmentSize   = 1024
	maxJWTSignatureSegmentSize = 768

	// Header name constants for consistency and security
	AuthorizationHeader = "Authorization"
	UserqHeader         = "userq"
	HilfmHeader         = "hilfm"
	BranchHeader        = "branch"
	OrgechHeader        = "orgeh"
	StellTXHeader       = "stellTX"
	KostlHeader         = "costCenter"
)

var jwtTokenPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)

// SecureHeaderContext holds validated header information to prevent resource exhaustion
type SecureHeaderContext struct {
	UserQ   string
	Hilfm   string
	Branch  string
	Orgeh   string
	StellTX string
	Kostl   string
}

// ValidateAndGetHeaders safely retrieves and validates headers to prevent resource exhaustion attacks
func ValidateAndGetHeaders(c *gin.Context, requiredHeaders ...string) (*SecureHeaderContext, error) {
	headerCtx := &SecureHeaderContext{}

	// Define header mappings and validation rules
	headerMappings := map[string]struct {
		target   *string
		required bool
	}{
		UserqHeader:   {&headerCtx.UserQ, false},
		HilfmHeader:   {&headerCtx.Hilfm, false},
		BranchHeader:  {&headerCtx.Branch, false},
		OrgechHeader:  {&headerCtx.Orgeh, false},
		StellTXHeader: {&headerCtx.StellTX, false},
		KostlHeader:   {&headerCtx.Kostl, false},
	}

	// Create a map of required headers for quick lookup
	requiredMap := make(map[string]bool)
	for _, header := range requiredHeaders {
		requiredMap[header] = true
	}

	// Validate each header
	for headerName, mapping := range headerMappings {
		headerValue := c.GetHeader(headerName)

		// Check size limit to prevent resource exhaustion
		if len(headerValue) > MaxHeaderSize {
			return nil, fmt.Errorf("%s header exceeds maximum allowed size (%d bytes)", headerName, MaxHeaderSize)
		}

		// Check if required header is missing
		if requiredMap[headerName] && headerValue == "" {
			return nil, fmt.Errorf("%s header is required", headerName)
		}

		// Store the validated header value
		if headerName == UserqHeader {
			// For userq, store as is
			*mapping.target = headerValue
		} else {
			// For other headers, trim leading zeros
			*mapping.target = strings.TrimLeft(headerValue, "0")
		}

	}

	return headerCtx, nil
}

// GetUserInfo safely extracts user information with header validation
func GetUserInfo(c *gin.Context, requiredHeaders ...string) (pernr, name, hilfm, branch, orgeh, kostl string, err error) {
	headerCtx, err := ValidateAndGetHeaders(c, requiredHeaders...)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	// Parse user header only after validation
	if headerCtx.UserQ != "" {
		var parseErr error
		pernr, name, parseErr = helper.ParseUserHeader(headerCtx.UserQ)
		if parseErr != nil {
			return "", "", "", "", "", "", parseErr
		}
	}

	return pernr, name, headerCtx.Hilfm, headerCtx.Branch, headerCtx.Orgeh, headerCtx.Kostl, nil
}

// GetApproverInfo safely extracts approver information with header validation for general services
func GetApproverInfo(c *gin.Context, requiredHeaders ...string) (pernr, name, hilfm, branch, orgeh, jabatan string, err error) {
	headerCtx, err := ValidateAndGetHeaders(c, requiredHeaders...)
	if err != nil {
		return "", "", "", "", "", "", err
	}

	// Parse user header only after validation
	if headerCtx.UserQ != "" {
		var parseErr error
		pernr, name, parseErr = helper.ParseUserHeader(headerCtx.UserQ)
		if parseErr != nil {
			return "", "", "", "", "", "", parseErr
		}
	}

	return pernr, name, headerCtx.Hilfm, headerCtx.Branch, headerCtx.Orgeh, headerCtx.StellTX, nil
}

// AuthenticationResult holds the result of authentication validation
type AuthenticationResult struct {
	Token        string
	Claims       map[string]interface{}
	UserqHeader  string
	HeadersToSet map[string]string
}

// ValidateAuthentication performs comprehensive authentication validation including header security
func ValidateAuthentication(c *gin.Context, jwtVerifyFunc func(string) (*map[string]interface{}, error), fillClaimsFunc func(map[string]interface{}) map[string]string) (*AuthenticationResult, error) {
	// Step 1: Extract and sanitize bearer token before verifier usage.
	token, err := extractAndValidateBearerToken(c)
	if err != nil {
		return nil, err
	}

	// Step 2: Verify JWT token
	claims, err := jwtVerifyFunc(token)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	claimsPayload := *claims

	// Step 3: Build userq header
	pernr, pernrOk := claimsPayload["pernr"].(string)
	nama, namaOk := claimsPayload["nama"].(string)

	if !pernrOk || !namaOk {
		return nil, fmt.Errorf("invalid claims: missing pernr or nama")
	}

	userqHeader := pernr + " | " + nama

	// Security: Validate userq header size
	if len(userqHeader) > MaxHeaderSize {
		return nil, fmt.Errorf("user header exceeds maximum size")
	}

	// Step 4: Process additional headers from claims
	headersToSet := fillClaimsFunc(claimsPayload)

	// Security: Validate all headers to set
	for k, v := range headersToSet {
		if len(k) > MaxHeaderSize {
			return nil, fmt.Errorf("header key '%s' exceeds maximum size", k)
		}
		if len(v) > MaxHeaderSize {
			return nil, fmt.Errorf("header value for '%s' exceeds maximum size", k)
		}
		if strings.ContainsAny(k, "\n\r") || strings.ContainsAny(v, "\n\r") {
			return nil, fmt.Errorf("header contains invalid characters")
		}
	}

	return &AuthenticationResult{
		Token:        token,
		Claims:       claimsPayload,
		UserqHeader:  userqHeader,
		HeadersToSet: headersToSet,
	}, nil
}

// extractAndValidateBearerToken performs strict pre-verification checks so
// untrusted header input cannot force expensive JWT verification work.
func extractAndValidateBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(AuthorizationHeader)

	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}
	if len(authHeader) > MaxHeaderSize {
		return "", fmt.Errorf("authorization header exceeds maximum size")
	}
	if strings.ContainsAny(authHeader, "\n\r") {
		return "", fmt.Errorf("authorization header contains invalid characters")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("bearer token is empty")
	}
	if len(token) > MaxJWTTokenSize {
		return "", fmt.Errorf("bearer token exceeds maximum size")
	}
	if strings.ContainsAny(token, " \t\n\r") {
		return "", fmt.Errorf("bearer token contains invalid whitespace")
	}

	// Require compact JWT syntax and base64url-safe characters only.
	if !jwtTokenPattern.MatchString(token) {
		return "", fmt.Errorf("invalid bearer token format")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid bearer token segments")
	}
	if len(parts[0]) > maxJWTHeaderSegmentSize || len(parts[1]) > maxJWTPayloadSegmentSize || len(parts[2]) > maxJWTSignatureSegmentSize {
		return "", fmt.Errorf("bearer token exceeds maximum size")
	}

	return token, nil
}
