package dynconfig

import "fmt"

type Segmentation struct {
	entry *Entry
	err   error
}

func NewSegmentation(entry *Entry) Segmentation {
	if _, err := entry.AsMap(); err == nil {
		return Segmentation{entry, nil}
	} else {
		err = fmt.Errorf("failed to parse JSON configuration : %w", err)
		return Segmentation{NewEntryWithRoot(nil, "/", err), err}
	}
}

func getSegmentation(entry *Entry, lastKey string, keys ...string) *Entry {
	var result *Entry
	if len(keys) == 0 {
		result = entry.Get(lastKey)
	} else {
		result = getSegmentation(entry.Get(keys[0]), lastKey, keys[1:]...)
	}

	if result.Exists() {
		return result
	}

	defaultEntry := entry.Get(lastKey)
	if defaultEntry.Exists() {
		return defaultEntry
	}

	defaultEntry = entry.Get("*").Get(lastKey)
	if defaultEntry.Exists() {
		return defaultEntry
	}

	defaultEntry = entry.Get("default").Get(lastKey)
	if defaultEntry.Exists() {
		return defaultEntry
	}

	return entry.Get(lastKey)
}

func (s *Segmentation) Get(keys ...string) *Entry {
	if s.err != nil {
		return s.entry
	}

	if len(keys) == 0 {
		return s.entry
	}

	lastKey := keys[len(keys)-1]
	return getSegmentation(s.entry, lastKey, keys[:len(keys)-1]...)
}
