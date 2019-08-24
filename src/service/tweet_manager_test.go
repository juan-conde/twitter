package service_test

import (
	"strings"
	"testing"

	"github.com/juan-conde/twitter/src/service"

	"github.com/juan-conde/twitter/src/domain"
)

var tweetWriter service.TweetWriter
var tweetManager *service.TweetManager

func setUp() {
	tweetWriter = service.NewMemoryTweetWriter() // Mock implementation
	tweetManager = service.NewTweetManager(tweetWriter)
}

func TestPublishedTweetIsSaved(t *testing.T) {
	setUp()
	var tweet domain.Tweet
	user := "grupoesfera"
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	tweetManager.PublishTweet(tweet)

	// Validation
	publishedTweets := tweetManager.GetTweets()
	publishedTweet := publishedTweets[0]
	if publishedTweet.GetUser() != user &&
		publishedTweet.GetText() != text {
		t.Errorf("Expected tweet is %s: %s \nbut is %s: %s",
			user, text, publishedTweet.GetUser(), publishedTweet.GetText())
	}
	if publishedTweet.GetDate() == nil {
		t.Error("Expected date can't be nil")
	}
}

func TestGetLastSavedTweet(t *testing.T) {
	setUp()
	var tweet domain.Tweet
	user := "grupoesfera"
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	tweetManager.PublishTweet(tweet)

	// Validation
	publishedTweet := tweetManager.GetTweet()
	if publishedTweet.GetUser() != user &&
		publishedTweet.GetText() != text {
		t.Errorf("Expected tweet is %s: %s \nbut is %s: %s",
			user, text, publishedTweet.GetUser(), publishedTweet.GetText())
	}
	if publishedTweet.GetDate() == nil {
		t.Error("Expected date can't be nil")
	}
}

func TestTweetWithoutUserIsNotPublished(t *testing.T) {
	setUp()
	var tweet domain.Tweet
	var user string
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	var err error
	_, err = tweetManager.PublishTweet(tweet)

	// Validation
	if err == nil {
		t.Error("Expected error did not appear")
	}
	if err != nil && err.Error() != "user is required" {
		t.Error("Expected error is user is required")
	}
}

func TestTweetWithoutTextIsNotPublished(t *testing.T) {
	setUp()
	var tweet domain.Tweet
	user := "Juan"
	var text string
	tweet = domain.NewTextTweet(user, text)
	var err error
	_, err = tweetManager.PublishTweet(tweet)

	// Validation
	if err == nil {
		t.Error("Expected error did not appear")
	}
	if err != nil && err.Error() != "text in tweet is required" {
		t.Error("Expected error is text in tweet is required")
	}
}

func TestTweetWhichExceeding140CharactersIsNotPublished(t *testing.T) {
	setUp()
	var tweet domain.Tweet
	user := "Juan"
	text := "AL GALLINERO YA SE LO PRENDIMO FUEGO, A SAN LORENZO LO CORRIMOS EN BOEDO, AVELLANEDA LO DEFIENDE UN POLICIA, HAY QUE PPPPP QUE SON LAS HINCHADAS UNIDAS"
	tweet = domain.NewTextTweet(user, text)
	var err error
	_, err = tweetManager.PublishTweet(tweet)

	// Validation
	if err == nil {
		t.Error("Expected error did not appear")
	}
	if err != nil && err.Error() != "text has over 140 characters" {
		t.Error("Expected error is text has over 140 characters")
	}
}

func TestCanPublishAndRetrieveMoreThanOneTweet(t *testing.T) {
	setUp()
	var tweet, secondTweet domain.Tweet
	user1 := "Juan"
	text1 := "Tweet 1 de Juan"
	tweet = domain.NewTextTweet(user1, text1)
	user2 := "Juan"
	text2 := "Tweet 2 de Juan"
	secondTweet = domain.NewTextTweet(user2, text2)
	tweetManager.PublishTweet(tweet)
	tweetManager.PublishTweet(secondTweet)

	// Validation
	publishedTweets := tweetManager.GetTweets()
	if len(publishedTweets) != 2 {
		t.Errorf("Expected size is 2 but was %d", len(publishedTweets))
		return
	}
	firstPublishedTweet := publishedTweets[0]
	secondPublishedTweet := publishedTweets[1]
	if !isValidTweet(t, firstPublishedTweet, 1, user1, text1) {
		return
	}
	if !isValidTweet(t, secondPublishedTweet, 2, user2, text2) {
		return
	}
}

