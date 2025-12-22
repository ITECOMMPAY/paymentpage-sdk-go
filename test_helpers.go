package paymentpage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func LoadJsonFromFile(t *testing.T, filename string) map[string]any {
	t.Helper()

	path := filepath.Join("testdata", filename)
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		t.Fatalf("failed to unmarshal JSON %s: %v", path, err)
	}

	return data
}
