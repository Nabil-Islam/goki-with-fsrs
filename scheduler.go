package main

import (
    "time"
    fsrs "github.com/open-spaced-repetition/go-fsrs/v4"
)

var f = fsrs.NewFSRS(fsrs.DefaultParam())

func Review(card *Card, rating fsrs.Rating) {
    now := time.Now()
    records := f.Repeat(card.FSRSCard, now)
    card.FSRSCard = records[rating].Card
}

func IsDue(card *Card) bool {
    return !time.Now().Before(card.FSRSCard.Due)
}
