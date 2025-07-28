package observability

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// generateTraceID generates a unique trace ID
func generateTraceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// generateSpanID generates a unique span ID
func generateSpanID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// getCommandType returns a string representation of the command type
func getCommandType(cmd tea.Cmd) string {
	if cmd == nil {
		return "nil"
	}
	
	// Use reflection to get the type name
	cmdType := reflect.TypeOf(cmd)
	if cmdType.Kind() == reflect.Func {
		return "func"
	}
	
	return cmdType.String()
}

// getMessageType returns a string representation of the message type
func getMessageType(msg tea.Msg) string {
	if msg == nil {
		return "nil"
	}
	
	msgType := reflect.TypeOf(msg)
	// Get the short type name without package
	typeName := msgType.String()
	if strings.Contains(typeName, ".") {
		parts := strings.Split(typeName, ".")
		typeName = parts[len(parts)-1]
	}
	
	return typeName
}

// copyMetadata creates a deep copy of metadata map
func copyMetadata(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// getAtomType returns a string representation of the atom type
func getAtomType(atom interface{}) string {
	if atom == nil {
		return "nil"
	}
	
	atomType := reflect.TypeOf(atom)
	// Handle pointers
	if atomType.Kind() == reflect.Ptr {
		atomType = atomType.Elem()
	}
	
	// Get the short type name without package
	typeName := atomType.String()
	if strings.Contains(typeName, ".") {
		parts := strings.Split(typeName, ".")
		typeName = parts[len(parts)-1]
	}
	
	return typeName
}