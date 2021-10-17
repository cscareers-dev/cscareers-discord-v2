package handlers

// Handler contains custom metadata around a Discord Event
type Handler struct{}

// New returns a Handler
func New() *Handler {
	return &Handler{}
}
