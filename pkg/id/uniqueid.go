package id

import "github.com/google/uuid"

// UniqueIDGenerator genera IDs únicos usando UUID
type UniqueIDGenerator struct{}

// NewUniqueIDGenerator crea un nuevo generador de UUIDs
func NewUniqueIDGenerator() *UniqueIDGenerator {
	return &UniqueIDGenerator{}
}

// Generate genera un nuevo UUID como string
func (g *UniqueIDGenerator) Generate() string {
	return uuid.New().String()
}
