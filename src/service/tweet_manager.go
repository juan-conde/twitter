package service

import (
	"errors"
	"strings"

	"github.com/juan-conde/twitter/src/domain"
)

// TweetManager struct that has all tweets and its methods
type TweetManager struct {
	tweets       ListOfTweets
	tweetsByUser TweetsByUser
	tweetWriter  TweetWriter
}

// Tweet type of tweet
type Tweet domain.Tweet

// ListOfTweets type of list of tweets
type ListOfTweets []domain.Tweet

// TweetsByUser type of map with key user and value list of its tweets
type TweetsByUser map[string]ListOfTweets

// NewTweetManager initialize a tweet manager
func NewTweetManager(tweetWriter TweetWriter) *TweetManager {
	tweetManager := new(TweetManager)
	tweetManager.tweets = make(ListOfTweets, 0)
	tweetManager.tweetsByUser = make(map[string]ListOfTweets)
	tweetManager.tweetWriter = tweetWriter
	return tweetManager
}

// GetTweet returns last tweet published
func (manager *TweetManager) GetTweet() Tweet {
	lastTweet := len(manager.tweets) - 1
	return manager.tweets[lastTweet]
}

// GetTweets returns list of tweets
func (manager *TweetManager) GetTweets() ListOfTweets {
	return manager.tweets
}

// GetTweetByID search a tweet by id
func (manager *TweetManager) GetTweetByID(id int) Tweet {
	return manager.tweets[id-1]
}

// GetTweetsByUser get user in map
func (manager *TweetManager) GetTweetsByUser(user string) ListOfTweets {
	return manager.tweetsByUser[user]
}

// CountTweetsByUser g
func (manager *TweetManager) CountTweetsByUser(user string) int {
	return len(manager.tweetsByUser[user])
}

// PublishTweet setear tweet
func (manager *TweetManager) PublishTweet(tweet domain.Tweet) (int, error) {
	if tweet.GetUser() == "" {
		return -1, errors.New("user is required")
	}
	if tweet.GetText() == "" {
		return -1, errors.New("text in tweet is required")
	}
	if len(tweet.GetText()) > 140 {
		return -1, errors.New("text has over 140 characters")
	}

	manager.tweets = append(manager.tweets, tweet)
	pos := len(manager.tweets)
	tweet.SetID(pos)

	elem, ok := manager.tweetsByUser[tweet.GetUser()]
	if ok {
		manager.tweetsByUser[tweet.GetUser()] = append(elem, tweet)
	} else {
		manager.tweetsByUser[tweet.GetUser()] = []domain.Tweet{tweet}
	}

	manager.tweetWriter.SaveTweet(tweet)
	return pos, nil
}

// SearchTweetsContaining search for word in all tweets and write to channel
func (manager *TweetManager) SearchTweetsContaining(query string, c chan Tweet) {
	go manager.saveTweetInChannel(query, c)
}

func (manager *TweetManager) saveTweetInChannel(query string, c chan Tweet) {
	tweets := manager.tweets
	tweetsSize := len(tweets)

	for i := 0; i < tweetsSize; i++ {
		if strings.Contains(tweets[i].GetText(), query) {
			c <- tweets[i]
		}
	}
}
