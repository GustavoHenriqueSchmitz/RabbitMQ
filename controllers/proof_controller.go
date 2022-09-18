package controllers

import (
	"rabbitmq/models"
	"rabbitmq/services"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var wg sync.WaitGroup

func Proofs(a *fiber.Ctx) error {

	proof1 := models.QueueSettings{Name: "Proof1", Durable: true, AutoDelete: false, Exclusive: false, NoWait: false}
	proof2 := models.QueueSettings{Name: "Proof2", Durable: true, AutoDelete: false, Exclusive: false, NoWait: false}
	exchangeConfig := models.ExchangeSettings{Name: "CreateProofs", Type: "direct", Durable: true, AutoDelete: false, Internal: false, NoWait: false}
	proof1Bind := models.BindingSettings{QueueName: "Proof1", RoutingKey: "p1", Exchange: "CreateProofs", NoWait: false}
	proof2Bind := models.BindingSettings{QueueName: "Proof2", RoutingKey: "p2", Exchange: "CreateProofs", NoWait: false}

	proof1Publisher := models.PublishSettings{Exchange: "CreateProofs", RoutingKey: "p1", Mandatory: false, Immediate: false, Publish: models.Message{ContentType: "text/plan", Body: models.Proof1{Name: "PostgresSQL", Questions1: "Q1", Question2: "Q2", Questions3: "Q3"}}}
	proof2Publisher := models.PublishSettings{Exchange: "CreateProofs", RoutingKey: "p2", Mandatory: false, Immediate: false, Publish: models.Message{ContentType: "text/plan", Body: models.Proof2{Name: "MongoDB", Questions1: "Q1", Question2: "Q2", Questions3: "Q3", Questions4: "Q4", Question5: "Q5", Questions6: "Q6"}}}

	amountUrl := a.Query("amount")
	proofUrl := a.Query("proof")

	amount, err := strconv.Atoi(amountUrl)
	if err != nil {
		return err
	}

	proof, err := strconv.Atoi(proofUrl)
	if err != nil {
		return err
	}

	services.CreateQueue(proof1)
	services.CreateQueue(proof2)
	services.CreateExchange(exchangeConfig)
	services.CreateBind(proof1Bind)
	services.CreateBind(proof2Bind)

	if proof == 1 {
		wg.Add(1)
		go func() {
			for cont := 0; cont <= amount; cont += 1 {

				err = services.PublishMessage(proof1Publisher)
				if err != nil {
					cont -= 1
					continue
				}
			}
		}()
		wg.Done()
	} else if proof == 2 {
		wg.Add(1)
		go func() {
			for cont := 0; cont < amount; cont += 1 {

				err = services.PublishMessage(proof2Publisher)
				if err != nil {
					cont -= 1
					continue
				}
			}
		}()
	}

	services.CreateProof()

	return nil
}
