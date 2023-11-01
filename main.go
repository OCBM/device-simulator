package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// AWS IoT MQTT endpoint
	awsIotEndpoint := "a18chofoi4appa-ats.iot.ap-southeast-1.amazonaws.com"
	topic := "con/topic"

	// Load client certificate and key
	certFile := "./certs/clientCert.crt"
	keyFile := "./certs/clientKey.key"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Error loading client certificate and key: %v", err)
	}

	// Load server certificate
	caCertFile := "./certs/serverCert.pem"
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error loading server certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // Set this to 'false' in production.
	}

	// Create an MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:8883", awsIotEndpoint))
	opts.SetClientID("test-con")
	opts.SetTLSConfig(tlsConfig)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to AWS IoT: %v", token.Error())
		os.Exit(1)
	}
	defer client.Disconnect(250)

	connection_message := "connected sucessfully"

	// Encode the JSON object as a string
	payload, err := json.Marshal(connection_message)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	// Publish a test message
	token := client.Publish(topic, 0, false, string(payload))
	token.Wait()
	fmt.Printf("Published message on topic: %s", topic)
}
