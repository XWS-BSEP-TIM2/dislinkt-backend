package ecoding

import (
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base64Coder struct{}

func NewBase64Coder() *Base64Coder {
	return &Base64Coder{}
}

func (c *Base64Coder) Encode(binary primitive.Binary) string {
	return base64.StdEncoding.EncodeToString(binary.Data)
}

func (c *Base64Coder) Decode(str string) (primitive.Binary, error) {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return primitive.Binary{}, err
	}
	return primitive.Binary{Data: decoded}, nil
}
