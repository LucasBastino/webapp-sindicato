package apperrors

import (
	"net/http"
)

func NewUnauthorizedError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "Tu sesión no es válida o ha expirado."
	}
	return NewAppError(
		err,
		"unauthorized",
		clientMsg,
		http.StatusUnauthorized,
		RenderTypeLogin,
	)
}

func NewForbiddenError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "No tiene permisos para realizar esta acción."
	}
	return NewAppError(
		err,
		"forbidden",
		clientMsg,
		http.StatusForbidden,
		RenderTypeToast,
	)
}

func NewDatabaseError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "Ocurrió un error interno, intentelo nuevamente en unos segundos."
	}
	return NewAppError(
		err,
		"database",
		clientMsg,
		http.StatusServiceUnavailable,
		RenderTypeToast,
	)
}

func NewBusinessError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "No se puede realizar la operación solicitada."
	}
	return NewAppError(
		err,
		"business",
		clientMsg,
		http.StatusConflict,
		RenderTypeToast,
	)
}

func NewNotFoundError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "El registro al que desea acceder no existe."
	}
	return NewAppError(
		err,
		"not_found",
		clientMsg,
		http.StatusNotFound,
		RenderTypeModal,
	)
}

func NewBadRequestError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "Solicitud inválida."
	}
	return NewAppError(
		err,
		"bad_request",
		clientMsg,
		http.StatusBadRequest,
		RenderTypeToast,
	)
}

func NewInternalError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "Ocurrió un error interno en el servidor, intentelo nuevamente en unos segundos."
	}
	return NewAppError(
		err,
		"internal",
		clientMsg,
		http.StatusInternalServerError,
		RenderTypeModal,
	)
}

func NewInvalidLicenseError(err error, clientMsg string) *AppError {
	if clientMsg == "" {
		clientMsg = "La licencia no es válida."
	}
	return NewAppError(
		err,
		"forbidden",
		clientMsg,
		http.StatusForbidden,
		RenderTypeInvalidLicense,
	)
}