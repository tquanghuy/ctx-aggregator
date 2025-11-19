package aggregator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildContextKey(t *testing.T) {
	ctxKey := buildContextKey("key1", "key2")
	assert.Equal(t, ctxKey, contextKey("ctxAggCtxKey_key1_key2"))
}
