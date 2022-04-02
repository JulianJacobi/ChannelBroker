package broker

import (
    "testing"
    "time"
)

func TestChannelBroker(t *testing.T) {
    expectChan := func (c <-chan string, v, n string) {
        select {
        case s, ok := <-c:
            if s != v || (ok && v == "") {
                t.Fatalf("Read unexpected value from channel %s: Expect %s, got %s", n, v, s)
            }
        case <-time.After(5*time.Second):
            if v != "" {
                t.Fatalf("Channel %s is empty", n)
            }
        }
    }

    broker := New[string]()
    c1 := broker.NewChannel()
    c2 := broker.NewChannel()
    c3 := broker.NewChannel()

    broker.Chan <- "Test"

    expectChan(c1.Chan, "Test", "c1")
    expectChan(c2.Chan, "Test", "c2")
    expectChan(c3.Chan, "Test", "c3")

    broker.RemoveChannel(c2)

    broker.Chan <- "Test#2"

    expectChan(c1.Chan, "Test#2", "c1")
    expectChan(c2.Chan, "", "c2")
    expectChan(c3.Chan, "Test#2", "c3")

    broker.Clear()

    if len(broker.chans) != 0 {
        t.Fatalf("Broker channels not 0 after Clear: %d", len(broker.chans))
    }
}
