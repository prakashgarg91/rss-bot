package feedbot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/ItalyPaleAle/rss-bot/bot"
	"github.com/ItalyPaleAle/rss-bot/feeds"
	pb "github.com/ItalyPaleAle/rss-bot/proto"
	"github.com/ItalyPaleAle/rss-bot/utils"
)

// FeedBot is the class that manages the RSS bot
type FeedBot struct {
	log     *log.Logger
	manager *bot.BotManager
	feeds   *feeds.Feeds
	ctx     context.Context
	cancel  context.CancelFunc
}

// Init the object
func (fb *FeedBot) Init(manager *bot.BotManager) error {
	// Init the logger
	fb.log = log.New(os.Stdout, "feed-bot: ", log.Ldate|log.Ltime|log.LUTC)

	// Store the manager
	fb.manager = manager

	return nil
}

// Start the background workers
func (fb *FeedBot) Start() error {
	// Context, that can be used to stop the bot
	fb.ctx, fb.cancel = context.WithCancel(context.Background())

	// Init the feeds object
	fb.feeds = &feeds.Feeds{}
	err := fb.feeds.Init(fb.ctx)
	if err != nil {
		return err
	}

	// Start the background worker
	go fb.backgroundWorker()
	fb.log.Println("FeedBot workers started")

	// Register all commands
	err = fb.registerRoutes()
	if err != nil {
		return err
	}

	return nil
}

// Stop the background processes
func (fb *FeedBot) Stop() {
	fb.cancel()
}

// In background, start updating feeds periodically and send messages on new posts
// Also watch for the stop message
func (fb *FeedBot) backgroundWorker() {
	// Sleep for 2 seconds
	time.Sleep(2 * time.Second)

	// Channel for receiving messages to send
	msgCh := make(chan feeds.UpdateMessage)
	fb.feeds.SetUpdateChan(msgCh)

	// Queue an update right away
	fb.feeds.QueueUpdate()

	// Ticker for updates
	ticker := time.NewTicker(viper.GetDuration("FeedUpdateInterval") * time.Second)
	for {
		select {
		// On the interval, queue an update
		case <-ticker.C:
			fb.feeds.QueueUpdate()

		// Send messages on new posts
		case msg := <-msgCh:
			// This method logs errors already
			fb.sendFeedUpdate(&msg)

		// Context canceled
		case <-fb.ctx.Done():
			// Stop the ticker
			ticker.Stop()
			return
		}
	}
}

// Sends a message with a feed's post
func (fb *FeedBot) sendFeedUpdate(msg *feeds.UpdateMessage) {
	// Send the post
	_, err := fb.manager.SendMessage(&pb.OutMessage{
		Recipient: strconv.FormatInt(int64(msg.ChatId), 10),
		Content: &pb.OutMessage_Text{
			Text: &pb.OutTextMessage{
				Text:      fb.formatUpdateMessage(msg),
				ParseMode: pb.ParseMode_HTML,
			},
		},
		DisableWebPagePreview: true,
	})
	if err != nil {
		fb.log.Printf("Error sending message to chat %d: %s\n", msg.ChatId, err.Error())
		return
	}

	// Send photo, if any
	// Note that this might fail, for example if the image is too big (>5MB)
	if msg.Post.Photo != "" {
		_, err = fb.manager.SendMessage(&pb.OutMessage{
			Recipient: strconv.FormatInt(int64(msg.ChatId), 10),
			Content: &pb.OutMessage_Photo{
				Photo: &pb.OutPhotoMessage{
					File: &pb.OutFileMessage{
						Location: &pb.OutFileMessage_Url{
							Url: msg.Post.Photo,
						},
					},
				},
			},
			// Do not send notifications for subsequent messages
			DisableNotification: true,
		})
		if err != nil {
			// Just log the error and continue
			fb.log.Printf("Error sending photo %s to chat %d: %s\n", msg.Post.Photo, msg.ChatId, err.Error())
		}
	}
}

// Formats a message with an update
func (fb *FeedBot) formatUpdateMessage(msg *feeds.UpdateMessage) string {
	// Note: the msg.Feed object might be nil when passed to this method
	out := ""
	if msg.Feed != nil {
		out += fmt.Sprintf("🎙 %s:\n", utils.EscapeHTMLEntities(msg.Feed.Title))
	}

	// Add the content
	out += fmt.Sprintf("📬 <b>%s</b>\n🕓 %s\n🔗 %s\n",
		utils.EscapeHTMLEntities(msg.Post.Title),
		utils.EscapeHTMLEntities(msg.Post.Date.UTC().Format("Mon, 02 Jan 2006 15:04:05 MST")),
		utils.EscapeHTMLEntities(msg.Post.Link),
	)
	return out
}

// Sends a response to a command
// For commands sent in private chats, this just sends a regular message
// In groups, this replies to a specific message
/*func (b *FeedBot) respondToCommand(m *tb.Message, msg interface{}, options ...interface{}) (out *tb.Message, err error) {
	// If it's a private chat, send a message, otherwise reply
	if m.Private() {
		out, err = b.bot.Send(m.Sender, msg, options...)
	} else {
		out, err = b.bot.Reply(m, msg, options...)
	}

	// Log errors
	if err != nil {
		b.log.Printf("Error sending message to chat %d: %s\n", m.Chat.ID, err.Error())
	}

	return
}*/

// Register all routes
func (fb *FeedBot) registerRoutes() (err error) {
	fb.manager.AddRoute("(?i)^add feed (.*)", fb.routeAdd)

	/*// Handler for callbacks
	b.bot.Handle(tb.OnCallback, func(cb *tb.Callback) {
		// Seems that we need to trim whitespaces from the data
		data := strings.TrimSpace(cb.Data)
		// The main command comes before the /
		pos := strings.Index(data, "/")
		cmd := data
		var userData string
		if pos > -1 {
			cmd = data[0:pos]
			userData = data[(pos + 1):]
		}

		switch cmd {
		// Cancel command removes all inline keyboards
		case "cancel":
			_, err := b.bot.Edit(cb.Message, "Ok, I won't do anything")
			if err != nil {
				b.log.Printf("Error canceling callback: %s\n", err.Error())
			}

		// Confirm removing a feed
		case "confirm-remove":
			b.callbackConfirmRemove(cb, userData)
		}
	})*/

	return err
}