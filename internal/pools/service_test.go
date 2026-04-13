package pools

import (
	"testing"

	"super-proxy-pool/internal/models"
)

func TestValidatePortConflict(t *testing.T) {
	pools := []models.ProxyPool{
		{ID: 1, Name: "a", ListenPort: 8080},
		{ID: 2, Name: "b", ListenPort: 8081},
	}

	if err := ValidatePortConflict(7890, pools, 0, 7890); err == nil {
		t.Fatalf("expected panel port conflict")
	}
	if err := ValidatePortConflict(7890, pools, 0, 8081); err == nil {
		t.Fatalf("expected pool port conflict")
	}
	if err := ValidatePortConflict(7890, pools, 2, 8081); err != nil {
		t.Fatalf("expected current pool port to be allowed, got %v", err)
	}
	if err := ValidatePortConflict(7890, pools, 0, 18080); err != nil {
		t.Fatalf("expected non-conflicting port, got %v", err)
	}
}
