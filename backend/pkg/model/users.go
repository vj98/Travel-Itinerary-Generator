package model

// User represents a user in the system.
type Users struct {
	ID            string `bson:"_id,omitempty"`
	Email         string `bson:"email"`
	Password      string `bson:"password"`
	Status        string `bson:"status"`
	Session       string `bson:"session"`
	Type          string `bson:"type"`
	LastUpdatedAt string `bson:"lastupdatedat"`
}
