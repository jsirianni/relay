package env

import (
    "os"
    "testing"
)

// fn type represents any function that returns a string
// and an error (all of the environment functions in env.go)
// in order to write reusable tests, see TestEnv() and testEnv()
type fn func() (string, error)

type TestCase struct {
    // function to call
    fn

    // environment variable to use
    envVar string

    // value to set in the environment that we expect
    // each function to return
    expect string
}

func TestConstants(t *testing.T) {
    if ENVLogLevel != "RELAY_LOG_LEVEL" {
        t.Errorf("envLogLevel should be RELAY_LOG_LEVEL")
    }
    if ENVTopic !=  "RELAY_TOPIC" {
        t.Errorf("envTopic should be RELAY_TOPIC")
    }
    if ENVSub != "RELAY_SUBSCRIPTION" {
        t.Errorf("envSub should be RELAY_SUBSCRIPTION")
    }
    if ENVFrontPort != "RELAY_FRONTEND_PORT" {
        t.Errorf("envFrontPort should be RELAY_FRONTEND_PORT")
    }
    if ENVQueueType != "RELAY_QUEUE_TYPE" {
        t.Errorf("envQueueType should be RELAY_QUEUE_TYPE")
    }
    if ENVGoogleProjectID != "RELAY_GOOGLE_PROJECT_ID" {
        t.Errorf("envProjectID should be RELAY_GOOGLE_PROJECT_ID")
    }
}

// use strings for defining the environment variable names instead of
// using the constants defined in env.go
func TestEnv(t *testing.T) {
    c := []TestCase{}

    c = append(c, TestCase{fn: FrontendPort, envVar: "RELAY_FRONTEND_PORT", expect: ""})
    c = append(c, TestCase{fn: FrontendPort, envVar: "RELAY_FRONTEND_PORT", expect: "test"})

    c = append(c, TestCase{fn: GoogleProjectID, envVar: "RELAY_GOOGLE_PROJECT_ID", expect: ""})
    c = append(c, TestCase{fn: GoogleProjectID, envVar: "RELAY_GOOGLE_PROJECT_ID", expect: "test"})

    c = append(c, TestCase{fn: LogLevel, envVar: "RELAY_LOG_LEVEL", expect: ""})
    c = append(c, TestCase{fn: LogLevel, envVar: "RELAY_LOG_LEVEL", expect: "test"})

    c = append(c, TestCase{fn: QueueType, envVar: "RELAY_QUEUE_TYPE", expect: ""})
    c = append(c, TestCase{fn: QueueType, envVar: "RELAY_QUEUE_TYPE", expect: "test"})

    c = append(c, TestCase{fn: Subscription, envVar: "RELAY_SUBSCRIPTION", expect: ""})
    c = append(c, TestCase{fn: Subscription, envVar: "RELAY_SUBSCRIPTION", expect: "test"})

    c = append(c, TestCase{fn: Topic, envVar: "RELAY_TOPIC", expect: ""})
    c = append(c, TestCase{fn: Topic, envVar: "RELAY_TOPIC", expect: "test"})

    for _, c := range c {
        os.Setenv(c.envVar, c.expect)
        testEnv(c, t)
    }
}

// testEnv takes a TestCase and checks the result
func testEnv(c TestCase, t *testing.T) {
    output, err := c.fn()

    // if expect is an empty string, we should get an env is not
    // set error
    if c.expect == "" && err == nil {
        t.Errorf(expectedErrorGotNil(c.envVar))
    // if expect is not empty we expect a nil error
    } else if c.expect != "" && err != nil {
        t.Errorf(expectedNilGotError(c.envVar, err))
    }
    // if output is not equal to expect, we set the wrong environment
    // variable or value
    if output != c.expect {
        t.Errorf(expectedStringGotString(c.envVar, c.expect, output))
    }
}

func expectedErrorGotNil(prefix string) string {
    return prefix+": expected an error when environment is not set, got nil"
}

func expectedNilGotError(prefix string, err error) string {
    return prefix+": expected a nil error when environment is set, got: " + err.Error()
}

func expectedStringGotString(prefix, expect, got string) string {
    return prefix+": expected '"+expect+"' when environment is set to '"+expect+"', got: '"+got+"'"
}
