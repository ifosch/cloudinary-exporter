package heracles

import (
	// "sort"
	// "strconv"
	// "strings"

	"github.com/prometheus/client_golang/prometheus"
)

type GaugeMetrics map[int]*prometheus.GaugeVec

func NewGaugeMetrics(
	namespace, name, docString string,
	labels prometheus.Labels,
) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        name,
			Help:        docString,
			ConstLabels: labels,
		},
		[]string{},
	)
}

// func (m GaugeMetrics) String() string {
// 	keys := make([]int, 0, len(m))
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	sort.Ints(keys)
// 	s := make([]string, len(keys))
// 	for i, k := range keys {
// 		s[i] = strconv.Itoa(k)
// 	}
// 	return strings.Join(s, ",")
// }
