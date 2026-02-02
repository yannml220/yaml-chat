package types

type Handler struct {
	JsonRequestValidator
	JsonResponseSender
	FormRequestValidator
}
