package main

import (
	"flag"
	"fmt"
	"os"

	"toppr/gjallarhorn/blow/event"
	"toppr/gjallarhorn/queue"
)

// Receive message from Queue with long polling enabled.
//
// Usage:
// 		go run sqs_longpolling_receive_message.go -n queue_name -t timeout
func main() {
	queueName := flag.String("n", "", "Queue name")
	timeoutPtr := flag.Int64("t", 20, "(Optional) Timeout in seconds for long polling")

	flag.Parse()

	if *queueName == "" {
		flag.PrintDefaults()
		exitErrorf("Queue name required")
	}

	svc := queue.NewSQSClient()

	// Need to convert the queue name into a URL. Make the getSQSUrl
	// API call to retrieve the URL. This is needed for receiving messages
	// from the queue.
	// Get URL of queue
	resultURL, err := queue.GetSQSURL(svc, queueName)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	qURL := resultURL.QueueUrl
	qReceiveMsgConfig := &queue.ReceiveMessageInput{
		WaitTimeSeconds: timeoutPtr,
		QueueUrl:        qURL,
	}
	qDeleteMsgConfig := &queue.DeleteMessageInput{
		QueueUrl: qURL,
	}

	for {
		// Receive a message from the SQS queue with long polling enabled.
		result, err := queue.GetSQSMessages(svc, qReceiveMsgConfig)
		if err != nil {
			fmt.Printf("Unable to receive message from queue %q, %v.", *qURL, err)
		}

		fmt.Printf("Received %d messages.\n", len(result.Messages))

		if len(result.Messages) < 1 {
			continue
		}
		// Process the messages
		for _, message := range result.Messages {

			err := event.ProcessMessage(fmt.Sprintf("%s", *message.Body))
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("Message posted")
			}

			qDeleteMsgConfig.ReceiptHandle = message.ReceiptHandle
			// Delete the message
			err = queue.DeleteSQSMessage(svc, qDeleteMsgConfig)
			if err != nil {
				fmt.Println("Delete Error", err)
				continue
			}
			fmt.Printf("Message %s deleted\n", *message.MessageId)
		}

	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
