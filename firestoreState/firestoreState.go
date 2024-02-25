package firestoreState

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FirestoreState struct {
	Brightness float32 `json:"brightness"`
	Color      string  `json:"color"`
	State      string  `json:"state"`
	Sync       bool    `json:"sync"`
	Transition float32 `json:"transition"`
	Write      bool    `json:"write"`
}

func GetDeviceState(deviceName string, collection *firestore.CollectionRef) (FirestoreState, error) {
    var state FirestoreState
	bgCtx := context.Background()

	snapshot, err := collection.Doc(deviceName).Get(bgCtx)
	if err != nil {
		return state, err
	}

	err = snapshot.DataTo(&state)
	if err != nil {
		return state, err
	}
	return state, nil
}
