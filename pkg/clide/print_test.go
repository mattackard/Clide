package clide

import "testing"

func TestPrompt(t *testing.T) {
	cfg := Config{
		User:      "test",
		Directory: "/",
	}

	if prompt(cfg) != "test:/$ " {
		t.Errorf("Prompt returned unextected string. Expected test:/$ , got %s", prompt(cfg))
	}
}
