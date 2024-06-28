package auth

import "golang.org/x/crypto/bcrypt"

// Check if password is valid
func CheckPasswordHash(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
