package main

import (
	"fmt"
	"os"
	"strconv"
)

// Converts a Context-Free Grammar into Chomsky Normal Form
func from_cfg_to_cnf(cfg *Grammar) *Grammar {
	fmt.Fprintln(os.Stdout, "\n════════════════════════════════════════")
	fmt.Fprintln(os.Stdout, "🔄  Convertir CFG a CNF")
	fmt.Fprintln(os.Stdout, "════════════════════════════════════════")

	// Paso 1: Agregar nuevo símbolo inicial
	fmt.Fprintf(os.Stdout, "🚀  Agregando nuevo símbolo inicial\n")
	cnf := add_new_initial(cfg)

	// Paso 2: Eliminar producciones epsilon
	fmt.Fprintln(os.Stdout, "🔧  Eliminando producciones epsilon...")
	cnf = eliminate_epsilon_productions(cnf)

	// Paso 3: Eliminar producciones unitarias
	fmt.Fprintln(os.Stdout, "🔧  Eliminando producciones unitarias...")
	cnf = eliminate_unit_productions(cnf)

	// Paso 4: Binarizar producciones
	fmt.Fprintln(os.Stdout, "🔧  Binarizando producciones...")
	cnf = binarize_productions(cnf)

	// Paso 5: Eliminar símbolos inútiles
	fmt.Fprintln(os.Stdout, "🔧  Eliminando símbolos inútiles...")
	cnf = remove_useless_symbols(cnf)

	fmt.Fprintln(os.Stdout, "✔️  Conversión a CNF completada.")
	fmt.Fprintln(os.Stdout, "════════════════════════════════════════")
	
	// Debug: imprimir gramática resultante
	fmt.Fprintln(os.Stdout, "📊 Gramática en CNF:")
	for nonTerminal, productions := range cnf.Productions {
		for _, production := range productions {
			fmt.Fprintf(os.Stdout, "  %s -> %v\n", nonTerminal, production)
		}
	}
	fmt.Fprintln(os.Stdout, "════════════════════════════════════════")
	
	return cnf
}

// Agrega nuevo símbolo inicial S'
func add_new_initial(cfg *Grammar) *Grammar {
	newInitial := "S'"
	// Solo agregar si no existe ya
	if _, exists := cfg.Productions[newInitial]; !exists {
		cfg.Productions[newInitial] = [][]string{{cfg.Initial}}
		cfg.Initial = newInitial
	}
	return cfg
}

// Elimina producciones epsilon (producciones vacías)
func eliminate_epsilon_productions(cfg *Grammar) *Grammar {
	// Paso 1: Encontrar todos los símbolos anulables (nullable)
	nullable := make(map[string]bool)
	changed := true
	
	// Inicializar: los símbolos que tienen producción epsilon directa
	for nonTerminal, productions := range cfg.Productions {
		for _, production := range productions {
			if len(production) == 1 && production[0] == "ε" {
				nullable[nonTerminal] = true
			}
		}
	}
	
	// Propagar anulabilidad
	for changed {
		changed = false
		for nonTerminal, productions := range cfg.Productions {
			if nullable[nonTerminal] {
				continue
			}
			for _, production := range productions {
				allNullable := true
				for _, symbol := range production {
					if symbol == "ε" {
						continue
					}
					if !nullable[symbol] {
						allNullable = false
						break
					}
				}
				if allNullable && len(production) > 0 {
					if !nullable[nonTerminal] {
						nullable[nonTerminal] = true
						changed = true
					}
				}
			}
		}
	}
	
	// Paso 2: Generar nuevas producciones sin epsilon
	newProductions := make(map[string][][]string)
	
	for nonTerminal, productions := range cfg.Productions {
		newProductions[nonTerminal] = [][]string{}
		
		for _, production := range productions {
			if len(production) == 1 && production[0] == "ε" {
				// Saltar producciones epsilon explícitas
				continue
			}
			
			// Generar todas las combinaciones posibles omitiendo símbolos anulables
			combinations := generate_combinations(production, nullable)
			
			for _, comb := range combinations {
				if len(comb) > 0 { // No agregar producciones vacías
					newProductions[nonTerminal] = append(newProductions[nonTerminal], comb)
				}
			}
		}
	}
	
	cfg.Productions = newProductions
	return cfg
}

// Genera combinaciones omitiendo símbolos anulables
func generate_combinations(production []string, nullable map[string]bool) [][]string {
	if len(production) == 0 {
		return [][]string{{}}
	}
	
	first := production[0]
	rest := production[1:]
	
	restCombinations := generate_combinations(rest, nullable)
	var result [][]string
	
	// Incluir el primer símbolo
	for _, comb := range restCombinations {
		newComb := append([]string{first}, comb...)
		result = append(result, newComb)
	}
	
	// Omitir el primer símbolo si es anulable
	if nullable[first] {
		result = append(result, restCombinations...)
	}
	
	return result
}

