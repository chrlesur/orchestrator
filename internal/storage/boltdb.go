package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"orchestrator/internal/context"
	"orchestrator/internal/job"
	"orchestrator/internal/pipeline"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	jobBucket      = "jobs"
	pipelineBucket = "pipelines"
	contextBucket  = "contexts"
)

// InitDB initialise la connexion à la base de données BoltDB
func InitDB(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Création des buckets nécessaires
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(jobBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(pipelineBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(contextBucket))
		if err != nil {
			return err
		}
		return nil
	})
}

// CloseDB ferme la connexion à la base de données
func CloseDB() {
	db.Close()
}

// SaveJob sauvegarde un job dans la base de données
func SaveJob(j *job.Job) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(jobBucket))
		buf, err := json.Marshal(j)
		if err != nil {
			return err
		}
		return b.Put([]byte(j.ID), buf)
	})
}

// GetJob récupère un job depuis la base de données
func GetJob(id string) (*job.Job, error) {
	var j job.Job
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(jobBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("job not found")
		}
		return json.Unmarshal(v, &j)
	})
	if err != nil {
		return nil, err
	}
	return &j, nil
}

// SavePipeline sauvegarde un pipeline dans la base de données
func SavePipeline(p *pipeline.Pipeline) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pipelineBucket))
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(p.ID), buf)
	})
}

// GetPipeline récupère un pipeline depuis la base de données
func GetPipeline(id string) (*pipeline.Pipeline, error) {
	var p pipeline.Pipeline
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pipelineBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("pipeline not found")
		}
		return json.Unmarshal(v, &p)
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// SaveContext sauvegarde un contexte dans la base de données
func SaveContext(c *context.Context) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(contextBucket))
		buf, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return b.Put([]byte(c.ID), buf)
	})
}

// GetContext récupère un contexte depuis la base de données
func GetContext(id string) (*context.Context, error) {
	var c context.Context
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(contextBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("context not found")
		}
		return json.Unmarshal(v, &c)
	})
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListJobs récupère tous les jobs depuis la base de données
func ListJobs() ([]*job.Job, error) {
	var jobs []*job.Job

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(jobBucket))
		return b.ForEach(func(k, v []byte) error {
			var j job.Job
			if err := json.Unmarshal(v, &j); err != nil {
				return err
			}
			jobs = append(jobs, &j)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// ListPipelines récupère tous les pipelines depuis la base de données
func ListPipelines() ([]*pipeline.Pipeline, error) {
	var pipelines []*pipeline.Pipeline

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pipelineBucket))
		return b.ForEach(func(k, v []byte) error {
			var p pipeline.Pipeline
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			pipelines = append(pipelines, &p)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return pipelines, nil
}
