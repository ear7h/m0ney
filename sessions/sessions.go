package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var g_sessions map[string]*Session

func init() {
	g_sessions = map[string]*Session{}
	fmt.Println("sessions initiated")
}


func Exists(key string) bool{
	return g_sessions[key] != nil
}

// Create adds a new session and returns the key string associated with the session
func Create(sess *Session) string {
	token := make([]byte, 10)
	rand.Read(token)


	tokenStr := base64.StdEncoding.EncodeToString(token)

	//do not overwrite sessions
	for Exists(tokenStr) {
		rand.Read(token)
		tokenStr = base64.StdEncoding.EncodeToString(token)
	}

	g_sessions[tokenStr] = sess

	return tokenStr
}

func Get(key string) *Session {
	return g_sessions[key]
}