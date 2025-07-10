package url

import "fmt"

type InMemoryRepository struct {
	urls map[int]*Url
}

func (r *InMemoryRepository) GetById(id int) (*Url, error) {
	return r.urls[id], nil
}

func (r *InMemoryRepository) GetByValue(shortened string) (*Url, error) {
	for _, url := range r.urls {
		if url.Shortened == shortened {
			return url, nil
		}
	}
	return nil, fmt.Errorf("could not find url with shortened value %s", shortened)
}

func (r *InMemoryRepository) Insert(item *Url) (*Url, error) {
	if item.Id == -1 {
		item.Id = len(r.urls)
	}
	r.urls[item.Id] = item
	return item, nil
}

func (r *InMemoryRepository) Update(item *Url) error {
	if _, exists := r.urls[item.Id]; !exists {
		return fmt.Errorf("url with id %d not found", item.Id)
	}
	r.urls[item.Id].Visits += 1
	return nil
}

func (r *InMemoryRepository) Next() (int, error) {
	return len(r.urls), nil
}

func NewRepository() *InMemoryRepository {
	return &InMemoryRepository{
		urls: make(map[int]*Url),
	}
}
