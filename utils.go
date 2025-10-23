package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readFile lee un archivo línea por línea y devuelve un slice de cadenas de texto.
func readFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		os.Exit(1)
	}
	defer file.Close()

	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			data = append(data, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	return data
}
