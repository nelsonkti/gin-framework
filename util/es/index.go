package es

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	jsoniter "github.com/json-iterator/go"
)

type Index struct {
	client    *elasticsearch.Client
	IndexName string
}

func NewIndex(client *Client, indexName string) *Index {
	return &Index{
		client:    client.Client(),
		IndexName: indexName,
	}
}

func (i *Index) Create() error {
	res, err := i.client.Indices.Create(i.IndexName)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}
	return nil
}

func (i *Index) Delete() error {
	res, err := i.client.Indices.Delete([]string{i.IndexName}, i.client.Indices.Delete.WithContext(context.Background()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting index: %s", res.String())
	}
	return nil
}

func (i *Index) AddDocument(docID string, doc interface{}) error {
	var buf bytes.Buffer
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("error encoding document: %s", err)
	}

	res, err := i.client.Index(
		i.IndexName,
		&buf,
		i.client.Index.WithDocumentID(docID),
		i.client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error adding document: %s", res.String())
	}
	return nil
}

func (i *Index) GetDocument(docID string) (map[string]interface{}, error) {
	res, err := i.client.Get(i.IndexName, docID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	var doc map[string]interface{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(res.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("error parsing document: %s", err)
	}
	return doc, nil
}

func (i *Index) DeleteDocument(docID string) error {
	res, err := i.client.Delete(i.IndexName, docID)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}
	return nil
}
