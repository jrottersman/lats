package helpers

import (
	"bytes"
	"strings"
	"testing"

	"github.com/manifoldco/promptui"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb ClosingBuffer) Close() error {
	return nil
}

func mockValidate(input string) error {
	return nil
}

func TestGeneratePrompt(t *testing.T) {
	pc := PromptContent{
		"foo",
		"bar",
	}
	resp := GeneratePrompt(pc, mockValidate)
	expected := "bar"
	if resp.Label != expected {
		t.Errorf("test failed expected %s, got %s", expected, resp.Label)
	}

	input := "baz"
	if resp.Validate(input) != nil {
		t.Errorf("validate function failed")
	}
}

func TestPromptInput(t *testing.T) {
	reader := ClosingBuffer{
		bytes.NewBufferString("Y\n"),
	}

	p := promptui.Prompt{
		Stdin: reader,
	}

	resp := PromptInput(p)
	expected := "Y"
	if !strings.EqualFold(resp, expected) {
		t.Errorf("expected %s, actual %s", expected, resp)
	}

}

func TestValidate(t *testing.T) {
	e := validate("")
	if e == nil {
		t.Errorf("expected validate to fail")
	}

	p := validate("foo")
	if p != nil {
		t.Errorf("this should be nil")
	}
}
