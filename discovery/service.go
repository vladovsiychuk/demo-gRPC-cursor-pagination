package discovery

import (
	"context"
	"sync"

	"github.com/vladovsiychuk/demo-grpc/util"
)

type Service interface {
	AddPost(context.Context, *AddPost) (*Post, error)
	GetPosts(ctx context.Context, cursor string) (Page[*Post], error)
	Close()
}

type service struct {
	posts    []*Post
	idToIdx  map[string]int
	pageSize uint
	mu       sync.RWMutex
}

func NewService(pageSize uint) Service {
	return &service{
		posts:    make([]*Post, 0, 10),
		idToIdx:  make(map[string]int),
		pageSize: pageSize,
	}
}

func (s *service) AddPost(ctx context.Context, p *AddPost) (*Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := util.MustGetNewID()
	post := &Post{
		ID:          id,
		OwnerID:     p.OwnerID,
		FrontPicUrl: p.FrontPicUrl,
		BackPicUrl:  p.BackPicUrl,
	}
	s.posts = append(s.posts, post)
	s.idToIdx[id] = len(s.posts) - 1
	return post, nil
}

func (s *service) GetPosts(ctx context.Context, cursor string) (Page[*Post], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.posts)
	if total == 0 {
		return Page[*Post]{}, nil
	}

	end := total
	if cursor != "" {
		if idx, ok := s.idToIdx[cursor]; ok {
			end = idx
		} else {
			return Page[*Post]{}, nil
		}
	}

	start := end - int(s.pageSize)
	if start < 0 {
		start = 0
	}
	if end <= start {
		return Page[*Post]{}, nil
	}

	page := make([]*Post, end-start)
	for i := range page {
		page[i] = s.posts[end-1-i]
	}

	nextCursor := ""
	if len(page) > 0 {
		nextCursor = page[len(page)-1].ID
	}

	return Page[*Post]{Data: page, Cursor: nextCursor}, nil
}

func (s *service) Close() {}
