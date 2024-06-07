package dict

import (
	"encoding/json"
	"io"
)

func UnmarshalJSON(p []byte) (Dict, error) {
	d := make(Dict)
	if err := json.Unmarshal(p, &d); err != nil {
		return nil, err
	}
	return d, nil
}

func UnmarshalJSONReader(r io.Reader) (Dict, error) {
	d := make(Dict)
	if err := json.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}
	return d, nil
}
