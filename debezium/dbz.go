package debezium

import (
	"encoding/json"
)

type DebeziumMsg struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Before json.RawMessage `json:"before"`
	After  json.RawMessage `json:"after"`
	Op     string          `json:"op"`
}
