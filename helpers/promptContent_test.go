package helpers

import "testing"

func TestGeneratePrompt(t *testing.T) {
	pc := PromptContent{
		"foo",
		"bar",
	}
	p := GeneratePrompt(pc)
	expected := "bar"
	if p.Label != expected {
		t.Errorf("test failed expected %s, got %s", expected, p.Label)
	}

	input := "baz"
	if p.Validate(input) != nil {
		t.Errorf("validate function failed")
	}
}
