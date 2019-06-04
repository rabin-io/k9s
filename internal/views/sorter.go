package views

import (
	"regexp"
	"strconv"
	"time"

	"github.com/derailed/k9s/internal/resource"
	res "k8s.io/apimachinery/pkg/api/resource"
	"vbom.ml/util/sortorder"
)

type rowSorter struct {
	rows  resource.Rows
	index int
	asc   bool
}

func (s rowSorter) Len() int {
	return len(s.rows)
}

func (s rowSorter) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
}

func (s rowSorter) Less(i, j int) bool {
	return less(s.asc, s.rows[i][s.index], s.rows[j][s.index])
}

// ----------------------------------------------------------------------------

type groupSorter struct {
	rows []string
	asc  bool
}

func (s groupSorter) Len() int {
	return len(s.rows)
}

func (s groupSorter) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
}

func (s groupSorter) Less(i, j int) bool {
	return less(s.asc, s.rows[i], s.rows[j])
}

// ----------------------------------------------------------------------------
// Helpers...

func less(asc bool, c1, c2 string) bool {
	if o, ok := isMetricSort(asc, c1, c2); ok {
		return o
	}

	if o, ok := isDurationSort(asc, c1, c2); ok {
		return o
	}

	if o, ok := isIntegerSort(asc, c1, c2); ok {
		return o
	}

	b := sortorder.NaturalLess(c1, c2)
	if asc {
		return b
	}
	return !b
}

func isDurationSort(asc bool, s1, s2 string) (bool, bool) {
	d1, ok1 := isDuration(s1)
	d2, ok2 := isDuration(s2)
	if !ok1 || !ok2 {
		return false, false
	}

	if asc {
		return d1 <= d2, true
	}
	return d1 >= d2, true
}

func isMetricSort(asc bool, c1, c2 string) (bool, bool) {
	q1, err1 := res.ParseQuantity(c1)
	q2, err2 := res.ParseQuantity(c2)
	if err1 != nil || err2 != nil {
		return false, false
	}

	if asc {
		return q1.Cmp(q2) <= 0, true
	}
	return q1.Cmp(q2) > 0, true
}

func isIntegerSort(asc bool, c1, c2 string) (bool, bool) {
	n1, err := strconv.Atoi(c1)
	if err != nil {
		return false, false
	}
	n2, _ := strconv.Atoi(c2)
	if asc {
		return n1 <= n2, true
	}
	return n1 > n2, true
}

var metricRX = regexp.MustCompile(`\A(\d+)(m|Mi)\z`)

func isMetric(s string) (string, bool) {
	if m := metricRX.FindStringSubmatch(s); len(m) == 3 {
		return m[1], true
	}
	return s, false
}

func isDuration(s string) (time.Duration, bool) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return d, false
	}
	return d, true
}
