package events

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvelopeUnmarshalJSON(t *testing.T) {
	t.Run("malformed top-level JSON", func(t *testing.T) {
		assert.Error(t, (&Envelope{}).UnmarshalJSON([]byte(`{"key":"value,`)))
	})
	t.Run("unrecognized event type", func(t *testing.T) {
		err := (&Envelope{}).UnmarshalJSON([]byte(`{"type":"unknown"}`))
		assert.True(t, errors.Is(err, ErrUnrecognizedEventType))
	})
	for _type, _struct := range event_type_instance_map {
		t.Run(string(_type), func(t *testing.T) {
			t.Run("successful", func(t *testing.T) {
				envelope := Envelope{}
				err := (&envelope).UnmarshalJSON([]byte(fmt.Sprintf(`{"type":"%v","event":{}}`, _type)))
				assert.Equal(t, reflect.New(reflect.TypeOf(_struct)).Interface(), envelope.Event)
				assert.Nil(t, err)
			})
			t.Run("null event payload", func(t *testing.T) {
				err := (&Envelope{}).UnmarshalJSON([]byte(fmt.Sprintf(`{"type":"%v","event":null}`, _type)))
				assert.True(t, errors.Is(err, ErrEnvelopeHasNoEventPayload))
			})
			t.Run("unsuccessful parsing", func(t *testing.T) {
				envelope := Envelope{}
				err := (&envelope).UnmarshalJSON([]byte(fmt.Sprintf(`{"type":"%v","event":{ broken-json }}`, _type)))
				assert.Nil(t, envelope.Event)
				assert.Error(t, err)
			})
		})
	}
}
