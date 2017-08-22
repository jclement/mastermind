// Package mastermind provides a game object for the mastermind game
package mastermind

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// GameRow represents a row in the game
type GameRow []int

// GameState represents the state of a game
type GameState byte

func (state GameState) String() string {
	switch state {
	case StateActive:
		return "ACTIVE"
	case StateWin:
		return "WON"
	case StateLost:
		return "LOST"
	default:
		panic("Unsupported state")
	}
}

// ErrGameNotActive is raised if the player tries to make a move on an inactive game
var ErrGameNotActive = errors.New("Game not active")

// ErrBadGuess is raised if the guess provided isn't valid (number of pegs or colors)
var ErrBadGuess = errors.New("Bad Guess")

const (
	// StateActive represents an in-progress game
	StateActive GameState = iota
	// StateWin represents a complete game that was successfully completed
	StateWin
	// StateLost represents a complete game that was not successfully completed
	StateLost
)

// Game represents state of a mastermind game
type Game struct {
	Solution   GameRow
	Guesses    []GameRow
	State      GameState
	maxColors  int
	maxGuesses int
}

// GuessResult represents the result of a guess
type GuessResult struct {
	NumCorrectPins int
	NumPartialPins int
	State          GameState
}

// SubmitGuess places a guess against an active game and returns the result
func (game *Game) SubmitGuess(g GameRow) (GuessResult, error) {
	if game.State != StateActive {
		return GuessResult{}, ErrGameNotActive
	}

	if len(g) != len(game.Solution) {
		return GuessResult{}, ErrBadGuess
	}

	for i := 0; i < len(g); i++ {
		if g[i] >= game.maxColors {
			return GuessResult{}, ErrBadGuess
		}
	}

	game.Guesses = append(game.Guesses, g)

	result := GuessResult{}

	used := make([]bool, len(game.Solution))

	for i := 0; i < len(game.Solution); i++ {
		if g[i] == game.Solution[i] {
			result.NumCorrectPins++
			used[i] = true
			continue
		} else {
			for j := 0; j < len(game.Solution); j++ {
				if g[i] == game.Solution[j] && !used[j] {
					result.NumPartialPins++
					used[j] = true
					continue
				}
			}
		}

	}

	if result.NumCorrectPins == len(game.Solution) {
		game.State = StateWin
	}

	if len(game.Guesses) >= game.maxGuesses {
		game.State = StateLost
	}

	// finally make sure the returned result reflects the updated state
	result.State = game.State

	return result, nil

}

// NewGame creates a new game object
func NewGame() *Game {

	maxColors := 4
	holeCount := 4
	allowDuplicates := false
	solution := make([]int, holeCount)

	usedColours := make([]bool, maxColors)

	for i := 0; i < holeCount; i++ {
		color := 0
		for {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(maxColors)))
			color = int(n.Int64())
			if allowDuplicates || !usedColours[color] {
				usedColours[color] = true
				break
			}
		}
		solution[i] = color
	}

	return &Game{
		Solution:   solution,
		Guesses:    make([]GameRow, 0),
		State:      StateActive,
		maxColors:  maxColors,
		maxGuesses: 8,
	}
}
