package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(input string) []string {
	input = strings.ToLower(input)

	re := regexp.MustCompile("([-A-Za-zА-Яа-я]{2,}|[A-Za-zА-Яа-я])")
	words := re.FindAllString(input, -1)

	freq := CalculateFrequency(words)

	sort.Slice(freq, freq.Less)

	return freq.Top10StringArray()
}

type WordFrequency struct {
	Word      string
	Frequency int
}

type ByFrequency []WordFrequency

func (wf ByFrequency) Len() int {
	return len(wf)
}

func (wf ByFrequency) Swap(i, j int) {
	wf[i], wf[j] = wf[j], wf[i]
}

func (wf ByFrequency) Less(i, j int) bool {
	less := wf[i].Frequency > wf[j].Frequency
	if wf[i].Frequency == wf[j].Frequency {
		less = wf[i].Word < wf[j].Word
	}
	return less
}

func (wf ByFrequency) Top10StringArray() []string {
	var top = make([]string, 0, 10)
	maxReturnedCount := 10
	if wf.Len() < maxReturnedCount {
		maxReturnedCount = wf.Len()
	}
	for i := 0; i < maxReturnedCount; i++ {
		top = append(top, wf[i].Word)
	}
	return top
}

func CalculateFrequency(words []string) ByFrequency {
	var frequency = make(map[string]int)
	for _, word := range words {
		_, ok := frequency[word]
		if ok {
			frequency[word]++
		} else {
			frequency[word] = 1
		}
	}
	var byFreq = make(ByFrequency, 0, len(frequency))
	for word, freq := range frequency {
		byFreq = append(byFreq, WordFrequency{
			Word:      word,
			Frequency: freq,
		})
	}
	return byFreq
}
