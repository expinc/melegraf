package metric

import (
	"fmt"
	"time"
)

type Tag struct {
	Key   string
	Value string
}

type Field struct {
	Key   string
	Value interface{}
}

type Metric struct {
	Name   string
	Tags   []Tag
	Fields []Field
	Time   time.Time
}

func (mt *Metric) Copy() Metric {
	tags := make([]Tag, len(mt.Tags))
	copy(tags, mt.Tags)

	fields := make([]Field, len(mt.Fields))
	copy(fields, mt.Fields)

	return Metric{
		Name:   mt.Name,
		Tags:   tags,
		Fields: fields,
		Time:   mt.Time,
	}
}

func (mt *Metric) AddTag(key, value string) error {
	for _, tag := range mt.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag with key '%s' already exists", key)
		}
	}
	mt.Tags = append(mt.Tags, Tag{Key: key, Value: value})
	return nil
}

func (mt *Metric) AddField(key string, value interface{}) error {
	for _, field := range mt.Fields {
		if field.Key == key {
			return fmt.Errorf("field with key '%s' already exists", key)
		}
	}
	mt.Fields = append(mt.Fields, Field{Key: key, Value: value})
	return nil
}

func (mt *Metric) GetTag(key string) (string, error) {
	for _, tag := range mt.Tags {
		if tag.Key == key {
			return tag.Value, nil
		}
	}
	return "", fmt.Errorf("tag with key '%s' not found", key)
}

func (mt *Metric) GetField(key string) (interface{}, error) {
	for _, field := range mt.Fields {
		if field.Key == key {
			return field.Value, nil
		}
	}
	return nil, fmt.Errorf("field with key '%s' not found", key)
}
