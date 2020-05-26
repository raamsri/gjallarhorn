package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Receive message from Queue with long polling enabled.
//
// Usage:
//    go run sqs_longpolling_receive_message.go -n queue_name -t timeout
func main() {
	namePtr := flag.String("n", "", "Queue name")
	timeoutPtr := flag.Int64("t", 20, "(Optional) Timeout in seconds for long polling")

	flag.Parse()

	if *namePtr == "" {
		flag.PrintDefaults()
		exitErrorf("Queue name required")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	// Need to convert the queue name into a URL. Make the GetQueueUrl
	// API call to retrieve the URL. This is needed for receiving messages
	// from the queue.
	resultURL, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: namePtr,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			exitErrorf("Unable to find queue %q.", *namePtr)
		}
		exitErrorf("Unable to queue %q, %v.", *namePtr, err)
	}

	for {

		// Receive a message from the SQS queue with long polling enabled.
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: resultURL.QueueUrl,
			AttributeNames: aws.StringSlice([]string{
				"SentTimestamp",
			}),
			MaxNumberOfMessages: aws.Int64(1),
			MessageAttributeNames: aws.StringSlice([]string{
				"All",
			}),
			WaitTimeSeconds: timeoutPtr,
		})
		if err != nil {
			exitErrorf("Unable to receive message from queue %q, %v.", *namePtr, err)
		}

		fmt.Printf("Received %d messages.\n", len(result.Messages))
		if len(result.Messages) > 0 {
			// fmt.Println(result.Messages)
			for _, message := range result.Messages {
				fmt.Println(message)
				_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      resultURL.QueueUrl,
					ReceiptHandle: message.ReceiptHandle,
				})

				if err != nil {
					fmt.Println("Delete Error", err)
					return
				}

				fmt.Printf("Message %s deleted\n", *message.MessageId)
			}
		}

	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
