package apperrors

type RenderType string

const (
	RenderTypeToast          RenderType = "toast"
	RenderTypeModal          RenderType = "modal"
	RenderTypeLogin          RenderType = "login"
	RenderTypeInvalidLicense RenderType = "invalidLicense"
)
