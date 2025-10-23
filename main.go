package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("✨ Bienvenido al Verificador de Frases ✨")
	fmt.Println("════════════════════════════════════════")

	fmt.Println("\n📚 Seleccione la gramática a usar:")
	fmt.Println("1. Gramática de inglés (input.txt)")
	fmt.Println("2. Gramática aritmética CNF (1-cnf.txt)")
	fmt.Print("Ingrese opción (1-2): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	opcion := strings.TrimSpace(scanner.Text())

	var grammarFile string
	var grammarName string

	switch opcion {
	case "1":
		grammarFile = "input.txt"
		grammarName = "inglés"
	case "2":
		grammarFile = "1-cnf.txt"
		grammarName = "aritmética CNF"
	default:
		grammarFile = "input.txt"
		grammarName = "inglés (por defecto)"
		fmt.Println("⚠️  Opción inválida, usando gramática por defecto")
	}

	if _, err := os.Stat(grammarFile); os.IsNotExist(err) {
		fmt.Printf("❌ El archivo '%s' no existe.\n", grammarFile)
		return
	}

	data := readFile(grammarFile)
	fmt.Printf("\n📖 Cargando reglas de la gramática de %s...\n", grammarName)

	rules := make(map[string][][]string)
	for _, line := range data {
		first, rest, ok := strings.Cut(line, "->")
		if !ok {
			fmt.Printf("⚠️  Regla inválida: %s\n", line)
			continue
		}

		first = strings.TrimSpace(first)
		transitions := strings.Split(strings.TrimSpace(rest), "|")

		for _, transition := range transitions {
			states := strings.Fields(strings.TrimSpace(transition))
			rules[first] = append(rules[first], states)
			fmt.Printf("✅  %s -> %s\n", first, strings.Join(states, " "))
		}
	}

	terminals := map[string]struct{}{}
	for _, prods := range rules {
		for _, prod := range prods {
			for _, sym := range prod {
				if _, isNonTerm := rules[sym]; !isNonTerm {
					terminals[sym] = struct{}{}
				}
			}
		}
	}

	initial := "S"
	if _, exists := rules["E"]; exists {
		initial = "E"
	}

	grammar := Grammar{Productions: rules, Terminals: terminals, Initial: initial}

	// Ya no convertimos a CNF, solo usamos las gramáticas válidas directamente
	fmt.Printf("\n🚀 Gramática cargada con %d no terminales\n", len(grammar.Productions))

	fmt.Println("\n💬 Ingrese la frase o expresión a analizar:")
	switch grammarFile {
	case "input.txt":
		fmt.Println("   Ejemplo: she eats a cake with a fork")
	case "1-cnf.txt":
		fmt.Println("   Ejemplo: id + id * id  o  ( id + id )")
	}
	fmt.Print("   > ")

	scanner.Scan()
	inputText := strings.TrimSpace(scanner.Text())
	if inputText == "" {
		fmt.Println("❌ No se ingresó ninguna frase.")
		return
	}

	sentence := strings.Fields(inputText)
	accepted, table := cykParse(grammar.Productions, sentence)

	if accepted {
		fmt.Println("✅ La frase pertenece al lenguaje.")
		tree := generateParseTree(table, grammar.Productions, sentence, grammar.Initial)
		if tree != nil {
			fmt.Println("\n🌳 Árbol sintáctico:")
			printTree(tree, 0)

			if err := saveTreeAsJSON(tree, "output/tree.json"); err != nil {
				fmt.Printf("⚠️  Error al guardar JSON: %v\n", err)
			} else {
				fmt.Println("📁 Árbol guardado en 'output/tree.json'")
			}
		} else {
			fmt.Println("⚠️  No se pudo construir el árbol sintáctico.")
		}
	} else {
		fmt.Println("❌ La frase NO pertenece al lenguaje.")
	}

	fmt.Println("════════════════════════════════════════")
}