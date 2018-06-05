// Package silhouette implements the silhouette cluster analysis algorithm
// See: https://en.wikipedia.org/wiki/Silhouette_(clustering)
package silhouette

import (
	"bytes"
	"io/ioutil"
	"math"

	"github.com/muesli/clusters"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type Score struct {
	K     int
	Score float64
}

type Partitioner interface {
	Partition(data clusters.Observations, k int) (clusters.Clusters, error)
}

func EstimateK(data clusters.Observations, m Partitioner) (Score, error) {
	scores, err := Scores(data, m)
	if err != nil {
		return Score{}, err
	}

	r := Score{
		K: -1,
	}
	for _, score := range scores {
		if r.K < 0 || score.Score > r.Score {
			r = score
		}
	}

	return r, nil
}

func Scores(data clusters.Observations, m Partitioner) ([]Score, error) {
	var r []Score

	for n := 2; n < 10; n++ {
		cc, err := m.Partition(data, n)
		if err != nil {
			return []Score{Score{K: 0, Score: -1.0}}, err
		}

		var si float64
		var sc int64
		for ci, c := range cc {
			for _, p := range c.Observations {
				ai := clusters.AverageDistance(p, c.Observations)
				_, bi := cc.Neighbour(p, ci)

				si += (bi - ai) / math.Max(ai, bi)
				sc++
			}
		}

		sd := si / float64(sc)
		r = append(r, Score{
			K:     n,
			Score: sd,
		})
	}

	// fmt.Printf("%+v\n", r)
	return r, nil
}

func Plot(scores []Score) {
	var series []chart.Series

	for _, s := range scores {
		series = append(series, chart.ContinuousSeries{
			Style: chart.Style{
				Show:        true,
				StrokeWidth: chart.Disabled,
				DotColor:    drawing.ColorBlue,
				DotWidth:    8,
			},
			XValues: []float64{float64(s.K)},
			YValues: []float64{s.Score},
		})
	}

	graph := chart.Chart{
		Height: 1024,
		Width:  1024,
		Series: series,
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("silhouette.png", buffer.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
