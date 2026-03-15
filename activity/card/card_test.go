// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package card

import (
	"encoding/json"
	"testing"
)

func TestContentTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"HeroCard", ContentTypeHeroCard, "application/vnd.microsoft.card.hero"},
		{"AdaptiveCard", ContentTypeAdaptiveCard, "application/vnd.microsoft.card.adaptive"},
		{"OAuthCard", ContentTypeOAuthCard, "application/vnd.microsoft.card.oauth"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.value != tc.expected {
				t.Errorf("got %q, want %q", tc.value, tc.expected)
			}
		})
	}
}

func TestHeroCardJSONRoundTrip(t *testing.T) {
	card := &HeroCard{
		Title:    "Test Title",
		Subtitle: "Test Subtitle",
		Text:     "Test text",
		Images: []*CardImage{
			{URL: "https://example.com/img.png", Alt: "an image"},
		},
		Buttons: []*CardAction{
			{Type: "openUrl", Title: "Click me", Value: "https://example.com"},
		},
	}

	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got HeroCard
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.Title != card.Title {
		t.Errorf("Title: got %q, want %q", got.Title, card.Title)
	}
	if got.Subtitle != card.Subtitle {
		t.Errorf("Subtitle: got %q, want %q", got.Subtitle, card.Subtitle)
	}
	if got.Text != card.Text {
		t.Errorf("Text: got %q, want %q", got.Text, card.Text)
	}
	if len(got.Images) != 1 || got.Images[0].URL != "https://example.com/img.png" {
		t.Errorf("Images not round-tripped correctly: %+v", got.Images)
	}
	if len(got.Buttons) != 1 || got.Buttons[0].Title != "Click me" {
		t.Errorf("Buttons not round-tripped correctly: %+v", got.Buttons)
	}
}

func TestHeroCardOmitsEmptyFields(t *testing.T) {
	card := &HeroCard{Title: "Only Title"}
	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if _, ok := m["subtitle"]; ok {
		t.Error("expected subtitle to be omitted but it was present")
	}
	if _, ok := m["images"]; ok {
		t.Error("expected images to be omitted but it was present")
	}
}

func TestAdaptiveCardInvokeActionRoundTrip(t *testing.T) {
	action := &AdaptiveCardInvokeAction{
		Type: "Action.Execute",
		ID:   "action-1",
		Verb: "doSomething",
		Data: map[string]interface{}{"key": "value"},
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got AdaptiveCardInvokeAction
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.Type != action.Type {
		t.Errorf("Type: got %q, want %q", got.Type, action.Type)
	}
	if got.ID != action.ID {
		t.Errorf("ID: got %q, want %q", got.ID, action.ID)
	}
	if got.Verb != action.Verb {
		t.Errorf("Verb: got %q, want %q", got.Verb, action.Verb)
	}
}

func TestAdaptiveCardInvokeValueRoundTrip(t *testing.T) {
	val := &AdaptiveCardInvokeValue{
		Action: &AdaptiveCardInvokeAction{
			Type: "Action.Execute",
			Verb: "refresh",
		},
		Authentication: &TokenExchangeInvokeRequest{
			ID:             "req-1",
			ConnectionName: "myConn",
			Token:          "tok123",
		},
		State: "magic-code",
	}

	data, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got AdaptiveCardInvokeValue
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.State != "magic-code" {
		t.Errorf("State: got %q, want %q", got.State, "magic-code")
	}
	if got.Action == nil || got.Action.Verb != "refresh" {
		t.Errorf("Action.Verb not round-tripped: %+v", got.Action)
	}
	if got.Authentication == nil || got.Authentication.Token != "tok123" {
		t.Errorf("Authentication.Token not round-tripped: %+v", got.Authentication)
	}
}

func TestAdaptiveCardInvokeResponseRoundTrip(t *testing.T) {
	resp := &AdaptiveCardInvokeResponse{
		StatusCode: 200,
		Type:       "application/vnd.microsoft.card.adaptive",
		Value:      map[string]interface{}{"result": "ok"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got AdaptiveCardInvokeResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.StatusCode != resp.StatusCode {
		t.Errorf("StatusCode: got %d, want %d", got.StatusCode, resp.StatusCode)
	}
	if got.Type != resp.Type {
		t.Errorf("Type: got %q, want %q", got.Type, resp.Type)
	}
}

func TestOAuthCardRoundTrip(t *testing.T) {
	card := &OAuthCard{
		Text:           "Sign in",
		ConnectionName: "myConnection",
		Buttons: []*CardAction{
			{Type: "signin", Title: "Sign in", Value: "https://signin.example.com"},
		},
		TokenExchangeResource: &TokenExchangeResource{
			ID:                  "exchange-1",
			URI:                 "https://token.example.com",
			ProviderDisplayName: "My Provider",
		},
	}

	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got OAuthCard
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.Text != card.Text {
		t.Errorf("Text: got %q, want %q", got.Text, card.Text)
	}
	if got.ConnectionName != card.ConnectionName {
		t.Errorf("ConnectionName: got %q, want %q", got.ConnectionName, card.ConnectionName)
	}
	if got.TokenExchangeResource == nil {
		t.Fatal("expected TokenExchangeResource to be non-nil")
	}
	if got.TokenExchangeResource.ProviderDisplayName != "My Provider" {
		t.Errorf("ProviderDisplayName: got %q, want %q", got.TokenExchangeResource.ProviderDisplayName, "My Provider")
	}
}

func TestCardActionJSONFieldNames(t *testing.T) {
	action := &CardAction{
		Type:        "openUrl",
		Title:       "Go",
		Image:       "https://img.example.com",
		Text:        "click me",
		DisplayText: "shown in feed",
		ImageAltText: "alt text",
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	for _, key := range []string{"type", "title", "image", "text", "displayText", "imageAltText"} {
		if _, ok := m[key]; !ok {
			t.Errorf("expected JSON field %q to be present", key)
		}
	}
}

func TestCardImageJSONRoundTrip(t *testing.T) {
	img := &CardImage{
		URL: "https://example.com/img.png",
		Alt: "description",
	}
	data, err := json.Marshal(img)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var got CardImage
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got.URL != img.URL {
		t.Errorf("URL: got %q, want %q", got.URL, img.URL)
	}
	if got.Alt != img.Alt {
		t.Errorf("Alt: got %q, want %q", got.Alt, img.Alt)
	}
}
