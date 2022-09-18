package services

import (
	"fmt"
	"rabbitmq/models"
	"sync"
)

var wg sync.WaitGroup

func CreateProof() error {

	consumerP1 := models.ConsumerSettings{QueueName: "Proof1", Name: "Create_proof1", AutoAck: false, Exclusive: false, NoLocal: false, NoWait: false}
	consumerP2 := models.ConsumerSettings{QueueName: "Proof2", Name: "Create_proof2", AutoAck: false, Exclusive: false, NoLocal: false, NoWait: false}

	messagesP1, err := CreateConsumer(consumerP1, false, true)

	for ms := range messagesP1 {
		fmt.Println(ms)
	}

	fmt.Println("t1")

	if err != nil {
		return err
	}

	messagesP2, err := CreateConsumer(consumerP2, false, true)

	for ms := range messagesP2 {
		fmt.Println(ms)
	}

	fmt.Println("t1")
	if err != nil {
		return err
	}

	return nil
}
