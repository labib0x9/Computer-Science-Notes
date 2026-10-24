package helloworld

import "testing"

func TestHelloWorld(t *testing.T) {
	name := "Labib"
	expected := "Hello " + name
	got := HelloWorld(name)

	if expected != got {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}
