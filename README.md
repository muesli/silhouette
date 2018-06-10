# silhouette
Silhouette cluster analysis implementation in Go

## What It Does

Silhouette refers to an [algorithm](https://en.wikipedia.org/wiki/Silhouette_(clustering))
used to interpret and validate the consistency within clusters of data.

The silhouette value is a measure of how similar an object is to its own cluster
compared to other clusters. The silhouette ranges from −1 to +1, where a high
value indicates that the object is well matched to its own cluster and poorly
matched to neighboring clusters.

If most objects have a high value, then the clustering configuration is
appropriate. If many points have a low or negative value, then the clustering
configuration may have too many or too few clusters.

## When You Should Use It

- When you have numeric, multi-dimensional data sets
- If you want to check whether your data set is clustered
- When you have a vague idea of the clustering in your data set
- You want to figure out the optimal clustering configuration

## Example
```go
import (
    "github.com/muesli/silhouette"
    "github.com/muesli/clusters"
    "github.com/muesli/kmeans"
)

// initialize your data set
// for the example we'll use three distinct clusters of data points
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

// silhouette will theoretically work with multiple clustering algorithms
// it's commonly used with k-means
km := kmeans.New()

// compute the average silhouette score (coefficient) for 2 to 8 clusters, using
// the k-means clustering algorithm
scores, err := silhouette.Scores(d, 8, km)
for _, s := range scores {
    fmt.Printf("k: %d (score: %.2f)\n", s.K, s.Score)
}

// estimate the amount of clusters in our data set
// this returns the k with the highest score (where 2 <= k <= 8)
k, score, err := silhouette.EstimateK(d, 8, km)

// k is usually 3 for this example, with a score close to 1.0
// note that k-means doesn't always converge optimally
...
}
```

## Development

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/silhouette)
[![Build Status](https://travis-ci.org/muesli/silhouette.svg?branch=master)](https://travis-ci.org/muesli/silhouette)
[![Coverage Status](https://coveralls.io/repos/github/muesli/silhouette/badge.svg?branch=master)](https://coveralls.io/github/muesli/silhouette?branch=master)
[![Go ReportCard](http://goreportcard.com/badge/muesli/silhouette)](http://goreportcard.com/report/muesli/silhouette)
