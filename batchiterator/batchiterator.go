package batchiterator

func NewIterator[t any](slice []t, windowSize int) func() ([]t, bool) {
	index := 0
	length := len(slice)
	return func() ([]t, bool) {
		if index >= length {
			return nil, false
		}
		end := index + windowSize
		if end > length {
			end = length
		}
		window := slice[index:end]
		index = end
		return window, true
	}
}
