package util

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	EmployeeID  uint     `json:"employee_id"`
	ClientID    uint     `json:"client_id"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"`   // "access" | "refresh"
	TokenSource string   `json:"token_source"` // "employee" | "client"
	jwt.RegisteredClaims
}

func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// employeeRoleLevel returns hierarchy level for employee roles.
func employeeRoleLevel(role string) int {
	switch role {
	case "employeeAdmin":
		return 4
	case "employeeSupervisor":
		return 3
	case "employeeAgent":
		return 2
	case "employeeBasic":
		return 1
	default:
		return 0
	}
}

func HasPermission(claims *Claims, perm string) bool {
	requiredLevel := employeeRoleLevel(perm)
	for _, p := range claims.Permissions {
		if p == perm {
			return true
		}
		// Hierarchical: higher employee role grants lower role access
		if requiredLevel > 0 && employeeRoleLevel(p) >= requiredLevel {
			return true
		}
	}
	return false
}

// ValidateAccountNumber returns true if the number is a valid 18-digit account number.
func ValidateAccountNumber(number string) bool {
	if len(number) != 18 {
		return false
	}
	for _, c := range number {
		if c < '0' || c > '9' {
			return false
		}
	}
	bankCode := number[:3]
	return bankCode == "111" || bankCode == "222" || bankCode == "333" || bankCode == "444"
}
