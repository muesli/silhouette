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

// KScore holds the score for a value of K
type KScore struct {
	K     int
	Score float64
}

// Partitioner interface which suitable clustering algorithms should implement
type Partitioner interface {
	Partition(data clusters.Observations, k int) (clusters.Clusters, error)
}

// EstimateK estimates the amount of clusters (k) along with the silhouette
// score for that value, using the given partitioning algorithm
func EstimateK(data clusters.Observations, kmax int, m Partitioner) (int, float64, error) {
	scores, err := Scores(data, kmax, m)
	if err != nil {
		return 0, -1.0, err
	}

	r := KScore{
		K: -1,
	}
	for _, score := range scores {
		if r.K < 0 || score.Score > r.Score {
			r = score
		}
	}

	return r.K, r.Score, nil
}

// Scores calculates the silhouette scores for all values of k between 2 and
// kmax, using the given partitioning algorithm
func Scores(data clusters.Observations, kmax int, m Partitioner) ([]KScore, error) {
	var r []KScore

	for k := 2; k <= kmax; k++ {
		s, err := Score(data, k, m)
		if err != nil {
			return r, err
		}

		r = append(r, KScore{
			K:     k,
			Score: s,
		})
	}

	return r, nil
}

// Score calculates the silhouette score for a given value of k, using the given
// partitioning algorithm
func Score(data clusters.Observations, k int, m Partitioner) (float64, error) {
	cc, err := m.Partition(data, k)
	if err != nil {
		return -1.0, err
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

	return si / float64(sc), nil
}

// Plot creates a graph of the silhouette scores
func Plot(filename string, scores []KScore) error {
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
		return err
	}
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}
