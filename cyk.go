package main

import (
	"fmt"
	"time"
)

// cykParse ejecuta el algoritmo CYK y mide el tiempo de ejecución.
func cykParse(grammar map[string][][]string, sentence []string) (bool, [][][]string) {
	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("🔍  Algoritmo CYK - Análisis de la frase")
	fmt.Println("════════════════════════════════════════")

	start := time.Now() // Inicio de la medición de tiempo

	n := len(sentence)
	T := make([][][]string, n)
	for i := range T {
		T[i] = make([][]string, n)
	}

	// DEBUG: Mostrar la frase completa
	fmt.Printf("📝 Frase a analizar: %v\n", sentence)
	fmt.Printf("🔢 Longitud de la frase: %d símbolos\n", n)

	// Inicialización
	for j := 0; j < n; j++ {
		fmt.Printf("\n🔎 Procesando símbolo %d: '%s'\n", j, sentence[j])
		for lhs, productions := range grammar {
			for _, rhs := range productions {
				if len(rhs) == 1 && rhs[0] == sentence[j] {
					T[j][j] = append(T[j][j], lhs)
					fmt.Printf("   ✅ Terminal encontrado: %s -> %s en T[%d][%d]\n", lhs, rhs[0], j, j)
				}
			}
		}
		if len(T[j][j]) == 0 {
			fmt.Printf("   ⚠️  No se encontraron producciones para '%s'\n", sentence[j])
		}
	}

	// Llenado de la tabla CYK
	for span := 2; span <= n; span++ {
		for i := 0; i <= n-span; i++ {
			j := i + span - 1
			fmt.Printf("\n🔍 Analizando rango T[%d][%d] (símbolos %d a %d)\n", i, j, i, j)

			for k := i; k < j; k++ {
				for lhs, productions := range grammar {
					for _, rhs := range productions {
						if len(rhs) == 2 {
							B, C := rhs[0], rhs[1]
							if contains(T[i][k], B) && contains(T[k+1][j], C) {
								if !contains(T[i][j], lhs) {
									T[i][j] = append(T[i][j], lhs)
									fmt.Printf("   🔗 %s -> %s %s en T[%d][%d]\n", lhs, B, C, i, j)
								}
							}
						}
					}
				}
			}
			if len(T[i][j]) > 0 {
				fmt.Printf("   ✅ T[%d][%d] final = %v\n", i, j, T[i][j])
			}
		}
	}

	printTable(T, sentence)

	// ✅ NUEVA VERIFICACIÓN GENERALIZADA
	finalSymbols := T[0][n-1]
	accepted := contains(finalSymbols, "S'") || contains(finalSymbols, "S") || contains(finalSymbols, "E")

	fmt.Printf("\n🎯 Resultado final:\n")
	fmt.Printf("   T[0][%d] = %v\n", n-1, finalSymbols)
	fmt.Printf("   ¿Contiene S', S o E? %v\n", accepted)

	elapsed := time.Since(start)
	fmt.Printf("⏱️  Tiempo de ejecución del algoritmo CYK: %s\n", elapsed)
	fmt.Println("════════════════════════════════════════")

	return accepted, T
}

// printTable imprime la tabla de CYK para fines de depuración.
func printTable(T [][][]string, sentence []string) {
	fmt.Println("\n╔══════════════════════════════════╗")
	fmt.Println("║           Tabla de CYK           ║")
	fmt.Println("╚══════════════════════════════════╝")
	n := len(sentence)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if len(T[i][j]) > 0 {
				fmt.Printf("T[%d][%d]: %v\n", i, j, T[i][j])
			}
		}
	}
	fmt.Println("════════════════════════════════════════")
}

// contains verifica si un slice contiene un elemento específico.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}