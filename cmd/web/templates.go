package main

import "github.com/Bude8/lets-go/internal/models"

// html/template only allows one item of dynamic data when rendering a template
// This allows multiple pieces of dynamic data
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
