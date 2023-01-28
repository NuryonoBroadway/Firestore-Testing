package helper

import "time"

type City struct {
	Name       string    `json:"name" firestore:"name,omitempty"`
	State      string    `json:"state" firestore:"state,omitempty"`
	Country    string    `json:"country" firestore:"country,omitempty"`
	Capital    bool      `json:"capital" firestore:"capital"`
	Population int       `json:"population" firestore:"population"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at" firestore:"deleted_at,omitempty"`
}
