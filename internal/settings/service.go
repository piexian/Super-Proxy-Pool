package settings

import (
	"context"
	"database/sql"
	"time"

	"super-proxy-pool/internal/config"
	"super-proxy-pool/internal/db"
	"super-proxy-pool/internal/models"
)

type Service struct {
	store *db.Store
	cfg   config.App
}

func NewService(store *db.Store, cfg config.App) *Service {
	return &Service{store: store, cfg: cfg}
}

func (s *Service) EnsureDefaults(ctx context.Context, passwordHash string) error {
	var count int
	if err := s.store.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM settings WHERE id = 1`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	now := time.Now().UTC()
	_, err := s.store.DB.ExecContext(ctx, `INSERT INTO settings (
		id, panel_host, panel_port, password_hash, speed_test_enabled, latency_test_url, speed_test_url,
		latency_timeout_ms, speed_timeout_ms, latency_concurrency, speed_concurrency,
		default_subscription_interval_sec, mihomo_controller_secret, failure_retry_count, log_level,
		speed_max_bytes, created_at, updated_at
	) VALUES (1, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?, 2, 'info', ?, ?, ?)`,
		s.cfg.PanelHost,
		s.cfg.PanelPort,
		passwordHash,
		config.DefaultLatencyURL(),
		config.DefaultSpeedURL(),
		config.DefaultLatencyTimeoutMS(),
		config.DefaultSpeedTimeoutMS(),
		config.DefaultLatencyConcurrency(),
		config.DefaultSpeedConcurrency(),
		config.DefaultSubscriptionIntervalSec(),
		s.cfg.DefaultControllerSecret,
		config.DefaultSpeedMaxBytes(),
		now,
		now,
	)
	return err
}

func (s *Service) Get(ctx context.Context) (models.Settings, error) {
	row := s.store.DB.QueryRowContext(ctx, `SELECT id, panel_host, panel_port, password_hash, speed_test_enabled,
		latency_test_url, speed_test_url, latency_timeout_ms, speed_timeout_ms, latency_concurrency,
		speed_concurrency, default_subscription_interval_sec, mihomo_controller_secret, failure_retry_count,
		log_level, speed_max_bytes, created_at, updated_at FROM settings WHERE id = 1`)
	return scanSettings(row)
}

func (s *Service) Update(ctx context.Context, current models.Settings) (models.Settings, bool, error) {
	existing, err := s.Get(ctx)
	if err != nil {
		return models.Settings{}, false, err
	}
	restartRequired := existing.PanelHost != current.PanelHost || existing.PanelPort != current.PanelPort
	current.ID = 1
	current.PasswordHash = existing.PasswordHash
	current.UpdatedAt = time.Now().UTC()
	_, err = s.store.DB.ExecContext(ctx, `UPDATE settings SET
		panel_host = ?, panel_port = ?, speed_test_enabled = ?, latency_test_url = ?, speed_test_url = ?,
		latency_timeout_ms = ?, speed_timeout_ms = ?, latency_concurrency = ?, speed_concurrency = ?,
		default_subscription_interval_sec = ?, mihomo_controller_secret = ?, failure_retry_count = ?,
		log_level = ?, speed_max_bytes = ?, updated_at = ? WHERE id = 1`,
		current.PanelHost, current.PanelPort, boolToInt(current.SpeedTestEnabled), current.LatencyTestURL, current.SpeedTestURL,
		current.LatencyTimeoutMS, current.SpeedTimeoutMS, current.LatencyConcurrency, current.SpeedConcurrency,
		current.DefaultSubscriptionIntervalSec, current.MihomoControllerSecret, current.FailureRetryCount,
		current.LogLevel, current.SpeedMaxBytes, current.UpdatedAt,
	)
	if err != nil {
		return models.Settings{}, false, err
	}
	updated, err := s.Get(ctx)
	return updated, restartRequired, err
}

func (s *Service) UpdatePasswordHash(ctx context.Context, hash string) error {
	_, err := s.store.DB.ExecContext(ctx, `UPDATE settings SET password_hash = ?, updated_at = ? WHERE id = 1`, hash, time.Now().UTC())
	return err
}

func scanSettings(scanner interface{ Scan(dest ...any) error }) (models.Settings, error) {
	var item models.Settings
	var speedEnabled int
	err := scanner.Scan(
		&item.ID, &item.PanelHost, &item.PanelPort, &item.PasswordHash, &speedEnabled,
		&item.LatencyTestURL, &item.SpeedTestURL, &item.LatencyTimeoutMS, &item.SpeedTimeoutMS,
		&item.LatencyConcurrency, &item.SpeedConcurrency, &item.DefaultSubscriptionIntervalSec,
		&item.MihomoControllerSecret, &item.FailureRetryCount, &item.LogLevel, &item.SpeedMaxBytes,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return models.Settings{}, err
	}
	item.SpeedTestEnabled = speedEnabled == 1
	return item, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
