================
ToyBroker
================

A MQTT ( http://mqtt.org ) broker written in Go.

- http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html


Restriction
-----------------------------

- Only support QoS=0.
- Not support Clean Session at CONNECT.
- Not support Will Retain, Will QoS, Will Flag ad CONNECT.
- Not support wildcard subscribe.


How to compile and execute
-----------------------------

::

   $ git clone git@github.com:nakagami/toybroker.git
   $ cd toybroker/toybroker
   $ go build
   $ ./toybroker


Customize
----------------------

You can wrote a customized broker.

At first::

   $ go get github.com/nakagami/toybroker

And simple example::

   package main

   import (
       "net"
       "github.com/nakagami/toybroker"
   )

   func main() {
       var topics toybroker.Topics = toybroker.NewMemoryTopics()
       var hook toybroker.Hook = toybroker.NewMemoryHook()

       listener, _ := net.Listen("tcp", ":1883")
       for {
           conn, _ := listener.Accept()
           go toybroker.MqttMainLoop(conn, topics, hook)
       }
   }
