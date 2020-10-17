package feeds

import (
	"github.com/Songmu/go-httpdate"

	"github.com/0x111/telegram-rss-bot/db"
	"github.com/0x111/telegram-rss-bot/models"
)

// Number of parallel requests to make
const parallelFetch = 4

// SetUpdateChan sets the channel to use for notify the bot of new messages for subscribers
func (f *Feeds) SetUpdateChan(ch chan<- UpdateMessage) {
	f.updateCh = ch
}

// QueueUpdate queues an update of the feeds
func (f *Feeds) QueueUpdate() {
	// The channel has a capacity of 1, which means that there can only be 1 running and one queued
	// This is so we don't have updates running in parallel, nor a situation in which updates are queued faster than they are completed
	select {
	// If there's already one request waiting, then return right away
	case f.waiting <- 1:
		break
	default:
		return
	}

	// Acquire the lock now (wait till we can) and then release the waiting lock
	f.semaphore <- 1
	<-f.waiting

	// Update the feeds in background
	// This is so the QueueUpdate method can return
	go func() {
		err := f.updateFeeds()
		if err != nil {
			f.log.Println("Error while updating feeds", err)
		}

		// Release the lock
		<-f.semaphore
	}()
}

type workerResult struct {
	Feed  *models.Feed
	Posts []Post
}

// Internal worker that fetches and processes feeds, in parallel
func (f *Feeds) updateWorker(id int, jobs <-chan *models.Feed, results chan<- workerResult) {
	for j := range jobs {
		res := workerResult{
			Feed: j,
		}
		f.log.Println("Worker", id, "started updating feed", j.ID)
		// Fetch new data from the feed
		posts, err := f.fetchFeed(j)
		if err != nil {
			// Error is already logged
			// Just move to the next post
			results <- res
			continue
		}
		res.Posts = posts
		f.log.Println("Worker", id, "finished updating feed", j.ID)
		results <- res
	}
}

// Worker that updates all feeds
func (f *Feeds) updateFeeds() error {
	f.log.Println("Started updating feeds")

	// Start background workers to parallelize requests
	jobs := make(chan *models.Feed)
	results := make(chan workerResult)
	for i := 1; i <= parallelFetch; i++ {
		go f.updateWorker(i, jobs, results)
	}

	// Select all feeds
	count := 0
	rows, err := db.GetDB().Queryx("SELECT * FROM feeds")
	if err != nil {
		rows.Close()
		return err
	}
	for rows.Next() {
		// If the context was canceled, return
		if err := f.ctx.Err(); err != nil {
			rows.Close()
			close(jobs)
			close(results)
			return err
		}

		// Read the row
		feed := models.Feed{}
		err = rows.StructScan(&feed)
		if err != nil {
			rows.Close()
			close(jobs)
			close(results)
			return err
		}

		// Get a worker to perform the request
		jobs <- &feed
		count++
	}
	rows.Close()
	if err != nil {
		close(jobs)
		close(results)
		return err
	}

	// Close the jobs channel
	close(jobs)

	// Read the results
	for i := 0; i < count; i++ {
		res := <-results

		// If there are new posts…
		if len(res.Posts) > 0 {
			// …first, update the feed object in the database
			f.setLastPost(res.Feed)

			// …second, notify subscribers
			// Ignore errors (already logged)
			_ = f.notifySubscribers(res.Feed, res.Posts)
		}
	}
	close(results)

	f.log.Println("Done updating feeds")

	return nil
}

// Fetches a feed and return the new posts only
// If there are new posts, the feed object is updated too as a side effect
func (f *Feeds) fetchFeed(feed *models.Feed) ([]Post, error) {
	// Request the data
	f.log.Printf("Updating feed %d (%s)\n", feed.ID, feed.Url)
	posts, err := f.RequestFeed(feed)
	if err != nil {
		f.log.Println("Error while fetching the feed:", err)
		return nil, err
	}

	// Get all new entries
	res := make([]Post, 0)
	if posts != nil && posts.Items != nil {
		after := feed.LastPostDate
		for _, el := range posts.Items {
			// Skip items with an invalid date
			parsePubDate, err := httpdate.Str2Time(el.Published, nil)
			if err != nil {
				f.log.Printf("Error in feed %s: skipping entry with invalid date '%s' (error: %s)\n", feed.Url, el.Published, err)
				continue
			}
			if el.Title == "" {
				f.log.Printf("Error in feed %s: skipping entry with empty title\n", feed.Url)
				continue
			}

			// Check if this is a new post
			if parsePubDate.After(after) {
				// Add it to the result
				res = append(res, Post{
					Title: el.Title,
					Link:  el.Link,
					Date:  parsePubDate,
				})

				// Look for the most recent post for updating the feed object
				if parsePubDate.After(feed.LastPostDate) {
					feed.LastPostTitle = el.Title
					feed.LastPostLink = el.Link
					feed.LastPostDate = parsePubDate
				}
			}
		}
	}

	return res, nil
}

// Update a feed in the database, setting the new details for the last post
// This doesn't return errors but it only logs them
func (f *Feeds) setLastPost(feed *models.Feed) {
	f.log.Printf("Updating last post for feed %d\n", feed.ID)

	// Note that we're not using a transaction here (because the update process can take a while), but there's only one of these methods that can be running at the same time
	// The bot can be deleting the feed in the meanwhile, but this would just make the next query fail (and that's why we're ignoring the error here)
	_, err := db.GetDB().Exec("UPDATE feeds SET feed_last_modified = ?, feed_etag = ?, feed_last_post_title = ?, feed_last_post_link = ?, feed_last_post_date = ? WHERE feed_id = ?", feed.LastModified, feed.ETag, feed.LastPostTitle, feed.LastPostLink, feed.LastPostDate, feed.ID)
	if err != nil {
		f.log.Printf("Error while updating the last post for feed %s (id: %d), but continuing to next. Error: %s\n", feed.Url, feed.ID, err)
	}
}

// Sends a notification to all subscribers when a new post is out
func (f *Feeds) notifySubscribers(feed *models.Feed, posts []Post) error {
	// Get the list of subscribers for this feed
	sub := &models.Subscription{}
	rows, err := db.GetDB().Queryx("SELECT chat_id FROM subscriptions WHERE feed_id = ?", feed.ID)
	defer rows.Close()
	if err != nil {
		f.log.Println("Error querying the database:", err)
		return err
	}
	subCount := 0
	for rows.Next() {
		// Read the row
		err = rows.StructScan(&sub)
		if err != nil {
			f.log.Println("Error reading a row:", err)
			return err
		}

		// Send the message to the channel
		for _, post := range posts {
			f.updateCh <- UpdateMessage{
				Feed:   feed,
				Post:   post,
				ChatId: sub.ChatID,
			}
		}
		subCount++
	}

	f.log.Printf("Found %d new posts in feed id %d, and notified %d subscribers\n", len(posts), feed.ID, subCount)

	return nil
}