// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

type User struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Verified        bool   `json:"verified"`
	IsAdmin         bool   `json:"is_admin"`
	EncryptedWallet string `json:"encrypted_wallet"`
	Passwhash       string `json:"passwhash"`
}
