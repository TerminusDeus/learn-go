package main

import (
	"fmt"
	"v/csrf"

	"github.com/go-stomp/stomp"
)

type (
	Consumer struct {
		callBack map[string]chan []byte
		qReply   string
		mq       *stomp.Conn
		sub      *stomp.Subscription
	}
)

func isErr(e error) {
	if e != nil {
		panic(e)
	}
}

func NewConsumer(host, qReply string) *Consumer {
	var e error

	cn := new(Consumer)
	cn.callBack = make(map[string]chan []byte)
	cn.qReply = qReply

	cn.mq, e = stomp.Dial("tcp", fmt.Sprintf("%s:61613", host))
	isErr(e)
	cn.sub, e = cn.mq.Subscribe(qReply, stomp.AckAuto)
	isErr(e)

	go func() {
		for {
			msg, e := cn.sub.Read()
			isErr(e)

			corrId := msg.Header.Get("correlation-id")
			if ch, ok := cn.callBack[corrId]; ok {
				ch <- msg.Body
			}
		}
	}()

	return cn
}

func (cn *Consumer) Request(qRequest, operation string, in []byte) []byte {
	// Correlation ID
	corrId, e := csrf.GenerateRandomString(32)
	isErr(e)

	// Register callback
	ch := make(chan []byte)
	cn.callBack[corrId] = ch

	// Request
	isErr(cn.mq.Send(
		qRequest,           // destination
		"application/json", // content-type
		in,                 // body
		stomp.SendOpt.Header("reply-to", cn.qReply),
		stomp.SendOpt.Header("correlation-id", corrId),
		stomp.SendOpt.Header("operation-name", operation),
	))

	return <-ch
}

func (cn *Consumer) Close() {
	isErr(cn.mq.Disconnect())
}