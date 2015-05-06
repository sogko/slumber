package users

//---- User Request API v0 ----

type CreateRequest_v0 struct {
	User User `json:"user"`
}

type UpdateListRequest_v0 struct {
	Action string   `json:"action"`
	IDs    []string `json:"ids"`
}

type UpdateRequest_v0 struct {
	User User `json:"user"`
}
