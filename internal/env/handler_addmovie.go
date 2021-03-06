package env

import (
	"fmt"
	"strings"

	"github.com/jzenoz/gramarr/internal/radarr"
	"github.com/jzenoz/gramarr/internal/util"

	"gopkg.in/tucnak/telebot.v2"
)

// HandleAddMovie handles telegram conversation flow for adding a movie
func (e *Env) HandleAddMovie(m *telebot.Message) {
	e.CM.StartConversation(newaddMovieConversation(e), m)
}

func newaddMovieConversation(e *Env) *addMovieConversation {
	return &addMovieConversation{env: e}
}

type addMovieConversation struct {
	currentStep            func(*telebot.Message)
	movieQuery             string
	movieResults           []radarr.Movie
	folderResults          []radarr.Folder
	selectedMovie          *radarr.Movie
	selectedQualityProfile *radarr.Profile
	selectedFolder         *radarr.Folder
	env                    *Env
}

func (c *addMovieConversation) Run(m *telebot.Message) {
	c.currentStep = c.AskMovie(m)
}

func (c *addMovieConversation) Name() string {
	return "addmovie"
}

func (c *addMovieConversation) CurrentStep() func(*telebot.Message) {
	return c.currentStep
}

func (c *addMovieConversation) AskMovie(m *telebot.Message) func(*telebot.Message) {
	util.Send(c.env.Bot, m.Sender, "What movie do you want to search for?")

	return func(m *telebot.Message) {
		c.movieQuery = m.Text

		movies, err := c.env.Radarr.SearchMovies(c.movieQuery)
		c.movieResults = movies

		// Search Service Failed
		if err != nil {
			util.SendError(c.env.Bot, m.Sender, "Failed to search movies.")
			c.env.CM.StopConversation(c)
			return
		}

		// No Results
		if len(movies) == 0 {
			msg := fmt.Sprintf("No movie found with the title '%s'", util.EscapeMarkdown(c.movieQuery))
			util.Send(c.env.Bot, m.Sender, msg)
			c.env.CM.StopConversation(c)
			return
		}

		// Found some movies! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d movies:*", len(movies)))
		for i, movie := range movies {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(movie.String())))
		}
		util.Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickMovie(m)
	}
}

func (c *addMovieConversation) AskPickMovie(m *telebot.Message) func(*telebot.Message) {

	// Send custom reply keyboard
	var options []string
	for _, movie := range c.movieResults {
		options = append(options, fmt.Sprintf("%s", movie))
	}
	options = append(options, "/cancel")
	util.SendKeyboardList(c.env.Bot, m.Sender, "Which one would you like to download?", options)

	return func(m *telebot.Message) {

		// Set the selected movie
		for i, opt := range options {
			if m.Text == opt {
				c.selectedMovie = &c.movieResults[i]
				break
			}
		}

		// Not a valid movie selection
		if c.selectedMovie == nil {
			util.SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickMovie(m)
			return
		}

		c.currentStep = c.AskPickMovieQuality(m)
	}
}

func (c *addMovieConversation) AskPickMovieQuality(m *telebot.Message) func(*telebot.Message) {

	profiles, err := c.env.Radarr.GetProfile("profile")

	// GetProfile Service Failed
	if err != nil {
		util.SendError(c.env.Bot, m.Sender, "Failed to get quality profiles.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, QualityProfile := range profiles {
		options = append(options, fmt.Sprintf("%v", QualityProfile.Name))
	}
	options = append(options, "/cancel")
	util.SendKeyboardList(c.env.Bot, m.Sender, "Which quality shall I look for?", options)

	return func(m *telebot.Message) {
		// Set the selected option
		for i := range options {
			if m.Text == options[i] {
				c.selectedQualityProfile = &profiles[i]
				break
			}
		}

		// Not a valid selection
		if c.selectedQualityProfile == nil {
			util.SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickMovieQuality(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *addMovieConversation) AskFolder(m *telebot.Message) func(*telebot.Message) {

	folders, err := c.env.Radarr.GetFolders()
	c.folderResults = folders

	// GetFolders Service Failed
	if err != nil {
		util.SendError(c.env.Bot, m.Sender, "Failed to get folders.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// No Results
	if len(folders) == 0 {
		util.SendError(c.env.Bot, m.Sender, "No destination folders found.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// Found folders!

	// Send the results
	var msg []string
	msg = append(msg, fmt.Sprintf("*Found %d folders:*", len(folders)))
	for i, folder := range folders {
		msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(folder.Path)))
	}
	util.Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))

	// Send the custom reply keyboard
	var options []string
	for _, folder := range folders {
		options = append(options, fmt.Sprintf("%s", folder.Path))
	}
	options = append(options, "/cancel")
	util.SendKeyboardList(c.env.Bot, m.Sender, "Which folder should it download to?", options)

	return func(m *telebot.Message) {
		// Set the selected folder
		for i, opt := range options {
			if m.Text == opt {
				c.selectedFolder = &c.folderResults[i]
				break
			}
		}

		// Not a valid folder selection
		if c.selectedMovie == nil {
			util.SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddMovie(m)
	}
}

func (c *addMovieConversation) AddMovie(m *telebot.Message) {
	_, err := c.env.Radarr.AddMovie(*c.selectedMovie, c.selectedQualityProfile.ID, c.selectedFolder.Path)

	// Failed to add movie
	if err != nil {
		util.SendError(c.env.Bot, m.Sender, "Failed to add movie.")
		c.env.CM.StopConversation(c)
		return
	}

	if c.selectedMovie.PosterURL != "" {
		photo := &telebot.Photo{File: telebot.FromURL(c.selectedMovie.PosterURL)}
		c.env.Bot.Send(m.Sender, photo)
	}

	// Notify User
	util.Send(c.env.Bot, m.Sender, "Movie has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added movie '%s'", util.DisplayName(m.Sender), util.EscapeMarkdown(c.selectedMovie.String()))
	util.SendAdmin(c.env.Bot, c.env.Users.Admins(), adminMsg)

	c.env.CM.StopConversation(c)
}
