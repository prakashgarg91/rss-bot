package main

import (
	"regexp"

	"github.com/ItalyPaleAle/rss-bot/bots/feedbot/feeds"
	pb "github.com/ItalyPaleAle/rss-bot/model"
)

var routeAddMatch = regexp.MustCompile("(?i)^add feed (.*)")

// Route for the "add feed" command
func (fb *FeedBot) routeAdd(m *pb.InMessage) {
	// Get the URL
	match := routeAddMatch.FindStringSubmatch(m.Text)
	if len(match) < 2 {
		fb.client.RespondToCommand(m, "I didn't understand this \"add feed\" message - is the URL missing?")
		return
	}
	url := match[1]

	// Send a message that we're working on it
	sent, err := fb.client.RespondToCommand(m, "Working on it…")
	if err != nil {
		fb.log.Printf("Error sending message to chat %d: %s\n", m.ChatId, err.Error())
		return
	}

	// Add the subscription
	post, err := fb.feeds.AddSubscription(url, m.ChatId)
	if err != nil {
		if err == feeds.ErrAlreadySubscribed {
			err := fb.client.EditTextMessage(sent, &pb.OutTextMessage{
				Text: "You've already subscribed this chat to the feed",
			}, nil)
			if err != nil {
				// Log errors only
				fb.log.Printf("Error sending message to chat %d: %s\n", m.ChatId, err.Error())
			}
		} else {
			// Log errors and then send a message
			fb.log.Printf("Error while adding feed to chat %d: %s\n", m.ChatId, err.Error())

			err := fb.client.EditTextMessage(sent, &pb.OutTextMessage{
				Text: "An internal error occurred",
			}, nil)
			if err != nil {
				// Log errors only
				fb.log.Printf("Error sending message to chat %d: %s\n", m.ChatId, err.Error())
			}
		}
		return
	}

	err = fb.client.EditTextMessage(sent, &pb.OutTextMessage{
		Text: "I've added the feed to this channel. Here is the last post they published:",
	}, nil)
	if err != nil {
		fb.log.Printf("Error sending message to chat %d: %s\n", m.ChatId, err.Error())
		return
	}
	fb.sendFeedUpdate(&feeds.UpdateMessage{
		Post:   *post,
		ChatId: m.ChatId,
	})
}