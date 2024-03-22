package client

import (
	"errors"
	"firebase.google.com/go/v4/messaging"
	pushererrors "github.com/topfreegames/pusher/errors"
)

// Firebase errors docs can be found here: https://firebase.google.com/docs/cloud-messaging/send-message#admin
var (
	ErrUnspecified         = errors.New("unspecified error")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUnregisteredDevice  = errors.New("unregistered device")
	ErrSenderIDMismatch    = errors.New("sender id mismatch")
	ErrQuotaExceeded       = errors.New("quota exceeded")
	ErrUnavailable         = errors.New("unavailable")
	ErrInternalServerError = errors.New("internal server error")
	ErrThirdParyAuthError  = errors.New("third party authentication error")
)

// TranslateError translates a Firebase error into a pusher error.
func translateError(err error) *pushererrors.PushError {
	switch {
	case messaging.IsInvalidArgument(err):
		return pushererrors.NewPushError("INVALID_JSON", err.Error())
	case messaging.IsUnregistered(err):
		return pushererrors.NewPushError("DEVICE_UNREGISTERED", err.Error())
	case messaging.IsSenderIDMismatch(err):
		return pushererrors.NewPushError("SENDER_ID_MISMATCH", err.Error())
	case messaging.IsQuotaExceeded(err):
		return pushererrors.NewPushError("DEVICE_MESSAGE_RATE_EXCEEDED", err.Error())
	case messaging.IsUnavailable(err):
		return pushererrors.NewPushError("UNAVAILABLE", err.Error())
	case messaging.IsInternal(err):
		return pushererrors.NewPushError("INTERNAL_SERVER_ERROR", err.Error())
	case messaging.IsThirdPartyAuthError(err):
		return pushererrors.NewPushError("THIRD_PARTY_AUTH_ERROR", err.Error())
	default:
		return pushererrors.NewPushError("UNKNOWN", err.Error())
	}

	return nil
}
