package sparkpost

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		cfg             Config
		wantErr         bool
		wantErrType     error
		wantErrContains string
	}{
		{
			name: "valid config",
			cfg: Config{
				APIKey: "test-key",
				APIUrl: "https://api.sparkpost.com",
			},
			wantErr: false,
		},
		{
			name:            "empty config",
			cfg:             Config{},
			wantErr:         true,
			wantErrType:     ErrInvalidConfig,
			wantErrContains: "invalid sparkpost configuration",
		},
		{
			name: "empty API key",
			cfg: Config{
				APIUrl: "https://api.sparkpost.com",
			},
			wantErr:         true,
			wantErrType:     ErrInvalidConfig,
			wantErrContains: "invalid sparkpost configuration",
		},
		{
			name: "empty API URL",
			cfg: Config{
				APIKey: "test-key",
			},
			wantErr:         true,
			wantErrType:     ErrInvalidConfig,
			wantErrContains: "invalid sparkpost configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, client)
				if tt.wantErrType != nil {
					require.ErrorIs(t, err, tt.wantErrType)
				}
				if tt.wantErrContains != "" {
					require.ErrorContains(t, err, tt.wantErrContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}
