package controller

// Action represents one of add, update, or delete
type Action string

const (
	// Add represents the add k8s event
	Add Action = "add"
	// Update represents the update k8s event
	Update Action = "update"
	// Delete represents the delete k8s event
	Delete Action = "delete"
)

// CacheKey combines a cache key with the Action associated with it
type CacheKey struct {
	Key string
	Act Action
}
