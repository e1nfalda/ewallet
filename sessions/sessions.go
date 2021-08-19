package sessions

/*
 * simple implement session
 * Assume system is a single node service, so just use memory as a
 * Features:
 *   * support multi devices login by using unique session id. Avoid relogin on other devices error
 */

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const MAX_LIFETIME_SECONDS = 60 * 60 * time.Second

type userSession struct {
	SessionID string
	Phone     string
	LastTime  time.Time // Last update time. any access may change this value
	// infos like login ip, login time things here
}

type sessionManager struct {
	Name     string
	Sessions map[string]*userSession
}

var globalSession *sessionManager

// GetInfo get Session Info by attr. If not exists, return nil
func GetInfo(sessionID string, attr string, byUser ...interface{}) (interface{}, error) {
	if sessionExpired(sessionID) {
		return nil, errors.New("session Expired")
	}

	r := reflect.ValueOf(globalSession.Sessions[sessionID])
	v := reflect.Indirect(r).FieldByName(attr).Interface()

	// flag hints access by user action
	if len(byUser) == 1 {
		updateSessionTime(sessionID)
	}

	return v, nil
}

// CreateSession create session for user, and register to global session manager
func CreateSession(phone string) string {
	sessionID := genSessionID(phone)
	session := userSession{
		SessionID: sessionID,
		Phone:     phone,
		LastTime:  time.Now(),
	}
	globalSession.Sessions[sessionID] = &session
	return sessionID
}

func updateSessionTime(sessionID string) {
	globalSession.Sessions[sessionID].LastTime = time.Now()
}

// SessionGC clear expired sessions
func (p *sessionManager) sessionGC() {
	ticker := time.NewTicker(300 * time.Second)
	go func() {
		for range ticker.C {
			for sid, _ := range p.Sessions {
				if sessionExpired(sid) {
					delete(p.Sessions, sid)
				}

			}
		}
	}()
}

func init() {
	globalSession = &sessionManager{"globalSession", map[string]*userSession{}}
}

// sessionExpired check session whether expired, return bool
func sessionExpired(sessionID string) bool {
	session, ok := globalSession.Sessions[sessionID]
	if !ok {
		return true
	}

	duration := time.Now().Sub(session.LastTime)
	if duration > MAX_LIFETIME_SECONDS {
		fmt.Println("bbbb")
		// no need gc here
		return true
	}
	return false
}

// genSessionID simple random algorithm
func genSessionID(phone string) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return fmt.Sprintf("%s-%d", phone, r1.Int31n(10000))
}
