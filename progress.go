package main

import (
	"fmt"
	"strings"
	"time"
)

type progressBar struct {
	total   int64
	current int64
	width   int
	start   time.Time
}

func newProgressBar(total int64) *progressBar {
	return &progressBar{
		total: total,
		width: 40,
		start: time.Now(),
	}
}

func (p *progressBar) add(n int64) {
	p.current += n
	p.render()
}

func (p *progressBar) render() {
	percent := float64(p.current) / float64(p.total) * 100
	filled := int(float64(p.width) * float64(p.current) / float64(p.total))

	elapsed := time.Since(p.start).Seconds()
	speed := float64(p.current) / elapsed
	remaining := time.Duration((float64(p.total-p.current) / speed)) * time.Second

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)

	fmt.Printf("\r[%s] %.1f%% | %s/%s | %s/s | ETA: %s",
		bar,
		percent,
		formatBytes(p.current),
		formatBytes(p.total),
		formatBytes(int64(speed)),
		formatDuration(remaining),
	)
}

func (p *progressBar) finish() {
	p.current = p.total
	p.render()
	fmt.Println()
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
