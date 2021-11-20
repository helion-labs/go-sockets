package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Location struct {
	Distance int
	Time     int
}

var Distance_chan chan []byte

func NewTLSConfig() (config *tls.Config, err error) {
	// Import trusted certificates from CAfile.pem.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("./certs/rootCA.pem")
	if err != nil {
		return
	}
	certpool.AppendCertsFromPEM(pemCerts)

	// Import client certificate/key pair.
	cert, err := tls.LoadX509KeyPair("./certs/cert.crt", "./certs/private.key")
	if err != nil {
		return
	}

	// Create tls.Config with desired tls properties
	config = &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
	return
}

func main_mqtt() {
	Distance_chan = make(chan []byte)

	tlsconfig, err := NewTLSConfig()
	if err != nil {
		log.Fatalf("failed to create TLS configuration: %v", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tls://a3dp44gpbbyb4r-ats.iot.us-west-2.amazonaws.com:8883")
	opts.SetClientID("clientID-home").SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)

	// Start the connection.
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}

	fmt.Println("Listening for new events.")
	if token := c.Subscribe("/topic/cat_location", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}

	for {
	}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	//Not 100% sure if it's safe to keep a reference to msg.Payload
	b := make([]byte, len(msg.Payload()))
	copy(b, msg.Payload())

	// Sanitize message
	l := Location{}
	err := json.Unmarshal(msg.Payload(), &l)
	if err != nil {
		fmt.Println("Strange... could not unmarshal json! error:", err)
		return
	}

	// Error if more than 10 meters away
	if l.Distance > 1000 {
		fmt.Println("Strange distance!")
		return
	}

	// send up the JSON byte representation
	// JSON on the browser will unmarshal
	// Non-blocking
	select {
	case Distance_chan <- b:
	default:
		fmt.Println("Channel blocked!")
	}
}
