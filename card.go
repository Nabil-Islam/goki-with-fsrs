package main

import fsrs "github.com/open-spaced-repetition/go-fsrs/v4"

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
  FSRSCard fsrs.Card `json:"fsrsCard"`
}

func (c Card) FilterValue() string { return c.Front }
func (c Card) Title() string       { return c.Front }
func (c Card) Description() string { return c.Back }

func NewCard(front, back string) *Card {
	return &Card{
		Front:    front,
		Back:     back,
		FSRSCard: fsrs.NewCard(),
	}
}

func (c *Card) EnsureFSRS() {
    // if FSRSCard was never saved, it'll be all zero values
    if c.FSRSCard.Reps == 0 && c.FSRSCard.Stability == 0 {
        c.FSRSCard = fsrs.NewCard()
    }
}
