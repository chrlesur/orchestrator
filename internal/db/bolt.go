package db

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/boltdb/bolt"
    "github.com/chrlesur/orchestrator/internal/models"
)

var jobBucket = []byte("jobs")
var pipelineBucket = []byte("pipelines")

type Store struct {
	db *bolt.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(jobBucket)
		if err != nil {
			return fmt.Errorf("could not create jobs bucket: %v", err)
		}
		_, err = tx.CreateBucketIfNotExists(pipelineBucket)
		if err != nil {
			return fmt.Errorf("could not create pipelines bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	return &Store{
		db: db,
	}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) SaveJob(j *models.Job) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(jobBucket)
        encoded, err := json.Marshal(j)
        if err != nil {
            return fmt.Errorf("could not encode job %s: %v", j.ID, err)
        }
        return b.Put([]byte(j.ID), encoded)
    })
}

func (s *Store) GetJob(id string) (*models.Job, error) {
    var j models.Job
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(jobBucket)
        v := b.Get([]byte(id))
        if v == nil {
            return fmt.Errorf("job %s not found", id)
        }
        return json.Unmarshal(v, &j)
    })
    if err != nil {
        return nil, err
    }
    return &j, nil
}

func (s *Store) SavePipeline(p *models.Pipeline) error {
    return s.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(pipelineBucket)
        encoded, err := json.Marshal(p)
        if err != nil {
            return fmt.Errorf("could not encode pipeline %s: %v", p.ID, err)
        }
        return b.Put([]byte(p.ID), encoded)
    })
}

func (s *Store) GetPipeline(id string) (*models.Pipeline, error) {
    var p models.Pipeline
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(pipelineBucket)
        v := b.Get([]byte(id))
        if v == nil {
            return fmt.Errorf("pipeline %s not found", id)
        }
        return json.Unmarshal(v, &p)
    })
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func (s *Store) GetAllJobs() ([]*models.Job, error) {
    var jobs []*models.Job
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(jobBucket)
        return b.ForEach(func(k, v []byte) error {
            var j models.Job
            if err := json.Unmarshal(v, &j); err != nil {
                return err
            }
            jobs = append(jobs, &j)
            return nil
        })
    })
    if err != nil {
        return nil, fmt.Errorf("could not get jobs: %v", err)
    }
    return jobs, nil
}

func (s *Store) GetAllPipelines() ([]*models.Pipeline, error) {
    var pipelines []*models.Pipeline
    err := s.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(pipelineBucket)
        return b.ForEach(func(k, v []byte) error {
            var p models.Pipeline
            if err := json.Unmarshal(v, &p); err != nil {
                return err
            }
            pipelines = append(pipelines, &p)
            return nil
        })
    })
    if err != nil {
        return nil, fmt.Errorf("could not get pipelines: %v", err)
    }
    return pipelines, nil
}