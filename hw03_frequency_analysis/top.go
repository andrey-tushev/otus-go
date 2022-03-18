package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// Разбивка по пробельным символам.
var spacesRe = regexp.MustCompile(`\s+`)

// Сколько слов нужно получить.
const maxTopSize = 10

func Top10(text string) []string {
	// Разобьем по словам
	words := spacesRe.Split(text, -1)

	// Посчитаем каждое слово
	counts := make(map[string]int)
	for _, word := range words {
		if word == "" {
			continue
		}
		counts[word]++
	}

	// Подготовим топ популярных слов
	type topRec struct {
		word  string
		count int
	}
	// Размер заранее известен. Сделаем массив указателей,
	// чтобы при сортировке не происходило копирование большого количества данных, а только копировались указатели.
	top := make([]*topRec, len(counts))

	// Заполним топ популярных слов
	i := 0
	for word, count := range counts {
		top[i] = &topRec{
			word:  word,
			count: count,
		}
		i++
	}

	// Отсортируем топ популярных слов
	sort.Slice(top, func(i, j int) bool {
		w1, w2 := top[i], top[j]

		if w1.count == w2.count {
			return strings.Compare(w1.word, w2.word) < 0
		}
		return w1.count > w2.count
	})

	// Возьмем из топа 10 самых популярных слов
	size := len(top)
	if size > maxTopSize {
		size = maxTopSize
	}
	res := make([]string, 0, size)
	for p := 0; p < size; p++ {
		res[p] = top[p].word
	}
	return res
}
