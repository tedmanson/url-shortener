package server

import (
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetTTL(t *testing.T) {
	type test struct {
		expire time.Time
		want   string
	}

	now := time.Now()

	tests := []test{
		{expire: time.Now().Add(4 * time.Hour), want: "max-age=3600"},
		{expire: time.Now().Add(100 * time.Second), want: "max-age=100"},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		addTTLHeader(w, now, tc.expire)
		got := w.Header().Get("Cache-Control")

		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
