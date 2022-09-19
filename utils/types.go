package utils

type MMap map[string]map[string]any

type SwaggerJson struct {
	Swagger     string          `json:"swagger"`
	Info        map[string]any  `json:"info"`
	Paths       map[string]MMap `json:"paths"`
	Definitions map[string]any  `json:"definitions"`
}
