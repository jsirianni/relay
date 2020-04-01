package auth

type Auth interface{
    ValidAPIKey(string) (bool, error)
    Type() string
}
