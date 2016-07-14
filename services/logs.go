package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/mateuszdyminski/logag/libcfg"
	"github.com/mateuszdyminski/logag/model"
	"github.com/mateuszdyminski/logag/ws"
	"gopkg.in/olivere/elastic.v3"
)

type LogService struct {
	cfg *libcfg.Cfg
	esc *elastic.Client
	db  *bolt.DB
	Ws  *ws.Hub
}

func NewLogService(cfg *libcfg.Cfg, esc *elastic.Client, ws *ws.Hub, db *bolt.DB) *LogService {
	return &LogService{cfg: cfg, esc: esc, db: db, Ws: ws}
}

func (s *LogService) AddLog(user string, logs []model.Log) error {
	// store logs in elasticsearch
	var enqued int
	bulkRequest := s.esc.Bulk()
	for i, log := range logs {
		logs[i].User = user

		if enqued > 0 && enqued%s.cfg.BatchSize == 0 {
			if _, err := bulkRequest.Do(); err != nil {
				return fmt.Errorf("Can't execute bulk. Err: %v", err)
			}

			logrus.Infof("Bulk with %v logs indexed!", s.cfg.BatchSize)

			bulkRequest = s.esc.Bulk()
		}

		bulkRequest.Add(
			elastic.NewBulkIndexRequest().
				Index("logs").
				Type("log").
				Id(fmt.Sprintf("%s%d", user, log.Time.Unix())).
				Doc(logs[i]))

		enqued++
	}

	// send batch if it's not sent already
	if bulkRequest.NumberOfActions() > 0 {
		toDo := bulkRequest.NumberOfActions()
		if _, err := bulkRequest.Do(); err != nil {
			return fmt.Errorf("Can't execute bulk. Err: %v", err)
		}

		logrus.Infof("Bulk with %v logs indexed!", toDo)
	}

	// send message to all WS clients
	for i := range logs {
		s.Ws.Broadcast <- &logs[i]
	}

	return nil
}

func (s *LogService) Search(query, level string, from, to time.Time, size, skip int) (*model.Response, error) {
	var elasticQuery elastic.Query
	emptyTime := time.Time{}

	if query == "" {
		if level == "" {
			if from == emptyTime {
				if to == emptyTime {
					elasticQuery = elastic.NewMatchAllQuery()
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			} else {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			}
		} else {
			if from == emptyTime {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			} else {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			}
		}
	} else {
		if level == "" {
			if from == emptyTime {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			} else {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			}
		} else {
			if from == emptyTime {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			} else {
				if to == emptyTime {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewMatchAllQuery())
				} else {
					elasticQuery = elastic.NewBoolQuery().
						Must(elastic.NewMatchQuery("msg", query)).
						Must(elastic.NewMatchQuery("level", level)).
						Must(elastic.NewRangeQuery("time").Gte(from)).
						Must(elastic.NewRangeQuery("time").Lte(to)).
						Must(elastic.NewMatchAllQuery())
				}
			}
		}
	}

	searchResult, err := s.esc.Search().
		Index("logs").
		Type("log").
		Query(elasticQuery).
		From(skip).Size(size).
		Do()
	if err != nil {
		return nil, err
	}

	response := model.Response{}
	logs := make([]model.Log, 0)
	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits != nil {
		logrus.Infof("For (query: %s, level: %s, from: %v, to: %v, size: %d, skip: %d) Found a total of %d logs", query, level, from, to, size, skip, searchResult.Hits.TotalHits)

		response.Total = searchResult.Hits.TotalHits

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			var l model.Log
			err := json.Unmarshal(*hit.Source, &l)
			if err != nil {
				return nil, err
			}

			l.Score = hit.Score

			logs = append(logs, l)
		}
	}
	response.Data = logs

	return &response, nil
}

func (s *LogService) CreateIndex() error {
	exists, err := s.esc.IndexExists("logs").Do()
	if err != nil {
		return err
	}

	if !exists {
		// Create an index if not exists
		if _, err = s.esc.
			CreateIndex("logs").
			BodyString(model.LogMapping).
			Do(); err != nil {
			return err
		}

		logrus.Infof("Index logs created!")
	} else {
		logrus.Infof("Index logs already exists!")
	}

	return nil
}

func (s *LogService) RegisterFilter(f model.Filter) error {
	if len(f.Keywords) == 0 && f.Level == "" {
		return fmt.Errorf("Filter should contain at least one filter options {keywords,level}")
	}

	s.Ws.RegisterFilter <- &f

	return nil
}

func (s *LogService) UnregisterFilter(id string) error {
	s.Ws.UnregisterFilter <- id

	return nil
}
