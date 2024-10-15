package main

import (
	"context"
	"errors"
	"reflect"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm/schema"
)

type YamlMap map[string]interface{}

// Implement the sql.Scanner interface
// This is called on retrieval. The YAML data in the DB is converted to map
func (ym *YamlMap) Scan(ctx context.Context, field *schema.Field,
	dst reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		*ym = nil
		return nil
	}
	bytes, ok := dbValue.([]byte)
	if !ok {
		return errors.New("failed to unmarshal YAML value: source data not bytes")
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(bytes, &m); err != nil {
		return err
	}

	*ym = m
	return nil
}

// Implement the Valuer interface
// This method converts the map into YAML to store it in the DB.
func (ym YamlMap) Value(ctx context.Context, field *schema.Field,
	dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if ym == nil {
		return nil, nil
	}
	return yaml.Marshal(ym)
}
