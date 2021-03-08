package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Wave struct {
	SampleRate beep.SampleRate
	phase      int
	freq       float64
	drain      bool
}

var enveloppe = map[int]float64{
	1: 0.8,
	2: 0.2,
	3: 0.1,
	4: 0.1,
	5: 0.1,
	6: 0.1,
	7: 0.1,
	8: 0.05,
	9: 0.01,
}

func (w *Wave) Stream(samples [][2]float64) (n int, ok bool) {
	if w.drain {
		return 0, false
	}
	period := w.SampleRate.N(time.Second / time.Duration(w.freq))
	var j int
	for i := range samples {
		j = i + w.phase
		t := float64(j) / float64(period)
		vol := float64(0)
		for h, v := range enveloppe {
			vol += math.Sin(t*math.Pi*2*float64(h)) * v
		}
		samples[i][0] = vol
		samples[i][1] = vol
	}
	w.phase = (j + 1) % period
	return len(samples), true
}

func (w *Wave) Err() error {
	return nil
}

type BufferedStreamer struct {
	SampleRate beep.SampleRate
	samples    Queue
	drain      bool
}

func (w *BufferedStreamer) Add(left float64, right float64) {
	pair := [2]float64{left, right}
	w.samples.Push(pair)
}

func (w *BufferedStreamer) Err() error {
	return nil
}

func (w *BufferedStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if w.drain {
		return 0, false
	}
	for i := range samples {
		front := w.samples.Pop()
		if front == nil {
			return i, true
		}
		pair, ok := front.([2]float64)
		if ok {
			//fmt.Println(i, pair)
			samples[i][0] = pair[0]
			samples[i][1] = pair[1]
		} else {
			panic(fmt.Sprintf("Should not happen : FIFO contains %T", front))
		}
	}
	//fmt.Println("Read", len(samples), "samples")
	return len(samples), true
}

func play(quit chan bool) {
	defer speaker.Close()
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	w := &beep.Mixer{}
	var b float64 = 110.0
	w.Add(&Wave{SampleRate: sr, freq: b},
		&Wave{SampleRate: sr, freq: b * 3 / 2},
		&Wave{SampleRate: sr, freq: b * 10 / 4})
	speaker.Play(w)
	select {
	//case <-quit:
	//	w.drain = true
	//	return
	}
}

func playStream(s beep.Streamer, sr beep.SampleRate) {
	fmt.Println("Start speaker")
	defer speaker.Close()
	defer fmt.Println("Close speaker")
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(s)
	select {}
}
