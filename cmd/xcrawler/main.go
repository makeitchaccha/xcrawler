package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/subcommands"
	twitterscraper "github.com/imperatrona/twitter-scraper"
	"github.com/makeitchaccha/xcrawler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	scraper    *twitterscraper.Scraper
	repository xcrawler.Repository
)

func init() {
	// subcommand
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&crawlCmd{}, "")
	subcommands.Register(&addCmd{}, "")

	flag.Parse()
}

// prepare scraper
func prepareScraper() {
	scraper = twitterscraper.New()
	scraper.SetAuthToken(twitterscraper.AuthToken{
		Token:     os.Getenv("CRAWLER_X_TOKEN"),
		CSRFToken: os.Getenv("CRAWLER_X_CSRF_TOKEN"),
	})

	if !scraper.IsLoggedIn() {
		panic("Invalid AuthToken")
	}
}

// database connection
func prepareDatabase() {
	fmt.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(os.Getenv("CRAWLER_DB_STRING")), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// no migrations: its managed by goose.
	repository = xcrawler.NewRepository(db)
}

func main() {
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

var _ subcommands.Command = (*addCmd)(nil)

type addCmd struct {
}

func (*addCmd) Name() string { return "add" }

func (*addCmd) Synopsis() string { return "Add a user to the database." }

func (*addCmd) Usage() string {
	return `add <username>:
	Add a user to the database.
`
}

func (c *addCmd) SetFlags(f *flag.FlagSet) {}

func (*addCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		return subcommands.ExitUsageError
	}

	username := f.Arg(0)

	prepareScraper()
	prepareDatabase()

	// check if user exists
	user, err := scraper.GetProfile(username)
	if err != nil {
		fmt.Println("Failed to fetch user from twitter: ", err)
		return subcommands.ExitFailure
	}

	id, err := strconv.ParseUint(user.UserID, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse user ID: ", err)
		return subcommands.ExitFailure
	}

	repository.SaveUser(xcrawler.User{
		ID:       id,
		Username: user.Username, // use username from twitter, to match the capitalization.
	})

	fmt.Printf("User %s added to the database\n", user.Username)
	return subcommands.ExitSuccess
}

var _ subcommands.Command = (*crawlCmd)(nil)

type crawlCmd struct{}

func (*crawlCmd) Name() string     { return "crawl" }
func (*crawlCmd) Synopsis() string { return "Crawl the users registered on a database." }
func (*crawlCmd) Usage() string {
	return `crawl:
	Crawl the users registered on a database.
`
}

func (*crawlCmd) SetFlags(f *flag.FlagSet) {}
func (*crawlCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	prepareScraper()
	prepareDatabase()
	timestamp := time.Now().Truncate(time.Minute)

	users, err := repository.FindAllUsers()
	if err != nil {
		panic(err)
	}

	histories, missingUsers := processUsers(users, timestamp)

	repository.SaveHistories(histories)
	showResults(missingUsers)

	return subcommands.ExitSuccess
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
