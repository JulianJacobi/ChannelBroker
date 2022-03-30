// ChannelBroker
// Write arbitrary objects to multiple channels, recieving channels
// can be added and removed dynamicaly.

package broker

import (
    "sync"
    "github.com/google/uuid"
)

// ChannelBroker represents the broker.
type ChannelBroker[T any] struct {
    Chan       chan<- T
    chans      map[uuid.UUID]chan T
    tChan      chan T
    addChan    chan *Channel[T]
    removeChan chan *Channel[T]
    lock       sync.Mutex
}

// Channel represents a receiving channel.
type Channel[T any] struct {
    Chan    <-chan T
    channel chan T
    uuid    uuid.UUID
}

// New channel broker instantiates a new ChannelBroker.
func New[T any]() *ChannelBroker[T] {
    tChan := make(chan T, 10)
    cb := &ChannelBroker[T]{
        Chan:       tChan,
        chans:      make(map[uuid.UUID]chan T),
        tChan:      tChan,
        addChan:    make(chan *Channel[T], 10),
        removeChan: make(chan *Channel[T], 10),
    }
    go cb.loop()
    return cb
}

// loop is an infinite loop to read from Broker Chan and distribute to receiving Channels
func (cb *ChannelBroker[T]) loop() {
    for {
        t := <-cb.tChan
        cb.distribute(t)
    }
}

// addChannel adds a channel to broker's channel map
func (cb *ChannelBroker[T]) addChannel(c *Channel[T]) {
    cb.lock.Lock()
    defer cb.lock.Unlock()
    cb.chans[c.uuid] = c.channel
}

// removeChannel from broker's channel map
func (cb *ChannelBroker[T]) removeChannel(c *Channel[T]) {
    cb.lock.Lock()
    defer cb.lock.Unlock()
    _, exists := cb.chans[c.uuid]
    if exists {
        delete(cb.chans, c.uuid)
    }
    close(c.channel)
}

// distribute object to all channels in broker's channel map
func (cb *ChannelBroker[T]) distribute(t T) {
    cb.lock.Lock()
    defer cb.lock.Unlock()
    for _, c := range cb.chans {
        c <- t
    }
}

// NewChannel creates new receiving Channel and adds it to the broker
func (cb *ChannelBroker[T]) NewChannel() *Channel[T] {
    channel := make(chan T, 10)
    c := &Channel[T]{
        Chan: channel,
        channel: channel,
        uuid: uuid.New(),
    }
    cb.addChannel(c)
    return c
}

// RemoveChannel from broker and closes the containing channel.
//
// Channel object should not be further used after removing it from broker,
// because the unlaying channel is closed after this call.
func (cb *ChannelBroker[T]) RemoveChannel(c *Channel[T]) {
    cb.removeChannel(c)
}
