package progressbar

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"time"
)

type ProgressBar struct {
	Size    int
	Current int

	Label  string
	Writer io.Writer

	percentage int
	first      bool
	done       bool

	lastUpdate time.Time
}

func NewProgressBar(size int) *ProgressBar {
	if size <= 0 {
		size = 1
	}
	return &ProgressBar{Size: size, Writer: os.Stdout}
}

func (p *ProgressBar) Update(progress int) {
	if p.done {
		panic("update called after done")
	}
	p.Current += progress
	old := p.percentage
	p.percentage = p.Current * 100 / p.Size
	if old != p.percentage && time.Since(p.lastUpdate) > 300*time.Millisecond {
		p.lastUpdate = time.Now()
		p.Draw()
	}
}

func (p *ProgressBar) Done() {
	p.done = true

	if p.percentage < 100 {
		p.percentage = 100
	}

	p.Draw()
}

func (p *ProgressBar) Draw() {
	var sb bytes.Buffer
	if !p.first {
		sb.WriteString("\n")
		p.first = true
	}
	sb.WriteString("\033[F\r")
	sb.WriteString("\r")
	sb.WriteString(p.Label)
	sb.WriteString(": ")
	sb.WriteString(strconv.Itoa(p.percentage))
	sb.WriteString("%")
	if p.done {
		sb.WriteString("\t DONE!")
	}
	sb.WriteString("\n")

	p.Writer.Write(sb.Bytes())
}
