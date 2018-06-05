package silhouette

import (
	"math/rand"
	"testing"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

const (
	randomSeed = 42
)

func TestScores(t *testing.T) {
	rand.Seed(randomSeed)

	var d clusters.Observations
	for x := 0; x < 64; x++ {
		d = append(d, clusters.Coordinates{
			rand.Float64() * 0.1,
			rand.Float64() * 0.1,
		})
	}
	for x := 0; x < 64; x++ {
		d = append(d, clusters.Coordinates{
			0.5 + rand.Float64()*0.1,
			0.5 + rand.Float64()*0.1,
		})
	}
	for x := 0; x < 64; x++ {
		d = append(d, clusters.Coordinates{
			0.9 + rand.Float64()*0.1,
			0.9 + rand.Float64()*0.1,
		})
	}

	km := kmeans.New()
	scores, err := Scores(d, km)
	if err != nil {
		t.Error(err)
		return
	}
	Plot(scores)

	km = kmeans.New()
	estimate, err := EstimateK(d, km)
	if estimate.K != 3 {
		t.Errorf("Expected k-value of 3, got %d", estimate.K)
	}
	// fmt.Printf("Best estimate: %+v %v\n", estimate, err)
}
