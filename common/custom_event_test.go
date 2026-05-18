package common

import (
	"net/http/httptest"
	"testing"
)

func TestCustomEventRenderWritesSSEHeadersAndData(t *testing.T) {
	recorder := httptest.NewRecorder()
	event := &CustomEvent{Data: "data: aud020"}

	if err := event.Render(recorder); err != nil {
		t.Fatalf("render custom event: %v", err)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("content type = %q, want text/event-stream", got)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("cache-control = %q, want no-cache", got)
	}
	if got := recorder.Body.String(); got != "data: aud020\n\n" {
		t.Fatalf("body = %q", got)
	}
}
