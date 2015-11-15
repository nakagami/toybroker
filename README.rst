================
ToyBroker
================

A MQTT broker.

Restriction
-----------------------------

Only support QoS=0.


How to compile and execute
-----------------------------

::

   $ git clone git@github.com:nakagami/toybroker.git
   $ cd toybroker/toybroker
   $ go build
   $ ./toybroker


Customize
----------------------

You can wrote your application.

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
