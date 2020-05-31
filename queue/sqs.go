package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ReceiveMessageInput provides input parameters to receive message
type ReceiveMessageInput sqs.ReceiveMessageInput

// DeleteMessageInput provides input parameters to delete message
type DeleteMessageInput sqs.DeleteMessageInput

// NewSQSClient creates a new instance of the SQS client with a session
// Inputs:
//		None
// Output:
// 		If success, an SQS client.
func NewSQSClient() *sqs.SQS {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	return svc
}

// GetSQSURL gets the URL of an Amazon SQS queue
// Inputs:
//		svc is the sqs client
//		queueName is the name of the queue
// Output:
//		If success, the URL of the queue and nil
//		Otherwise, an empty string and an error from the call to
func GetSQSURL(svc *sqs.SQS, queue *string) (*sqs.GetQueueUrlOutput, error) {

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return urlResult, nil
}

// GetSQSMessages gets the messages from an Amazon SQS queue
// Inputs:
//		svc is the sqs client
//		ReceiveMessageInput provides input parameters to receive message
// Output:
//		If success, the latest message and nil
//		Otherwise, nil and an error from the call to ReceiveMessage
func GetSQSMessages(svc *sqs.SQS, input *ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            input.QueueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     input.WaitTimeSeconds,
	})
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

// DeleteSQSMessage deletes the message from an Amazon SQS queue
// Inputs:
//		svc is the sqs client
//		DeleteMessageInput provides input parameters to delete message
// Output:
//		If success, nill
//		Otherwise, an error from the call to DeleteMessage
func DeleteSQSMessage(svc *sqs.SQS, input *DeleteMessageInput) error {
	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      input.QueueUrl,
		ReceiptHandle: input.ReceiptHandle,
	})

	return err
}
