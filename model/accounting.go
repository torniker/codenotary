package model

import (
	"encoding/json"

	"github.com/torniker/codenotary/immudb"
)

type Accounting struct {
	ID      string `json:"id,omitempty"`
	Number  string `json:"account_number"`
	Name    string `json:"account_name"`
	IBAN    string `json:"iban"`
	Address string `json:"address"`
	Amount  int    `json:"amount"`
	Type    string `json:"type"`
}

func FromSearchResult(result immudb.SearchResult) (accounting []Accounting, err error) {
	for _, revision := range result.Revisions {
		var a Accounting
		err = json.Unmarshal(revision.Document, &a)
		if err != nil {
			return
		}
		accounting = append(accounting, a)
	}
	return
}
