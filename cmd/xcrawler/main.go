package main

import (
	"fmt"
	"os"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
	"github.com/makeitchaccha/xcrawler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	scraper    *twitterscraper.Scraper
	repository xcrawler.Repository
)

// prepare scraper
func init() {
	scraper = twitterscraper.New()
	scraper.SetAuthToken(twitterscraper.AuthToken{
		Token:     os.Getenv("CRAWLER_X_TOKEN"),
		CSRFToken: os.Getenv("CRAWLER_X_CSRF_TOKEN"),
	})
}

// database connection
func init() {
	fmt.Println("Connecting to database...")
	fmt.Println("DB_STRING: ", os.Getenv("CRAWLER_DB_STRING"))
	db, err := gorm.Open(postgres.Open(os.Getenv("CRAWLER_DB_STRING")), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// no migrations: its managed by goose.
	repository = xcrawler.NewRepository(db)
}

func main() {
	timestamp := time.Now().Truncate(time.Minute)

	users, err := repository.FindAllUsers()
	if err != nil {
		panic(err)
	}

	histories, missingUsers := processUsers(users, timestamp)

	repository.SaveHistories(histories)
	showResults(missingUsers)
}

func processUsers(users []xcrawler.User, timestamp time.Time) ([]xcrawler.History, []string) {
	missingUsers := make([]string, 0)
	histories := make([]xcrawler.History, 0)

	for _, user := range users {
		strID := fmt.Sprintf("%d", user.ID)
		profile, err := scraper.GetProfileByID(strID)
		if err != nil {
			missingUsers = append(missingUsers, fmt.Sprintf("%d with error: %s", user.ID, err.Error()))
			continue
		}

		history := xcrawler.History{
			ID:             user.ID,
			CreatedAt:      timestamp,
			FavoritesCount: int64(profile.LikesCount),
			TweetCount:     int64(profile.TweetsCount),
		}

		histories = append(histories, history)
	}

	return histories, missingUsers
}

func showResults(missingUsers []string) {
	if len(missingUsers) > 0 {
		fmt.Println("Missing users: ", missingUsers)
		fmt.Println("Please check if the user is deleted or private.")
	} else {
		fmt.Println("Crawling completed successfully")
	}
}
