![tests](https://github.com/JulianJacobi/ChannelBroker/actions/workflows/test.yml/badge.svg)
[![Go version](https://img.shields.io/github/go-mod/go-version/JulianJacobi/ChannelBroker.svg)](https://github.com/JulianJacobi/ChannelBroker)
[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/JulianJacobi/ChannelBroker)

# Channel Broker

A helper for golang that provides a broker based on channels.
So you send to one channel receive by multiple channels.
Receving channels can be created and removed at runtime.

## Usage Example

    package main
    
    import (
        "fmt"
        "github.com/JulianJacobi/ChannelBroker"
    )
    
    func main() {
        b := broker.New[string]()
    
        c1 := b.NewChannel()
        c2 := b.NewChannel()
    
        b.Chan <- "Test"
    
        r1 := <-c1.Chan
        fmt.Println(r1)
    
        r2 := <-c2.Chan
        fmt.Println(r2)
    }
