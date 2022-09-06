package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/AugustoEMoreira/learning_go_clean_architecture/entity"
	"google.golang.org/api/iterator"
)

type PostRepository interface {
	Save(p *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}

type repo struct{}

const (
	projectId      string = "golearn-21536"
	collectionName string = "posts"
)

func NewPostRepository() PostRepository {
	return &repo{}
}

func (*repo) Save(p *entity.Post) (*entity.Post, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
		return nil, err
	}
	defer client.Close()
	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"ID":    p.ID,
		"Title": p.Title,
		"Text":  p.Text,
	})
	if err != nil {
		log.Fatalf("Failed to add new post: %v", err)
		return nil, err
	}
	return p, nil
}

func (*repo) FindAll() ([]entity.Post, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
		return nil, err
	}
	defer client.Close()
	var posts []entity.Post

	iter := client.Collection(collectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate the posts: %v", err)
			return nil, err
		}
		post := entity.Post{
			ID:    doc.Data()["ID"].(int64),
			Title: doc.Data()["Title"].(string),
			Text:  doc.Data()["Text"].(string),
		}
		posts = append(posts, post)
	}
	return posts, nil
}
