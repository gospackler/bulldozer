package queue

import (
	"testing"
)

func TestQ(t *testing.T) {
	q := New()

	q.Add(10)
	q.Add(20)

	elem := q.Remove()
	if elem != nil {
		t.Log(elem)
	} else {
		t.Errorf("expected 10")
	}
}
