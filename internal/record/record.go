package record

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CheckRecord struct {
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
}

type v2Log struct {
	Target string    `json:"target"`
	Result string    `json:"result"`
	TS     time.Time `json:"ts"`
}

var ErrBadFormat = errors.New("unrecognised format")

func ParseLine(line string) (CheckRecord, error) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "{") {
		var record v2Log
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return CheckRecord{}, fmt.Errorf("Parse v2: %w", err)
		}
		return CheckRecord{
			URL:       record.Target,
			Status:    record.Result,
			CheckedAt: record.TS,
		}, nil
	}
	parts := strings.Split(line, "|")
	if len(parts) != 3 {
		return CheckRecord{}, fmt.Errorf("Parse v1: %w", ErrBadFormat)
	}
	t, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return CheckRecord{}, fmt.Errorf("Parse v1 time: %w", err)
	}
	return CheckRecord{
		URL:       parts[1],
		Status:    parts[2],
		CheckedAt: t,
	}, nil
}
