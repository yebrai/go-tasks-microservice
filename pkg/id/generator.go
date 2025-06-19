package id

// Generator define el contrato para generar IDs únicos
type Generator interface {
	Generate() string
}
