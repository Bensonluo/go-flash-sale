package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://admin:admin@127.0.0.1:5672/Bensonl"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	QueueName string
	Exchange string
	Key string
	Mqurl string
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ {
		QueueName: queueName,
		Exchange: exchange,
		Key: key,
		Mqurl: MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "connection error!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "fail to obtain channel!")
	return rabbitmq
}

func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "",  "")

}

func (r *RabbitMQ) PublishSimple(message string) {
	//apply for queue
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})
}

func (r *RabbitMQ) ConsumeSimple() {
	//apply for queue
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.channel.Consume(
		  r.QueueName,
		  "",
		  true,
		  false,
		  false,
		  false,
		  nil)

	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//go routine
	go func() {
		for d := range msgs {
			//deal with msgs
			log.Printf("Received a message: %s", d.Body)
			fmt.Println(d.Body)
		}
	}()

	log.Printf("[*] waiting for messages, To exit, press CTRL + C" )
	<-forever
}

func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "failed to connect rabbitmq!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "failed to open a channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishPub(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnError(err, "Failed to declare!")

	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ReceiveSub() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare")
	q, err := r.channel.QueueDeclare(
		"", //random name
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare a queue")

	//bind to exchange
	err = r.channel.QueueBind(
		q.Name, //random name
		//empty key
		"",
		r.Exchange,
		false,
		nil)

	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("Exit by Clicking CTRL+C")
	<-forever
}


func NewRabbitMQRouting(exchangeName string,routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err,"failed to connect rabbitmq!")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "failed to open a channel")
	return rabbitmq
}


func (r *RabbitMQ) PublishRouting(message string )  {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", //required
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnError(err, "Failed to declare")

	err = r.channel.Publish(
		r.Exchange,
		r.Key, //required
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ReceiveRouting() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare")
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		r.Key, //required
		r.Exchange,
		false,
		nil)

	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("Exit by Clicking CTRL+C")
	<-forever
}


func NewRabbitMQTopic(exchangeName string,routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err,"failed to connect rabbitmq!")

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "failed to open a channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishTopic(message string )  {

	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic", //topic
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnError(err, "Failed to declare")


	err = r.channel.Publish(
		r.Exchange,
		r.Key, //key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}


//use * to match single word and # to match mutiple words
func (r *RabbitMQ) RecieveTopic() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare")

	q, err := r.channel.QueueDeclare(
		"", //random
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare a queue")

	err = r.channel.QueueBind(
		q.Name,
		//empty key
		r.Key,
		r.Exchange,
		false,
		nil)


	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("exit by clicking CTRL+C")
	<-forever
}
