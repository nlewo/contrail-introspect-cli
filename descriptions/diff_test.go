package descriptions

import "testing"

func TestDiff(t *testing.T) {
	a := []string{"192.168.1.3/32", "192.168.10.4/32", "192.168.10.6/32", "192.168.10.7/32", "192.168.10.8/32", "192.168.10.9/32", "192.168.10.10/32", "192.168.10.11/32", "192.168.20.3/32"}
	b := []string{"169.254.169.254/32", "192.168.1.3/32", "192.168.10.0/24", "192.168.10.1/32", "192.168.10.2/32", "192.168.10.4/32", "192.168.10.6/32", "192.168.10.7/32", "192.168.10.8/32", "192.168.10.9/32", "192.168.10.10/32", "192.168.10.11/32", "192.168.20.3/32"}
	aOutput := []string{}
	bOutput := []string{"169.254.169.254/32", "192.168.10.0/24", "192.168.10.1/32", "192.168.10.2/32"}

	a, b = diff(a, b)

	if len(a) != len(aOutput) {
		t.Fatalf("a and aOutput don't have the same size len(a)=%d vs len(aOutput)=%d", len(a), len(aOutput))
	}
	if len(b) != len(bOutput) {
		t.Fatalf("b and bOutput don't have the same size len(b)=%d vs len(bOutput)=%d", len(b), len(bOutput))
	}
	for i, _ := range a {
		if a[i] != aOutput[i] {
			t.Fatalf("a and aOutput have not the same value at index %d : a[i]=%s vs aOutput[i]=%s", i, a[i], aOutput[i])
		}
	}
	for i, _ := range b {
		if b[i] != bOutput[i] {
			t.Fatalf("b and bOutput have not the same value at index %d : b[i]=%s vs bOutput[i]=%s", i, b[i], bOutput[i])
		}
	}
}
