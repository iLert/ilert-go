package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	// for the channel type "ALERT" the channel id is the id of the alert
	channelID := int64(0) // your specific alert id
	channelType := &ilert.ChatChannelType.Alert

	created, err := client.CreateChatChannelMessage(&ilert.CreateChatChannelMessageInput{
		ChannelID:   &channelID,
		ChannelType: channelType,
		ChatChannelMessage: &ilert.ChatChannelMessage{
			Content:     "deploy rolled back, watching the error rate",
			ContentType: ilert.ChatChannelMessageContentType.Text,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Message:\n\n %+v\n", created.ChatChannelMessage)

	// reactions are written through their own endpoints and come back on the message
	reacted, err := client.AddChatChannelMessageReaction(&ilert.AddChatChannelMessageReactionInput{
		ChannelID:            &channelID,
		ChannelType:          channelType,
		ChatChannelMessageID: &created.ChatChannelMessage.ID,
		Code:                 ilert.String("thumbsup"),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reactions: %+v\n", reacted.ChatChannelMessage.Reactions)

	unreacted, err := client.RemoveChatChannelMessageReaction(&ilert.RemoveChatChannelMessageReactionInput{
		ChannelID:            &channelID,
		ChannelType:          channelType,
		ChatChannelMessageID: &created.ChatChannelMessage.ID,
		Code:                 ilert.String("thumbsup"),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reactions after removing: %+v\n", unreacted.ChatChannelMessage.Reactions)

	// only the author of a message may edit it, the API answers 403 otherwise
	edited, err := client.UpdateChatChannelMessage(&ilert.UpdateChatChannelMessageInput{
		ChannelID:            &channelID,
		ChannelType:          channelType,
		ChatChannelMessageID: &created.ChatChannelMessage.ID,
		ChatChannelMessage: &ilert.ChatChannelMessage{
			Content:     "deploy rolled back, error rate recovering",
			ContentType: ilert.ChatChannelMessageContentType.Text,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Edited: %t\n", edited.ChatChannelMessage.Edited)

	single, err := client.GetChatChannelMessage(&ilert.GetChatChannelMessageInput{
		ChannelID:            &channelID,
		ChannelType:          channelType,
		ChatChannelMessageID: &created.ChatChannelMessage.ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Message by id:\n\n %+v\n", single.ChatChannelMessage)

	// the first reply turns the message into a thread
	reply, err := client.CreateChatChannelThreadReply(&ilert.CreateChatChannelThreadReplyInput{
		ChannelID:   &channelID,
		ChannelType: channelType,
		ThreadID:    &created.ChatChannelMessage.ID,
		ChatChannelThreadReply: &ilert.ChatChannelThreadReply{
			Content:     "error rate back to normal",
			ContentType: ilert.ChatChannelMessageContentType.Text,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reply:\n\n %+v\n", reply.ChatChannelThreadReply)

	replies, err := client.GetChatChannelThreadReplies(&ilert.GetChatChannelThreadRepliesInput{
		ChannelID:   &channelID,
		ChannelType: channelType,
		ThreadID:    &created.ChatChannelMessage.ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Found %d replies\n", len(replies.ChatChannelThreadReplies))

	singleReply, err := client.GetChatChannelThreadReply(&ilert.GetChatChannelThreadReplyInput{
		ChannelID:                &channelID,
		ChannelType:              channelType,
		ThreadID:                 &created.ChatChannelMessage.ID,
		ChatChannelThreadReplyID: &reply.ChatChannelThreadReply.ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reply by id:\n\n %+v\n", singleReply.ChatChannelThreadReply)

	editedReply, err := client.UpdateChatChannelThreadReply(&ilert.UpdateChatChannelThreadReplyInput{
		ChannelID:                &channelID,
		ChannelType:              channelType,
		ThreadID:                 &created.ChatChannelMessage.ID,
		ChatChannelThreadReplyID: &reply.ChatChannelThreadReply.ID,
		ChatChannelThreadReply: &ilert.ChatChannelThreadReply{
			Content:     "error rate back to normal, closing this out",
			ContentType: ilert.ChatChannelMessageContentType.Text,
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reply edited: %t\n", editedReply.ChatChannelThreadReply.Edited)

	reactedReply, err := client.AddChatChannelThreadReplyReaction(&ilert.AddChatChannelThreadReplyReactionInput{
		ChannelID:                &channelID,
		ChannelType:              channelType,
		ThreadID:                 &created.ChatChannelMessage.ID,
		ChatChannelThreadReplyID: &reply.ChatChannelThreadReply.ID,
		Code:                     ilert.String("eyes"),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reply reactions: %+v\n", reactedReply.ChatChannelThreadReply.Reactions)

	if _, err := client.RemoveChatChannelThreadReplyReaction(&ilert.RemoveChatChannelThreadReplyReactionInput{
		ChannelID:                &channelID,
		ChannelType:              channelType,
		ThreadID:                 &created.ChatChannelMessage.ID,
		ChatChannelThreadReplyID: &reply.ChatChannelThreadReply.ID,
		Code:                     ilert.String("eyes"),
	}); err != nil {
		log.Fatalln("ERROR:", err)
	}

	// replies are deleted softly as well
	deletedReply, err := client.DeleteChatChannelThreadReply(&ilert.DeleteChatChannelThreadReplyInput{
		ChannelID:                &channelID,
		ChannelType:              channelType,
		ThreadID:                 &created.ChatChannelMessage.ID,
		ChatChannelThreadReplyID: &reply.ChatChannelThreadReply.ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Reply deleted: %t\n", deletedReply.ChatChannelThreadReply.Deleted)

	messages, err := client.GetChatChannelMessages(&ilert.GetChatChannelMessagesInput{
		ChannelID:   &channelID,
		ChannelType: channelType,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	for _, message := range messages.ChatChannelMessages {
		// ilert uses other styles internally, for example for the steps of its AI agents,
		// so only render the default ones in your own interface
		if message.MessageStyle != ilert.ChatChannelMessageStyle.Default {
			continue
		}
		log.Printf("%s: %s\n", message.AuthorType, message.Content)
	}

	// deletes are soft: the message stays in the list with Deleted set
	deleted, err := client.DeleteChatChannelMessage(&ilert.DeleteChatChannelMessageInput{
		ChannelID:            &channelID,
		ChannelType:          channelType,
		ChatChannelMessageID: &created.ChatChannelMessage.ID,
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Deleted: %t, content is now %q\n", deleted.ChatChannelMessage.Deleted, deleted.ChatChannelMessage.Content)
}
