package handlers

// Handler describes what actions to take for a given event
type Handler interface {
	ObjectAdded(obj interface{}) error
	ObjectDeleted(obj interface{}) error
	ObjectUpdated(ojb interface{}) error
}

// Default is an empty Handler
type Default struct{}

// ObjectAdded does nothing
func (d *Default) ObjectAdded(obj interface{}) error {
	return nil
}

// ObjectDeleted does nothing
func (d *Default) ObjectDeleted(obj interface{}) error {
	return nil
}

// ObjectUpdated does nothing
func (d *Default) ObjectUpdated(obj interface{}) error {
	return nil
}
