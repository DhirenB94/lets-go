package main

import models "dhiren.brahmbhatt/snippetbox/pkg"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}