// Elimina producciones unitarias (A -> B)
func eliminate_unit_productions(cfg *Grammar) *Grammar {
	// Para cada no terminal, calcular su cierre unitario
	unitClosures := make(map[string]map[string]bool)
	
	for nonTerminal := range cfg.Productions {
		closure := make(map[string]bool)
		closure[nonTerminal] = true
		changed := true
		
		for changed {
			changed = false
			for A := range closure {
				for _, production := range cfg.Productions[A] {
					if len(production) == 1 {
						B := production[0]
						if _, isNonTerminal := cfg.Productions[B]; isNonTerminal {
							if !closure[B] {
								closure[B] = true
								changed = true
							}
						}
					}
				}
			}
		}
		unitClosures[nonTerminal] = closure
	}
	
	// Construir nuevas producciones
	newProductions := make(map[string][][]string)
	
	for A, closure := range unitClosures {
		newProductions[A] = [][]string{}
		
		for B := range closure {
			if A == B {
				continue // Saltar autorreferencias
			}
			
			// Agregar todas las producciones no unitarias de B
			for _, production := range cfg.Productions[B] {
				if len(production) != 1 || (len(production) == 1 && cfg.Productions[production[0]] == nil) {
					// Solo agregar si no es producción unitaria o es terminal
					exists := false
					for _, existing := range newProductions[A] {
						if equalSlices(existing, production) {
							exists = true
							break
						}
					}
					if !exists {
						newProductions[A] = append(newProductions[A], production)
					}
				}
			}
		}
		
		// También mantener las producciones originales no unitarias de A
		for _, production := range cfg.Productions[A] {
			if len(production) != 1 || (len(production) == 1 && cfg.Productions[production[0]] == nil) {
				exists := false
				for _, existing := range newProductions[A] {
					if equalSlices(existing, production) {
						exists = true
						break
					}
				}
				if !exists {
					newProductions[A] = append(newProductions[A], production)
				}
			}
		}
	}
	
	cfg.Productions = newProductions
	return cfg
}

// Binariza producciones (convierte a máximo 2 símbolos por producción)
func binarize_productions(cfg *Grammar) *Grammar {
	generator := construct_generator()
	newProductions := make(map[string][][]string)
	
	// Primero copiar todas las producciones existentes
	for nonTerminal, productions := range cfg.Productions {
		newProductions[nonTerminal] = [][]string{}
		for _, production := range productions {
			if len(production) <= 2 {
				// Ya está en forma binaria o menos
				newProductions[nonTerminal] = append(newProductions[nonTerminal], production)
			} else {
				// Necesita binarización
				currentSymbol := nonTerminal
				for i := 0; i < len(production)-2; i++ {
					newSymbol := generator()
					
					// Crear nueva producción
					newProductions[currentSymbol] = append(newProductions[currentSymbol], 
						[]string{production[i], newSymbol})
					
					currentSymbol = newSymbol
				}
				// Última producción
				newProductions[currentSymbol] = append(newProductions[currentSymbol], 
					production[len(production)-2:])
			}
		}
	}
	
	cfg.Productions = newProductions
	return cfg
}

// Elimina símbolos inútiles (no alcanzables o no generativos)
func remove_useless_symbols(cfg *Grammar) *Grammar {
	// Paso 1: Eliminar símbolos no generativos
	generative := make(map[string]bool)
	
	// Inicializar: los terminales son generativos
	for terminal := range cfg.Terminals {
		generative[terminal] = true
	}
	
	changed := true
	for changed {
		changed = false
		for nonTerminal, productions := range cfg.Productions {
			if generative[nonTerminal] {
				continue
			}
			for _, production := range productions {
				allGenerative := true
				for _, symbol := range production {
					if !generative[symbol] {
						allGenerative = false
						break
					}
				}
				if allGenerative {
					generative[nonTerminal] = true
					changed = true
					break
				}
			}
		}
	}
	
	// Eliminar símbolos no generativos
	newProductions := make(map[string][][]string)
	for nonTerminal, productions := range cfg.Productions {
		if generative[nonTerminal] {
			newProductions[nonTerminal] = [][]string{}
			for _, production := range productions {
				allGenerative := true
				for _, symbol := range production {
					if !generative[symbol] {
						allGenerative = false
						break
					}
				}
				if allGenerative {
					newProductions[nonTerminal] = append(newProductions[nonTerminal], production)
				}
			}
		}
	}
	cfg.Productions = newProductions
	
	// Paso 2: Eliminar símbolos no alcanzables
	reachable := make(map[string]bool)
	reachable[cfg.Initial] = true
	changed = true
	
	for changed {
		changed = false
		for nonTerminal, productions := range cfg.Productions {
			if !reachable[nonTerminal] {
				continue
			}
			for _, production := range productions {
				for _, symbol := range production {
					if _, isNonTerminal := cfg.Productions[symbol]; isNonTerminal {
						if !reachable[symbol] {
							reachable[symbol] = true
							changed = true
						}
					}
				}
			}
		}
	}
	
	// Eliminar símbolos no alcanzables
	finalProductions := make(map[string][][]string)
	for nonTerminal, productions := range cfg.Productions {
		if reachable[nonTerminal] {
			finalProductions[nonTerminal] = productions
		}
	}
	cfg.Productions = finalProductions
	
	return cfg
}

// Genera nombres únicos para nuevos no terminales
func construct_generator() func() string {
	count := 0
	return func() string {
		count++
		return "X" + strconv.Itoa(count)
	}
}

// Compara si dos slices son iguales
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}