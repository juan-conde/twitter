package domain

import "time"

// Tweet tweet interface gonna be implemented by every type of tweet
type Tweet interface {
	PrintableTweet() string
	GetUser() string
	GetText() string
	GetID() int
	GetDate() *time.Time
	SetID(int)
}

// TextTweet formato del Text Tweet
type TextTweet struct {
	User string
	Text string
	//Lo uso como puntero porque siempre que trabajo con estructuras conviene el puntero
	Date *time.Time
	ID   int
}

// ImageTweet formato del Image Tweet
type ImageTweet struct {
	TextTweet
	URL string
}

// QuoteTweet formato del Quote Tweet
type QuoteTweet struct {
	TextTweet
	TweetCitado Tweet
}

// NewTextTweet devuelvo una estructura tweer con user y text
func NewTextTweet(user, text string) *TextTweet {
	date := time.Now()
	textTweet := TextTweet{
		user,
		text,
		&date,
		-2,
	}
	return &textTweet
}

// NewImageTweet a
func NewImageTweet(user, text, url string) *ImageTweet {
	date := time.Now()
	textTweet := TextTweet{
		user,
		text,
		&date,
		-2,
	}
	imageTweet := ImageTweet{
		textTweet,
		url,
	}
	return &imageTweet
}

// NewQuoteTweet a
func NewQuoteTweet(user, text string, quotedTweet Tweet) *QuoteTweet {
	date := time.Now()
	textTweet := TextTweet{
		user,
		text,
		&date,
		-2,
	}
	finalTweet := QuoteTweet{
		textTweet,
		quotedTweet,
	}
	return &finalTweet
}

func (textTweet *TextTweet) PrintableTweet() string {
	text := "@" + textTweet.User + ": " + textTweet.Text
	return text
}

func (imageTweet *ImageTweet) PrintableTweet() string {
	text := "@" + imageTweet.User + ": " + imageTweet.Text + " - URL: " + imageTweet.URL
	return text
}

func (quoteTweet *QuoteTweet) PrintableTweet() string {
	text := quoteTweet.TextTweet.PrintableTweet() + " - " + quoteTweet.TweetCitado.PrintableTweet()
	return text
}

func (textTweet *TextTweet) GetUser() string { return textTweet.User }

func (textTweet *TextTweet) GetText() string { return textTweet.Text }

func (textTweet *TextTweet) GetID() int { return textTweet.ID }

func (textTweet *TextTweet) GetDate() *time.Time { return textTweet.Date }

func (textTweet *TextTweet) SetID(id int) { textTweet.ID = id }
