package ex02

import "testing"

func TestMultiplex(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	out := multiplex(ch1, ch2, ch3)

	expected := []string{"Message 1 on ch1", "Message 2 on ch2", "Message 3 on ch3"}
	go func() {
		ch1 <- "Message 1 on ch1"
		ch2 <- "Message 2 on ch2"
		ch3 <- "Message 3 on ch3"
		close(ch1)
		close(ch2)
		close(ch3)
	}()

	var received []string
	for msg := range out {
		received = append(received, msg.(string))
	}

	if len(expected) != len(received) {
		t.Errorf("Expected %d messages, but received %d", len(expected), len(received))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != received[i] {
			t.Errorf("Expected message %q, but received %q", expected[i], received[i])
		}
	}
}
