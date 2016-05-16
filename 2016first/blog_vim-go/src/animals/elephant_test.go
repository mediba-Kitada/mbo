package animals

import "testing"

func TestElephantFeed(t *testing.T) {
	expect := "Grass"
	autual := ElephantFeed()

	if expect != autual {
		t.Errorf("%s != %s", expect, autual)
	}
}
