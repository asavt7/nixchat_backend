package model

// PagedQuery provide params for query result pagination
type PagedQuery struct {
	Size   int `validate:"min=0"`
	Offset int `validate:"min=0"`
}
