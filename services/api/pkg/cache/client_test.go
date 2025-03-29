package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTest initializes a miniredis instance, Redis client, and context with a test logger.
func setupTest(t *testing.T) (*miniredis.Miniredis, cache.CacheClient, context.Context) {
	mr, err := miniredis.Run()
	require.NoError(t, err, "Failed to start miniredis")

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	//logger := zaptest.NewLogger(t)
	// ctx := apicontext.WithLogger(context.Background(), logger)

	return mr, cache.NewRedisCache(client), context.Background()
}

func TestFormatKey(t *testing.T) {
	rc := &cache.RedisCache{} // Client not needed for this test

	tests := []struct {
		name      string
		prefix    string
		id        interface{}
		wantKey   string
		wantError bool
	}{
		{
			name:    "nil id",
			prefix:  "user",
			id:      nil,
			wantKey: "user",
		},
		{
			name:    "string id",
			prefix:  "user",
			id:      "123",
			wantKey: "user:123",
		},
		{
			name:    "int id",
			prefix:  "post",
			id:      456,
			wantKey: "post:456",
		},
		{
			name:    "struct id",
			prefix:  "data",
			id:      struct{ Name string }{Name: "test"},
			wantKey: `data:{"Name":"test"}`,
		},
		{
			name:      "unmarshalable id",
			prefix:    "error",
			id:        make(chan int),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := rc.FormatKey(tt.prefix, tt.id)
			if tt.wantError {
				assert.Error(t, err, "Expected an error")
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tt.wantKey, key, "Formatted key mismatch")
		})
	}
}

func TestSetAndGet(t *testing.T) {
	mr, rc, ctx := setupTest(t)
	defer mr.Close()

	// Create a test organization ID
	orgId := uuid.MustParse("f4149aae-7c15-450c-a5a9-da358955a22a")
	// Add org ID to context
	ctx = apicontext.AddAuthToContext(ctx, "", uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), []uuid.UUID{orgId})

	type testStruct struct {
		Name string
		Age  int
	}

	t.Run("set and get struct", func(t *testing.T) {
		key := "testkey"
		expectedKey := fmt.Sprintf("%s:%s", orgId.String(), key)
		value := testStruct{Name: "Alice", Age: 30}

		err := rc.Set(ctx, key, value, 0)
		require.NoError(t, err, "Set failed")

		// Verify the key in redis has org prefix
		exists := mr.Exists(expectedKey)
		assert.True(t, exists, "Key with org prefix should exist")

		var result testStruct
		err = rc.Get(ctx, key, &result)
		require.NoError(t, err, "Get failed")
		assert.Equal(t, value, result, "Retrieved value mismatch")
	})

	t.Run("get non-existent key", func(t *testing.T) {
		key := "nonexistent"
		expectedKey := fmt.Sprintf("%s:%s", orgId.String(), key)
		var result testStruct
		err := rc.Get(ctx, key, &result)
		require.Error(t, err, "Expected error for non-existent key")
		assert.Contains(t, err.Error(), fmt.Sprintf("key not found: %s", expectedKey), "Error message mismatch")
	})

	t.Run("get invalid JSON", func(t *testing.T) {
		key := "invalidjson"
		expectedKey := fmt.Sprintf("%s:%s", orgId.String(), key)
		mr.Set(expectedKey, "invalid json")

		var result testStruct
		err := rc.Get(ctx, key, &result)
		require.Error(t, err, "Expected unmarshal error")

		var jsonSyntaxErr *json.SyntaxError
		var jsonTypeErr *json.UnmarshalTypeError

		isJSONErr := errors.As(err, &jsonSyntaxErr) || errors.As(err, &jsonTypeErr)
		assert.True(t, isJSONErr, "Error should wrap a JSON parsing error")
	})

	t.Run("set with expiry", func(t *testing.T) {
		key := "expiringkey"
		expectedKey := fmt.Sprintf("%s:%s", orgId.String(), key)
		value := "test value"
		expiry := time.Second

		err := rc.Set(ctx, key, value, expiry)
		require.NoError(t, err, "Set with expiry failed")

		exists, err := rc.Exists(ctx, key)
		require.NoError(t, err, "Exists check failed")
		assert.True(t, exists, "Key should exist before expiry")

		// Verify the actual key in redis has org prefix
		exists = mr.Exists(expectedKey)
		assert.True(t, exists, "Key with org prefix should exist")

		mr.FastForward(expiry)
		exists, err = rc.Exists(ctx, key)
		require.NoError(t, err, "Exists check after expiry failed")
		assert.False(t, exists, "Key should expire after TTL")
	})
}

