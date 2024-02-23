package firestoreState

type FirestoreState struct {
	Brightness float32 `json:"brightness"`
	Color      string  `json:"color"`
	State      string  `json:"state"`
	Sync       bool    `json:"sync"`
	Transition float32 `json:"transition"`
	Write      bool    `json:"write"`
}
