package controller

type Page struct {
	Count int `json:"count"`
	Page  int `json:"page"`
}

type ListResponse struct {
	List interface{} `json:"list"`
	Page `json:"page"`
}

type Filter struct {
	Page   int    `form:"_page"`
	Size   int    `form:"_size"`
}
