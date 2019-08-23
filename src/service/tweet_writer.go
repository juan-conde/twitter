package service

import (
	"fmt"
	"os"
)

//TweetWriter interface
type TweetWriter interface {
	GetAllTweets() ListOfTweets
	SaveTweet(Tweet)
}

//MemoryTweetWriter memory writer
type MemoryTweetWriter struct {
	Tweets ListOfTweets
}

// NewMemoryTweetWriter constructor
func NewMemoryTweetWriter() *MemoryTweetWriter {
	return new(MemoryTweetWriter)
}

// GetAllTweets get all tweets in memory
func (writer *MemoryTweetWriter) GetAllTweets() ListOfTweets {
	return writer.Tweets
}

// SaveTweet save tweet in memory
func (writer *MemoryTweetWriter) SaveTweet(tweet Tweet) {
	writer.Tweets = append(writer.Tweets, tweet)
}

// FileTweetWriter writer in file
type FileTweetWriter struct {
	file *os.File
}

// NewFileTweetWriter file writer constructor
func NewFileTweetWriter() *FileTweetWriter {
	file, err := os.OpenFile(
		"/Users/jconde/Desktop/fileWriter",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	// file, err := os.Create("/Users/jconde/Desktop/fileWriter")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fileWriter := new(FileTweetWriter)
	fileWriter.file = file
	return fileWriter
}

// SaveTweet save tweet in file
func (writer *FileTweetWriter) SaveTweet(tweet Tweet) {
	go writeString(writer, tweet)
}

func writeString(writer *FileTweetWriter, tweet Tweet) {
	writer.file.WriteString(tweet.PrintableTweet())
}

// GetAllTweets get all tweets in file
func (writer *FileTweetWriter) GetAllTweets() ListOfTweets {
	return nil
}
