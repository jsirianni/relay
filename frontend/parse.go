package main

import (
    "strings"
    "net/http"
    "net"
    "encoding/json"

    "github.com/jsirianni/relay/common/message"

    "github.com/pkg/errors"
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
        p.Log.Trace("invalid api key:" + apiKey)
        return nil, err
    }

    i, err := parsePayload(req)
    if err != nil {
        p.Log.Trace(errors.Wrap(err, "failed to parse json body"))
        return nil, err
    }
    m.SetText(i.Text)

    safeJson, err := m.BytesSafe()
    if err != nil {
        p.Log.Error(err)
    }
    p.Log.Info("new message: " + string(safeJson))

    return m.Bytes()
}

func parseAddress(req *http.Request) (string, error) {
    raw := strings.Split(req.RemoteAddr, ":")[0]
    if raw == "[" {
        raw = "127.0.0.1"
    }
    addr := net.ParseIP(raw)
    if addr == nil {
        p.Log.Trace(errors.Wrap(errors.New(invalidIPError), "failed to parse address from '" + raw + "'"))
        return "", errors.New(invalidIPError)
    }
    return addr.String(), nil
}

func parseAPIKey(req *http.Request) (string, error) {
    apiKey := req.Header.Get(apiKeyHeader)
    if apiKey == "" {
        p.Log.Trace(missingAPIKeyHeader)
        return "", errors.New(missingAPIKeyHeader)
    }
    return apiKey, nil
}

func parsePayload(req *http.Request) (IncomingRequest, error) {
    d := json.NewDecoder(req.Body)
    i := IncomingRequest{}
    err := d.Decode(&i)
    return i, err
}
