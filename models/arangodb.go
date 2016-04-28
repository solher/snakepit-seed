package models

type Document struct {
	// The document handle. Format: ':collection/:key'
	ID string `json:"_id,omitempty"`
	// The document's revision token. Changes at each update.
	Rev string `json:"_rev,omitempty"`
	// The document's unique key.
	Key string `json:"_key,omitempty"`
}

func NewDocument(id, rev, key string) Document {
	return Document{
		ID:  id,
		Rev: rev,
		Key: key,
	}
}

type Edge struct {
	Document
	// Reference to another document. Format: ':collection/:key'
	From string `json:"_from,omitempty"`
	// Reference to another document. Format: ':collection/:key'
	To string `json:"_to,omitempty"`
}

func NewEdge(id, rev, key, from, to string) Edge {
	return Edge{
		Document: NewDocument(id, rev, key),
		From:     from,
		To:       to,
	}
}
