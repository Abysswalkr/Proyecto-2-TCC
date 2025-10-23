package main

// Grammar representa una gramática libre de contexto
type Grammar struct {
	Productions map[string][][]string
	Initial     string
	Terminals   map[string]struct{}
}