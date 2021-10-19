package models

type Player struct {
	Storage    map[Resource]float64
	Extractors []Extractor
}
