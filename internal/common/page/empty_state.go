package page

import "fmt"

type EmptyState struct {
	Icon        string
	Title       string
	Description string
	ActionHref  string
	ActionLabel string
}

func NewEmptyState(icon, entityLabel, addHref, addLabel string) EmptyState {
	return EmptyState{
		Icon:        icon,
		Title:       "No se encontraron resultados",
		Description: fmt.Sprintf("No hay %s para mostrar. Probá ajustar la búsqueda o agregá uno nuevo.", entityLabel),
		ActionHref:  addHref,
		ActionLabel: addLabel,
	}
}

func NewNoResultsEmptyState(icon, entityLabel string) EmptyState {
	return EmptyState{
		Icon:        icon,
		Title:       "No se encontraron resultados",
		Description: fmt.Sprintf("No hay %s que coincidan con la búsqueda.", entityLabel),
	}
}
