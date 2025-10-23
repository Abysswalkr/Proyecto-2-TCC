package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TreeNode struct {
	Value    string     `json:"value"`
	Children []*TreeNode `json:"children"`
}

// generateParseTree genera el árbol de análisis sintáctico a partir de la tabla generada por el algoritmo CYK.
func generateParseTree(table [][][]string, grammar map[string][][]string, sentence []string, startSymbol string) *TreeNode {
	n := len(sentence)
	if !contains(table[0][n-1], startSymbol) {
		return nil
	}

	return buildTree(table, grammar, sentence, 0, n-1, startSymbol)
}

// buildTree construye recursivamente el árbol de análisis sintáctico.
func buildTree(table [][][]string, grammar map[string][][]string, sentence []string, i, j int, symbol string) *TreeNode {
	node := &TreeNode{Value: symbol}
	if i == j {
		node.Children = append(node.Children, &TreeNode{Value: sentence[i]})
		return node
	}

	for k := i; k < j; k++ {
		for _, rhs := range grammar[symbol] {
			if len(rhs) == 2 && contains(table[i][k], rhs[0]) && contains(table[k+1][j], rhs[1]) {
				leftChild := buildTree(table, grammar, sentence, i, k, rhs[0])
				rightChild := buildTree(table, grammar, sentence, k+1, j, rhs[1])
				node.Children = append(node.Children, leftChild, rightChild)
				return node
			}
		}
	}
	return node
}

// printTree imprime el árbol de análisis sintáctico en una forma jerárquica.
func printTree(node *TreeNode, level int) {
	if node == nil {
		return
	}
	fmt.Printf("%s%s\n", strings.Repeat("  ", level), node.Value)
	for _, child := range node.Children {
		printTree(child, level+1)
	}
}

// saveTreeAsJSON guarda el árbol sintáctico como un archivo JSON.
func saveTreeAsJSON(root *TreeNode, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(root)
}
