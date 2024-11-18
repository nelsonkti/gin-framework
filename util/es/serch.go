package es

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	jsoniter "github.com/json-iterator/go"
)

type SearchOptions struct {
	Query map[string]interface{} // 查询条件
	Sort  []string               // 排序字段
	From  int                    // from
	Size  int                    // size
}

type SearchResult struct {
	Hits     []map[string]interface{} // 搜索结果
	Total    int                      // total （文档总数）
	Took     int                      // took（耗费多少毫秒）
	ScrollID string                   // scroll
}

func (i *Index) Search(options SearchOptions) (*SearchResult, error) {
	var buf bytes.Buffer
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(&buf).Encode(options.Query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	searchOpts := []func(request *esapi.SearchRequest){
		i.client.Search.WithContext(context.Background()),
		i.client.Search.WithIndex(i.IndexName),
		i.client.Search.WithBody(&buf),
		i.client.Search.WithTrackTotalHits(true),
	}

	if options.From != 0 {
		searchOpts = append(searchOpts, i.client.Search.WithFrom(options.From))
	}
	if options.Size != 0 {
		searchOpts = append(searchOpts, i.client.Search.WithSize(options.Size))
	}
	if len(options.Sort) > 0 {
		searchOpts = append(searchOpts, i.client.Search.WithSort(options.Sort...))
	}
	res, err := i.client.Search(searchOpts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result map[string]interface{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response body: %s", err)
	}

	hits := make([]map[string]interface{}, 0)
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var item = hit.(map[string]interface{})["_source"].(map[string]interface{})
		//追加_id进去
		if v, ok := hit.(map[string]interface{})["_id"].(string); ok {
			item["_id"] = v
		}
		hits = append(hits, item)
	}

	totalHits := int(result["hits"].(map[string]interface{})["total"].(float64))
	took := int(result["took"].(float64))
	scrollID, _ := result["_scroll_id"].(string)

	return &SearchResult{
		Hits:     hits,
		Total:    totalHits,
		Took:     took,
		ScrollID: scrollID,
	}, nil
}
