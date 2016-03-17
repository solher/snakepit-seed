package models

type Dashboard struct {
	// The document handle. Format: ':collection/:key'
	ID string `json:"_id,omitempty"`
	// The document's revision token. Changes at each update.
	Rev string `json:"_rev,omitempty"`
	// The document's unique key.
	Key  string      `json:"_key,omitempty"`
	Name string      `json:name,omitempty"`
	Spec interface{} `json:spec,omitempty"`
}
