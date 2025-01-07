package model

type Category struct {
	Id          int    `json:"id" orm:"id,primary"`
	Name        string `json:"name" orm:"name"`
	Description string `json:"description" orm:"description"`
	ParentId    int    `json:"parent_id" orm:"parent_id"`
}
