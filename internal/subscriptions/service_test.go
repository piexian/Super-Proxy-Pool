package subscriptions

import (
	"testing"
	"time"

	"super-proxy-pool/internal/models"
)

func TestShouldSyncSubscription(t *testing.T) {
	now := time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC)
	recent := now.Add(-5 * time.Minute)
	old := now.Add(-2 * time.Hour)

	cases := []struct {
		name string
		item models.Subscription
		want bool
	}{
		{
			name: "disabled subscription",
			item: models.Subscription{Enabled: false, SyncIntervalSec: 3600},
			want: false,
		},
		{
			name: "never synced",
			item: models.Subscription{Enabled: true, SyncIntervalSec: 3600},
			want: true,
		},
		{
			name: "not due yet",
			item: models.Subscription{Enabled: true, SyncIntervalSec: 3600, LastSyncAt: &recent},
			want: false,
		},
		{
			name: "due now",
			item: models.Subscription{Enabled: true, SyncIntervalSec: 3600, LastSyncAt: &old},
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := shouldSyncSubscription(tc.item, now); got != tc.want {
				t.Fatalf("shouldSyncSubscription() = %v, want %v", got, tc.want)
			}
		})
	}
}
