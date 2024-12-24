package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/4madness7/gator/internal/database"
	"github.com/google/uuid"
)

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("'login' expects 1 <username> parameter.")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User %s does not exist in the database.\nError: %w", username, err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User as been set to '%s'\n", username)
	return nil
}

func registerHander(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("'register' expects 1 <username> parameter.")
	}
	current_time := time.Now()
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      cmd.args[0],
	}

	newUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}

	fmt.Println("User created and logged successfully.")
	fmt.Println("===== DEBUG =====")
	fmt.Printf(`User Data {
    id:         %v
    created at: %v
    updated at: %v
    name:       %s
}`,
		newUser.ID, newUser.CreatedAt, newUser.UpdatedAt, newUser.Name)
	fmt.Println()

	return nil
}

func resetHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("'reset' does not expect any arguments.")
	}
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not reset DB: %w", err)
	}

	fmt.Println("Datadase reset successfully.")
	return nil
}

func usersHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("'users' does not expect any arguments.")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("No users in the database")
	}

	for _, user := range users {
		currentStr := ""
		if s.cfg.CurrentUserName == user.Name {
			currentStr = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, currentStr)
	}
	return nil
}

func aggHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("'agg' expects 1 arguments. Ex. gator agg <time_between_reqs>.")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return nil
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func addfeedHander(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("'addfeed' expects 2 arguments. Ex. gator addfeed <name> <url>.")
	}
	current_time := time.Now()
	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	fmt.Println("Feed created and followed successfully.")
	fmt.Println("===== DEBUG =====")
	fmt.Printf(`Feed Data {
    id:         %v
    created at: %v
    updated at: %v
    name:       %s
    url:        %s
    user id:    %v
}`,
		newFeed.ID, newFeed.CreatedAt, newFeed.UpdatedAt, newFeed.Name, newFeed.Url, newFeed.UserID)
	fmt.Println()

	return nil
}

func feedsHandler(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("'feeds' does not expect any arguments.")
	}
	rows, err := s.db.GetFeedsWithUser(context.Background())
	if err != nil {
		return err
	}

	for _, row := range rows {
		fmt.Printf("User: %s | Feed Name: %s | Feed Url: %s\n", row.UserName, row.FeedName, row.Url)
	}
	return nil
}

func followHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("'follow' expects 1 argument. Ex. gator follow <url>.")
	}
	current_time := time.Now()
	feed, err := s.db.GetFeedWithURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully followed '%s'\n", feed.Name)
	fmt.Println("===== DEBUG =====")
	fmt.Printf(`FeedFollow Data {
    id:         %v
    created at: %v
    updated at: %v
    user id:    %v
    feed id:    %v
}`,
		newFeedFollow.ID, newFeedFollow.CreatedAt, newFeedFollow.UpdatedAt, newFeedFollow.UserID, newFeedFollow.FeedID)
	fmt.Println()

	return nil
}

func followingHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("'following' does not expect any arguments")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Followed feeds:")
	for _, follow := range follows {
		fmt.Printf(" - %s\n", follow.FeedName)
	}

	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("'unfollow' expects 1 argument. Ex. gator unfollow <url>.")
	}

	feed, err := s.db.GetFeedWithURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' unfollowed successfully.\n", feed.Name)

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	current_time := time.Now()
	err = s.db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID:        feed.ID,
			UpdatedAt: current_time,
		},
	)
	if err != nil {
		return err
	}
	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Println("=== Feed fetched ===")
	fmt.Printf("Title: %s\n", rss.Channel.Title)
	fmt.Printf("Link: %s\n", rss.Channel.Link)
	fmt.Printf("Description: %s\n", rss.Channel.Description)
	fmt.Printf("Items: \n")
	for _, item := range rss.Channel.Item {
		fmt.Printf("  - %s\n", item.Title)
	}
	return nil
}
