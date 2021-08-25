package tests

import (
	"ewallet/sessions"
	"testing"
)

/**
 * test file
 */

func TestSession(t *testing.T) {
	t.Log("start test session...")
	sessionID := sessions.CreateSession("123456")
	t.Log("creat session success.", sessionID)

	t.Log("test get session info...")
	getedSessionID, err := sessions.GetInfo(sessionID, "SessionID")
	if getedSessionID.(string) != sessionID {
		t.Error("error get sessionID", getedSessionID, sessionID)
	}
	if err != nil {
		t.Error("error get Info", err)
	}

	t.Log("test get session info2...")
	sessionLastTime, err := sessions.GetInfo(sessionID, "LastTime")
	if err != nil {
		t.Error("error get Info", err)
	}
	t.Log("test get Session Info", sessionLastTime)
}