func TestDelete(t *testing.T) {
	mr, rc, ctx := setupTest(t)
	defer mr.Close()

	key := "deletekey"
	require.NoError(t, rc.Set(ctx, key, "value", 0), "Setup: Set failed")

	t.Run("delete existing key", func(t *testing.T) {
		err := rc.Delete(ctx, key)
		require.NoError(t, err, "Delete failed")

		exists, err := rc.Exists(ctx, key)
		require.NoError(t, err, "Exists check failed")
		assert.False(t, exists, "Key should be deleted")
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		err := rc.Delete(ctx, "nonexistent")
		require.NoError(t, err, "Deleting non-existent key should not error")
	})
}

func TestExists(t *testing.T) {
	mr, rc, ctx := setupTest(t)
	defer mr.Close()

	key := "existskey"

	exists, err := rc.Exists(ctx, key)
	require.NoError(t, err, "Initial exists check failed")
	assert.False(t, exists, "Key should not exist initially")

	require.NoError(t, rc.Set(ctx, key, "value", 0), "Setup: Set failed")

	exists, err = rc.Exists(ctx, key)
	require.NoError(t, err, "Exists check failed")
	assert.True(t, exists, "Key should exist after set")
}

func TestIncrementDecrement(t *testing.T) {
	mr, rc, ctx := setupTest(t)
	defer mr.Close()

	key := "counter"

	t.Run("increment new key", func(t *testing.T) {
		val, err := rc.Increment(ctx, key)
		require.NoError(t, err, "Increment failed")
		assert.Equal(t, int64(1), val, "Incremented value mismatch")
	})

	t.Run("increment again", func(t *testing.T) {
		val, err := rc.Increment(ctx, key)
		require.NoError(t, err, "Increment failed")
		assert.Equal(t, int64(2), val, "Incremented value mismatch")
	})

	t.Run("decrement", func(t *testing.T) {
		val, err := rc.Decrement(ctx, key)
		require.NoError(t, err, "Decrement failed")
		assert.Equal(t, int64(1), val, "Decremented value mismatch")
	})

	t.Run("decrement new key", func(t *testing.T) {
		key := "counter2"
		val, err := rc.Decrement(ctx, key)
		require.NoError(t, err, "Decrement failed")
		assert.Equal(t, int64(-1), val, "Decremented value mismatch")
	})
}

func TestClient(t *testing.T) {
	_, rc, _ := setupTest(t)

	client := rc.Client()
	require.NotNil(t, client, "Client() returned nil")

	// Verify client connectivity
	_, err := client.Ping(context.Background()).Result()
	assert.NoError(t, err, "Client Ping failed")
}

func TestSetWithInvalidValue(t *testing.T) {
	_, rc, ctx := setupTest(t)

	key := "invalidvalue"
	invalidValue := func() {} // Unmarshallable value

	err := rc.Set(ctx, key, invalidValue, 0)
	require.Error(t, err, "Expected marshaling error")

}

func TestGetUnmarshalTypeError(t *testing.T) {
	mr, rc, ctx := setupTest(t)
	defer mr.Close()

	key := "typeerror"
	mr.Set(key, `"hello"`) // JSON string stored

	var num int
	err := rc.Get(ctx, key, &num)
	require.Error(t, err, "Expected unmarshal error")
}
