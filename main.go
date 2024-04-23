package main

import (
	"fmt"
	"os"
)

type BracketHandler interface {
	IsOpenBracket(rune) bool
	GetClosingBracket(rune) (rune, bool)
	IsMatchingClosingBracket(rune, rune) bool
}

type CodeBracketHandler struct {
	brackets map[rune]rune
}

func (h *CodeBracketHandler) IsOpenBracket(char rune) bool {
	_, ok := h.brackets[char]
	return ok
}

func (h *CodeBracketHandler) GetClosingBracket(char rune) (rune, bool) {
	closing, ok := h.brackets[char]
	return closing, ok
}

func (h *CodeBracketHandler) IsMatchingClosingBracket(open rune, close rune) bool {
	return h.brackets[open] == close
}

func NewCodeBracketHandler() *CodeBracketHandler {
	return &CodeBracketHandler{
		brackets: map[rune]rune{
			'(': ')',
			'{': '}',
			'[': ']',
		},
	}
}

func validateBrackets(text string, handler BracketHandler, filePath string) []string {
	var stack []rune
	var errors []string

	for i, char := range text {
		lineNumber := i + 1
		columnNumber := i + 1

		if handler.IsOpenBracket(char) {
			stack = append(stack, char)
		} else if closing, ok := handler.GetClosingBracket(char); ok {
			if len(stack) == 0 || !handler.IsMatchingClosingBracket(stack[len(stack)-1], closing) {
				errors = append(errors, fmt.Sprintf("Invalid bracket %c found at %s:%d:%d", char, filePath, lineNumber, columnNumber))
			} else {
				stack = stack[:len(stack)-1]
			}
		}
	}

	for _, openBracket := range stack {
		lineNumber := len(text) + 1
		columnNumber := 1
		errors = append(errors, fmt.Sprintf("Invalid bracket %c found at %s:%d:%d", openBracket, filePath, lineNumber, columnNumber))
	}

	return errors
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: invalid_brackets <file_path>")
		return
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	text := string(data)
	handler := NewCodeBracketHandler()
	errors := validateBrackets(text, handler, filePath)

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}

