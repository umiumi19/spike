package main

import (
	"context"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetGreeting(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	ctx := context.Background() // define your necessary context

	resp := api.GetCtx(ctx, "/greeting/world") // provide it using the 'Ctx' suffixed methods
	if !strings.Contains(resp.Body.String(), "Hello, world!") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestPutReview(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	resp := api.Post("/reviews", map[string]any{
		"author": "daniel",
		"rating": 5,
	})

	if resp.Code != 201 {
		t.Fatalf("Unexpected status code: %d", resp.Code)
	}
}

func TestPutReviewError(t *testing.T) {
	_, api := humatest.New(t)

	addRoutes(api)

	resp := api.Post("/reviews", map[string]any{
		"rating": 10,
	})

	if resp.Code != 422 {
		t.Fatalf("Unexpected status code: %d", resp.Code)
	}
}