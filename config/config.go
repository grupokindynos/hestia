package config

import (
	"errors"
)

var (
	ErrorNoAuth            = errors.New("you are not authorized")
	ErrorSignature         = errors.New("your signatue could not be parsed")
	ErrorHeaderToken       = errors.New("no token found in header")
	ErrorFbInitializeAuth  = errors.New("unable to initialize auth client")
	ErrorNoBody            = errors.New("body was not found in request")
	ErrorWrongMessage      = errors.New("signed message is not on known hosts")
	ErrorNoHeaderSignature = errors.New("no signature found in header")
	ErrorSignatureParse    = errors.New("could not parse header signature")
	ErrorNoUserInformation = errors.New("unable to get user information")
	ErrorUnmarshal         = errors.New("unable to unmarshal object")
	ErrorMissingID         = errors.New("missing id param")
	ErrorInfoDontMatchUser = errors.New("information requested doesn't match for this user")
	ErrorCoinDataGet       = errors.New("unable to get coin information")
	ErrorInvalidPassword   = errors.New("could not decrypt using master password")
	ErrorConfigDataGet     = errors.New("unable to get config information")
	ErrorDecryptJWE        = errors.New("unable to decrypt jwe")
	ErrorDBStore           = errors.New("unable to store information to database")
	ErrorNotFound          = errors.New("information not found")
	ErrorAlreadyExists     = errors.New("object already exists")
)
