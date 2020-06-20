package auth

import (
    "github.com/google/uuid"
)

type Auth interface{
    ValidAPIKey(uuid.UUID) (bool, error)
    Type() string
}
