package worker

// Worker model
type Worker struct {
	ID    int    `bson:"_id"`
	name  string `json:"name"`
	email string `json:"email"`
}

// Workers is an array of Product objects
type Workers []Worker
