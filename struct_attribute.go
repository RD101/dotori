package main

import "encoding/json"

type JsonAttributes []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (u *JsonAttributes) ToAttributes(jsonAttributeString string) (map[string]string, error) {
	attr := make(map[string]string)
	jsonAttributes := JsonAttributes{}
	err := json.Unmarshal([]byte(jsonAttributeString), &jsonAttributes)
	if err != nil {
		return attr, err
	}
	for _, attribute := range jsonAttributes {
		if attribute.Key == "" || attribute.Value == "" {
			continue
		}
		attr[attribute.Key] = attribute.Value
	}
	return attr, nil
}
