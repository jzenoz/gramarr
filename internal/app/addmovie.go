package app

import "C"
import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/tommy647/gramarr/internal/users"

	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

func NewAddMovieConversation(e *Service) *AddMovieConversation {
	return &AddMovieConversation{env: e, bot: e.Bot}
}

type AddMovieConversation struct {
	currentStep            func(interface{})
	movieQuery             string
	movieResults           []radarr.Movie
	folderResults          []radarr.Folder
	selectedMovie          *radarr.Movie
	selectedQualityProfile *radarr.Profile
	selectedFolder         *radarr.Folder
	env                    *Service
	bot                    Bot
}

func (c *AddMovieConversation) Run(m interface{}) {
	c.currentStep = c.AskMovie(m)
}

func (c *AddMovieConversation) Name() string {
	return "addmovie"
}

func (c *AddMovieConversation) CurrentStep() func(interface{}) {
	return c.currentStep
}

func (c *AddMovieConversation) AskMovie(m interface{}) func(interface{}) {
	user := users.User{}                                          // @todo: fix
	_ = c.bot.Send(user, "What movie do you want to search for?") // @todo: handle error

	return func(m interface{}) {
		c.movieQuery = c.bot.GetText(m)
		if c.movieQuery == "" {
			log.Println("empty message??")
			return
		}

		movies, err := c.env.Radarr.SearchMovies(c.movieQuery)
		c.movieResults = movies

		// Search Service Failed
		if err != nil {
			_ = c.bot.Send(user, "Failed to search movies.") // @todo: handle error
			c.env.CM.StopConversation(c)
			return
		}

		// No Results
		if len(movies) == 0 {
			msg := fmt.Sprintf("No movie found with the title '%s'", util.EscapeMarkdown(c.movieQuery))
			_ = c.bot.Send(user, msg) // @todo: handle error
			c.env.CM.StopConversation(c)
			return
		}

		// Found some movies! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d movies:*", len(movies)))
		for i, movie := range movies {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(movie.String())))
		}
		_ = c.bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error
		c.currentStep = c.AskPickMovie(m)
	}
}

func (c *AddMovieConversation) AskPickMovie(m interface{}) func(interface{}) {
	user := users.User{} // @todo: fix
	// Send custom reply keyboard
	var options []string
	for _, movie := range c.movieResults {
		options = append(options, fmt.Sprintf("%s", movie))
	}
	options = append(options, "/cancel")
	_ = c.bot.SendKeyboardList(user, "Which one would you like to download?", options) // @todo: handle error

	return func(m interface{}) {
		// Set the selected movie
		for i, opt := range options {
			if c.bot.GetText(m) == opt {
				c.selectedMovie = &c.movieResults[i]
				break
			}
		}

		// Not a valid movie selection
		if c.selectedMovie == nil {
			_ = c.bot.Send(user, "Invalid selection.") // @todo: handle error
			c.currentStep = c.AskPickMovie(m)
			return
		}

		c.currentStep = c.AskPickMovieQuality(m)
	}
}

func (c *AddMovieConversation) AskPickMovieQuality(m interface{}) func(interface{}) {
	user := users.User{} // @todo: fix
	profiles, err := c.env.Radarr.GetProfile("profile")

	// GetProfile Service Failed
	if err != nil {
		_ = c.bot.Send(user, "Failed to get quality profiles.") // @todo: handle error
		c.env.CM.StopConversation(c)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, QualityProfile := range profiles {
		options = append(options, fmt.Sprintf("%v", QualityProfile.Name))
	}
	options = append(options, "/cancel")
	_ = c.bot.SendKeyboardList(user, "Which quality shall I look for?", options) // @todo: handle error

	return func(m interface{}) {
		// Set the selected option
		for i := range options {
			if c.bot.GetText(m) == options[i] {
				c.selectedQualityProfile = &profiles[i]
				break
			}
		}

		// Not a valid selection
		if c.selectedQualityProfile == nil {
			_ = c.bot.Send(user, "Invalid selection.") // @todo: handle error
			c.currentStep = c.AskPickMovieQuality(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddMovieConversation) AskFolder(m interface{}) func(interface{}) {
	user := users.User{} // @todo: fix
	folders, err := c.env.Radarr.GetFolders()
	c.folderResults = folders

	// GetFolders Service Failed
	if err != nil {
		_ = c.bot.Send(user, "Failed to get folders.") // @todo: handle error
		c.env.CM.StopConversation(c)
		return nil
	}

	// No Results
	if len(folders) == 0 {
		_ = c.bot.Send(user, "No destination folders found.") // @todo: handle error
		c.env.CM.StopConversation(c)
		return nil
	}

	// Found folders!

	// Send the results
	var msg []string
	msg = append(msg, fmt.Sprintf("*Found %d folders:*", len(folders)))
	for i, folder := range folders {
		msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(filepath.Base(folder.Path))))
	}
	_ = c.bot.Send(user, strings.Join(msg, "\n")) // @todo: handle error

	// Send the custom reply keyboard
	var options []string
	for _, folder := range folders {
		options = append(options, fmt.Sprintf("%s", filepath.Base(folder.Path)))
	}
	options = append(options, "/cancel")
	_ = c.bot.SendKeyboardList(user, "Which folder should it download to?", options) // @todo: handle error

	return func(m interface{}) {
		// Set the selected folder
		for i, opt := range options {
			if c.bot.GetText(m) == opt {
				c.selectedFolder = &c.folderResults[i]
				break
			}
		}

		// Not a valid folder selection
		if c.selectedMovie == nil {
			_ = c.bot.Send(user, "Invalid selection.") // @todo: handle error
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddMovie(m)
	}
}

func (c *AddMovieConversation) AddMovie(m interface{}) {
	user := users.User{} // @todo: fix this, user in context
	_, err := c.env.Radarr.AddMovie(*c.selectedMovie, c.selectedQualityProfile.ID, c.selectedFolder.Path)

	// Failed to add movie
	if err != nil {
		_ = c.bot.Send(user, "Failed to add movie.") // @todo: handle error
		c.env.CM.StopConversation(c)
		return
	}

	if c.selectedMovie.PosterURL != "" {
		photo := &tb.Photo{File: tb.FromURL(c.selectedMovie.PosterURL)} // @todo: refactor here
		_ = c.bot.Send(user, photo)                                     // @todo: handle error
	}

	// Notify User
	_ = c.bot.Send(user, "Movie has been added!") // @todo: handle error

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added movie '%s'", user.DisplayName(), util.EscapeMarkdown(c.selectedMovie.String()))
	_ = c.bot.SendToAdmins(adminMsg) // @todo: refactor

	c.env.CM.StopConversation(c)
}
