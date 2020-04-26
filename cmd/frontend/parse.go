package main

import (
    "strings"
    "net/http"
    "net"
    "io/ioutil"

    "github.com/jsirianni/relay/internal/message"

    "github.com/pkg/errors"
    "github.com/google/uuid"
)

func parseMessage(req *http.Request) ([]byte, error) {
    m := message.New()
    m.SetTime()

    addr, err := parseAddress(req)
    if err != nil {
        return nil, err
    }
    m.SetAddress(addr)

    apiKey, err := parseAPIKey(req)
    if err != nil {
        return nil, err
    }
    if err := m.SetAPIKey(apiKey); err != nil {
        f.Log.Trace("invalid api key:" + apiKey)
        return nil, err
    }

    p, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return nil, err
    }
    if err := m.ParsePayload(p); err != nil {
        return nil, err
    }

    // safely log the message without the APIKey
    safeJson, err := m.BytesSafe()
    if err != nil {
        f.Log.Error(err)
    }
    f.Log.Info("new message: " + string(safeJson))

    // return the message as json
    return m.Bytes()
}

func parseAddress(req *http.Request) (string, error) {
    raw := strings.Split(req.RemoteAddr, ":")[0]
    if raw == "[" {
        raw = "127.0.0.1"
    }
    addr := net.ParseIP(raw)
    if addr == nil {
        f.Log.Trace(errors.Wrap(errors.New(invalidIPError), "failed to parse address from '" + raw + "'"))
        return "", errors.New(invalidIPError)
    }
    return addr.String(), nil
}

func parseAPIKey(req *http.Request) (string, error) {
    apiKey := req.Header.Get(apiKeyHeader)
    if apiKey == "" {
        return "", errors.New(missingAPIKeyHeader)
    }

     _, err := uuid.Parse(apiKey)
    return apiKey, err
}
