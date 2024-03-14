package models

type Object = map[string]interface{}
type StrObject = map[string]string
type SubObject = map[string]Object

type Result struct {
	Item Object   `json:"item,omitempty"`
	List []Object `json:"list,omitempty"`
}
