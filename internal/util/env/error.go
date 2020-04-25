package env

import (
    "strings"
)

const notSetERR = "is not set in the environment"

func IsEnvNotSetError(err error) bool {
    if strings.Contains(err.Error(), notSetERR) {
        return true
    }
    return false
}
