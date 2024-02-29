package immudb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OperatorEqual          Operator = "EQ"
	OperatorNotEqual       Operator = "NE"
	OperatorLessThan       Operator = "LT"
	OperatorLessOrEqual    Operator = "LE"
	OperatorGreaterThan    Operator = "GT"
	OperatorGreaterOrEqual Operator = "GE"
	OperatorLike           Operator = "LIKE"
)

type (
	Client struct {
		baseURL string
		headers map[string]string
		ledger  string
	}

	Query struct {
		Expressions []Expression `json:"expressions,omitempty"`
		Limit       int          `json:"limit,omitempty"`
	}
	Expression struct {
		FieldComparisons []FieldComparison `json:"fieldComparisons,omitempty"`
		OrderBy          []OrderBy         `json:"orderBy,omitempty"`
	}
	Operator        string
	FieldComparison struct {
		Field    string      `json:"field,omitempty"`
		Operator Operator    `json:"operator,omitempty"`
		Value    interface{} `json:"value"`
	}
	OrderBy struct {
		Field string `json:"field,omitempty"`
		Desc  bool   `json:"desc,omitempty"`
	}

	Collection struct {
		Name        string  `json:"name"`
		IDFieldName string  `json:"idFieldName"`
		Fields      []Field `json:"fields"`
		Indexes     []Index `json:"indexes"`
	}
	Field struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	Index struct {
		Fields   []string `json:"fields"`
		IsUnique bool     `json:"isUnique"`
	}

	SearchResult struct {
		Revisions []Revision `json:"revisions"`
		SearchId  string     `json:"searchId"`
		Page      int        `json:"page"`
		PerPage   int        `json:"perPage"`
	}
	Revision struct {
		Document      json.RawMessage `json:"document"`
		Revision      string          `json:"revision"`
		TransactionId string          `json:"transactionId"`
	}

	CountResult struct {
		Collection string `json:"collection"`
		Count      int    `json:"count"`
	}

	SearchRequest struct {
		Query    Query  `json:"query,omitempty"`
		Page     int    `json:"page"`
		PerPage  int    `json:"perPage"`
		SearchId string `json:"searchId,omitempty"`
		KeepOpen bool   `json:"keepOpen,omitempty"`
	}
)

func New() *Client {
	return &Client{
		baseURL: "https://vault.immudb.io/ics/api/v1",
		headers: map[string]string{
			"accept":       "application/json",
			"content-type": "application/json",
			// TODO: move API-Key to env
			"X-API-Key": "default.6xtovhtM2-6PCRzdCjjJ5w.F2Tmi0RC8X9qTRBJK2hcLFetIcn1v8CYeWPZsO16JYAlIzXj",
		},
		ledger: "default",
	}
}

func (c Client) addHeaders(req *http.Request) *http.Request {
	for key, val := range c.headers {
		req.Header.Add(key, val)
	}
	return req
}

func (c Client) SearchDocument(ctx context.Context, collection string, query Query, page, perPage int) (ret SearchResult, err error) {
	sr := SearchRequest{
		Page:    page,
		PerPage: perPage,
		Query:   query,
	}
	data, err := json.Marshal(sr)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ledger/%s/collection/%s/documents/search", c.baseURL, c.ledger, collection), bytes.NewReader(data))
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &ret)
	return
}

func (c Client) CountDocument(ctx context.Context, collection string, query Query) (ret CountResult, err error) {
	data, err := json.Marshal(query)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ledger/%s/collection/%s/documents/count", c.baseURL, c.ledger, collection), bytes.NewReader(data))
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &ret)
	return
}

func (c Client) CreateDocument(ctx context.Context, collection string, document any) (err error) {
	data, err := json.Marshal(document)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/ledger/%s/collection/%s/document", c.baseURL, c.ledger, collection), bytes.NewReader(data))
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s\n", body)
		err = errors.New("status code NOT OK")
		return
	}
	return
}

func (c Client) ReplaceDocument(ctx context.Context, collection string, document any, query Query) (err error) {
	reqBody := make(map[string]interface{})
	reqBody["document"] = document
	reqBody["query"] = query
	data, err := json.Marshal(reqBody)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ledger/%s/collection/%s/document", c.baseURL, c.ledger, collection), bytes.NewReader(data))
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s\n", body)
		err = errors.New("status code NOT OK")
		return
	}
	return
}

func (c Client) Collection(ctx context.Context, collection string) (ret Collection, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/ledger/%s/collection/%s", c.baseURL, c.ledger, collection), nil)
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("status code NOT OK")
		return
	}
	err = json.Unmarshal(body, &ret)
	return
}

func (c Client) CollectionCreate(ctx context.Context, collection Collection) (err error) {
	data, err := json.Marshal(collection)
	if err != nil {
		return
	}
	reqBody := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/ledger/%s/collection/%s", c.baseURL, c.ledger, collection.Name), reqBody)
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("status code NOT OK")
		return
	}
	return
}

func (c Client) DropCollection(ctx context.Context, collection string) (err error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/ledger/%s/collection/%s", c.baseURL, c.ledger, collection), nil)
	if err != nil {
		return
	}
	c.addHeaders(req)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("status code NOT OK")
		return
	}
	return
}
