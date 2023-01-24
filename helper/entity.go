package helper

type City struct {
	Name       string `json:"name" firestore:"name,omitempty"`
	State      string `json:"state" firestore:"state,omitempty"`
	Country    string `json:"country" firestore:"country,omitempty"`
	Capital    bool   `json:"capital" firestore:"capital"`
	Population int    `json:"population" firestore:"population"`
}
