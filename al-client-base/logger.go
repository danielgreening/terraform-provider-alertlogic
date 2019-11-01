package al_client_base

import (
	"log"
	"strings"
)

type DebugLogger struct{}

func (l DebugLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.Printf("[DEBUG] [al_client] %s", strings.Join(tokens, " "))
}
