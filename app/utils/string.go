package utils

import (
	"strconv"

	"github.com/google/uuid"
)

func GenerateUuid() string {
	uuid := uuid.New().String()

	return uuid
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
