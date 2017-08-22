package mastermind

import (
	"testing"
)

var tests = []struct {
	solution            GameRow
	guess               GameRow
	expectedError       bool
	expectedCorrectPins int
	expectedPartialPins int
	expectedState       GameState
}{
	{GameRow{1, 2}, GameRow{3, 0}, false, 0, 0, StateActive},
	{GameRow{1, 2}, GameRow{1, 0}, false, 1, 0, StateActive},
	{GameRow{1, 2}, GameRow{2, 1}, false, 0, 2, StateActive},
	{GameRow{1, 2}, GameRow{1, 2}, false, 2, 0, StateWin},
	{GameRow{0, 1, 2, 3}, GameRow{0, 1, 2, 3}, false, 4, 0, StateWin},
	{GameRow{0, 0, 0, 0}, GameRow{0, 1, 2, 3}, false, 1, 0, StateActive},
	{GameRow{0, 1, 2, 3}, GameRow{0, 0, 0, 0}, false, 1, 0, StateActive},
	{GameRow{0, 0, 2, 3}, GameRow{1, 0, 0, 2}, false, 1, 2, StateActive},
	{GameRow{1, 2}, GameRow{1, 2, 3}, true, 0, 0, StateActive},
	{GameRow{1, 2}, GameRow{1, 4}, true, 0, 0, StateActive},
}

func TestSimpleCases(t *testing.T) {
	for testIndex, test := range tests {
		game := Game{
			Solution:   test.solution,
			Guesses:    make([]GameRow, 0),
			State:      StateActive,
			maxColors:  4,
			maxGuesses: 8,
		}
		res, err := game.SubmitGuess(test.guess)
		if err == nil && test.expectedError {
			t.Fatalf("Test %v - Expected an error but didn't get one", testIndex)
		}
		if err != nil && !test.expectedError {
			t.Fatalf("Test %v - Didn't expect an error but got one", testIndex)
		}
		if res.NumCorrectPins != test.expectedCorrectPins {
			t.Fatalf("Test %v - Expected %v correct pins but got %v", testIndex, test.expectedCorrectPins, res.NumCorrectPins)
		}
		if res.NumPartialPins != test.expectedPartialPins {
			t.Fatalf("Test %v - Expected %v partial pins but got %v", testIndex, test.expectedPartialPins, res.NumPartialPins)
		}
		if res.State != test.expectedState {
			t.Fatalf("Test %v - Expected state %v but got %v", testIndex, test.expectedState, res.State)
		}
	}
}

func TestLose(t *testing.T) {
	game := Game{
		Solution:   GameRow{0, 0},
		Guesses:    make([]GameRow, 0),
		State:      StateActive,
		maxColors:  4,
		maxGuesses: 2,
	}

	res, _ := game.SubmitGuess(GameRow{1, 1})
	if res.State != StateActive {
		t.Fatal("Expected Active")
	}

	res, _ = game.SubmitGuess(GameRow{1, 1})
	if res.State != StateLost {
		t.Fatal("Expected Lost")
	}
}

func TestGuessClosed(t *testing.T) {
	for _, state := range []GameState{StateLost, StateWin} {

		game := Game{
			Solution:   GameRow{0, 0},
			Guesses:    make([]GameRow, 0),
			State:      state,
			maxColors:  4,
			maxGuesses: 2,
		}

		_, err := game.SubmitGuess(GameRow{1, 1})
		if err != ErrGameNotActive {
			t.Fatal("Expected Error")
		}
	}
}
