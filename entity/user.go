package entity

import "time"

type User struct {
	UUID       string    `json:"uuid"`
	NameCrypt  []byte    `json:"__encrypted__data_name_crypt"`
	NameHash   []byte    `json:"__encrypted__data_name_hash"`
	EmailCrypt []byte    `json:"__encrypted__data_email_crypt"`
	CreatedAt  time.Time `json:"created_at"`
}
