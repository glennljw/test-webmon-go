package record

import "testing"

func TestParseLine(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantStatus string
		wantErr    bool
	}{
		{"v1 valid", "2026-10-01T10:00:00Z|https://a.gov.sg|up", "up", false},
		{"v2 valid", `{"target":"https://a.gov.sg","result":"down","ts":"2026-10-02T10:00:00Z"}`, "down", false},
		{"invalid", "garbage", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.Status != tt.wantStatus {
				t.Errorf("status = %q, want %q", got.Status, tt.wantStatus)
			}
		})
	}

}