func TestCanRetrieveTweetById(t *testing.T) {
	setUp()
	var tweet domain.Tweet
	var id int
	user := "grupoesfera"
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	id, _ = tweetManager.PublishTweet(tweet)

	// Validation
	publishedTweet := tweetManager.GetTweetByID(id)
	if !isValidTweet(t, publishedTweet, id, user, text) {
		return
	}
}

func TestCanRetrieveTheTweetsSentByAnUser(t *testing.T) {
	setUp()
	var tweet, secondTweet, thirdTweet domain.Tweet
	user := "grupoesfera"
	anotherUser := "nick"
	text := "This is my first tweet"
	secondText := "This is my second tweet"
	tweet = domain.NewTextTweet(user, text)             //id 1
	secondTweet = domain.NewTextTweet(user, secondText) // id 2
	thirdTweet = domain.NewTextTweet(anotherUser, text) // id 3
	tweetManager.PublishTweet(tweet)
	tweetManager.PublishTweet(secondTweet)
	tweetManager.PublishTweet(thirdTweet)

	// Operation
	tweets := tweetManager.GetTweetsByUser(user)
	countTweets := tweetManager.CountTweetsByUser(user)

	// Validation
	if countTweets != 2 { /* handle error */
		t.Errorf("Expected size is 2 but was %d", len(tweets))
		return
	}
	firstPublishedTweet := tweets[0]
	secondPublishedTweet := tweets[1]
	if !isValidTweet(t, firstPublishedTweet, 1, user, text) {
		return
	}
	if !isValidTweet(t, secondPublishedTweet, 2, user, secondText) {
		return
	}
}

func isValidTweet(t *testing.T, tweet domain.Tweet, id int, user string, text string) bool {
	if tweet.GetUser() == user && tweet.GetText() == text && tweet.GetID() == id {
		return true
	}
	return false
}

func TestPublishedTweetIsSavedToExternalResource(t *testing.T) {
	setUp()
	var tweet domain.Tweet // Fill the tweet with data
	user := "grupoesfera"
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	id, _ := tweetManager.PublishTweet(tweet)

	// Validation
	savedTweet := tweetWriter.GetAllTweets()[0]
	//    savedTweet := memoryWriter.GetLastSavedTweet()

	if savedTweet == nil {
		t.Errorf("NULL TWEET")
		return
	}
	if savedTweet.GetID() != id {
		t.Errorf("DIFFERENT ID")
		return
	}
}

func TestWrtieInFile(t *testing.T) {
	writer := service.NewFileTweetWriter()
	manager := service.NewTweetManager(writer)
	var tweet domain.Tweet // Fill the tweet with data
	user := "Juan"
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	// Operation
	id, _ := manager.PublishTweet(tweet)

	if id == 0 {
		t.Errorf("DIFFERENT ID")
		return
	}
}

func TestCanSearchForTweetContainingText(t *testing.T) {
	setUp()
	// Create and publish a tweet
	var tweet domain.Tweet
	user := "grupoesfera"
	text := "This is my first tweet"
	tweet = domain.NewTextTweet(user, text)
	tweetManager.PublishTweet(tweet)

	// Operation
	searchResult := make(chan service.Tweet)
	query := "first"
	tweetManager.SearchTweetsContaining(query, searchResult)

	// Validation
	foundTweet := <-searchResult //estoy leyendo del channel

	if foundTweet == nil {
		t.Errorf("NULL")
		return
	}
	if !strings.Contains(foundTweet.GetText(), query) {
		t.Errorf("Word not found in tweet")
		return
	}
}
