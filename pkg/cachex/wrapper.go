package cachex

import "encoding/json"

type Wrapper struct {
	Exist bool            `json:"exist"`
	Data  json.RawMessage `json:"data"`
}
