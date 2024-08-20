package impls

import "encoding/json"

type Topic interface {
	ExecuteEvent(string, string, []byte) error
}

func Marshal[T any](value []byte) (T, error) {
	var data T
	if err := json.Unmarshal(value, &data); err != nil {
		return data, err
	}
	return data, nil
}
