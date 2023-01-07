package bee

import (
	"context"
	"github.com/luminosita/honeycomb/pkg/http"
	"github.com/luminosita/honeycomb/pkg/http/handlers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBeeServer(t *testing.T) {
	s := NewBeeServer(&Config{})

	routes := s.Routes(context.Background())

	assert.NotNil(t, routes)
	assert.Equal(t, 3, len(routes))

	for _, v := range routes {
		if v.Type != http.STATIC {
			_, ok := v.Handler.(handlers.Handler)
			assert.True(t, ok)
		}
	}
}
