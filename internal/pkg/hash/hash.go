// Package hash инкапсулирует хеширование и проверку паролей (bcrypt).
package hash

import "golang.org/x/crypto/bcrypt"

// Hash возвращает bcrypt-хеш пароля.
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Check сравнивает пароль с ранее сохранённым хешем.
func Check(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
