package main

import (
	"github.com/GabiBizdoc/golang-playground/pkg/progressbar"
	"time"
)

func main() {
	pg := progressbar.NewProgressBar(100)
	pg.Label = "Progress"
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second / 2)
		pg.Update(3)
	}
	pg.Done()
}
