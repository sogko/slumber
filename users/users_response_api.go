package users

//---- User Response API v0 ----

type ListResponse_v0 struct {
	Users   Users  `json:"users,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type CreateResponse_v0 struct {
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type UpdateListResponse_v0 struct {
	Action  string   `json:"action,omitempty"`
	IDs     []string `json:"ids,omitempty"`
	Message string   `json:"message,omitempty"`
	Success bool     `json:"success"`
}

type DeleteAllResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type GetResponse_v0 struct {
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type ConfirmUserResponse_v0 struct {
	Code    string `json:"code,omitempty"`
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type UpdateResponse_v0 struct {
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type DeleteResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}
