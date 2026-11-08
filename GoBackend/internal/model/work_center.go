package model

type WorkCenter struct {
	ID             int64    `json:"id"`
	Code           string   `json:"code"`
	Title          string   `json:"title"`
	Active         bool     `json:"active"`
	OperationTypes []string `json:"operation_types"`
}

type CreateWorkCenterRequest struct {
	Code           string   `json:"code"            binding:"required,min=1,max=32"`
	Title          string   `json:"title"           binding:"required,min=1,max=100"`
	Active         *bool    `json:"active"`
	OperationTypes []string `json:"operation_types" binding:"required,min=1,max=20,dive,required,min=1,max=20"`
}

func (r CreateWorkCenterRequest) IsActiveOrDefault() bool {
	if r.Active == nil {
		return true
	}
	return *r.Active
}
