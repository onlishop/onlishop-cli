package account_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ProducerEndpoint struct {
	c *Client
}

func (c *Client) Producer() (*ProducerEndpoint, error) {
	return &ProducerEndpoint{c: c}, nil
}

type Query struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
type ListExtensionCriteria struct {
	Limit         int    `schema:"limit,omitempty"`
	Offset        int    `schema:"offset,omitempty"`
	OrderBy       string `schema:"orderBy,omitempty"`
	OrderSequence string `schema:"orderSequence,omitempty"`
	Query         *Query `json:"query,omitempty"`
}

func (e ProducerEndpoint) Extensions(ctx context.Context, criteria *ListExtensionCriteria) ([]Extension, error) {
	form := url.Values{}

	if criteria.Limit != 0 {
		form.Set("limit", strconv.Itoa(criteria.Limit))
	}
	if criteria.Offset != 0 {
		form.Set("offset", strconv.Itoa(criteria.Offset))
	}
	if criteria.OrderBy != "" {
		form.Set("orderBy", criteria.OrderBy)
	}
	if criteria.OrderSequence != "" {
		form.Set("orderSequence", criteria.OrderSequence)
	}

	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/plugins?%s", ApiUrl, form.Encode()), nil)

	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, err
	}

	var extensions []Extension
	if err := json.Unmarshal(body, &extensions); err != nil {
		return nil, fmt.Errorf("list_extensions: %v", err)
	}

	return extensions, nil
}

func (e ProducerEndpoint) GetExtensionByName(ctx context.Context, name string) (*Extension, error) {
	criteria := ListExtensionCriteria{
		Query: &Query{
			Type:  "equals",
			Field: "productNumber",
			Value: name,
		},
	}

	extensions, err := e.Extensions(ctx, &criteria)
	if err != nil {
		return nil, err
	}

	for _, ext := range extensions {
		if strings.EqualFold(ext.Name, name) {
			return e.GetExtensionById(ctx, ext.Id)
		}
	}

	return nil, fmt.Errorf("cannot find Extension by name %s", name)
}

func (e ProducerEndpoint) GetExtensionById(ctx context.Context, id string) (*Extension, error) {
	errorFormat := "GetExtensionById: %v"

	// Create it
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/plugins/%s", ApiUrl, id), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var extension Extension
	if err := json.Unmarshal(body, &extension); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return &extension, nil
}

type Extension struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Generation struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"generation"`
	StandardLocale string `json:"standardLocale"`
	Infos          []*struct {
		Id                 int          `json:"id"`
		Locale             string       `json:"locale"`
		Name               string       `json:"name"`
		Description        string       `json:"description"`
		InstallationManual string       `json:"installationManual"`
		ShortDescription   string       `json:"shortDescription"`
		Highlights         string       `json:"highlights"`
		Features           string       `json:"features"`
		MetaTitle          string       `json:"metaTitle"`
		MetaDescription    string       `json:"metaDescription"`
		Tags               []StoreTag   `json:"tags"`
		Videos             []StoreVideo `json:"videos"`
		Faqs               []StoreFaq   `json:"faqs"`
		SupportInfo        interface{}  `json:"supportInfo"`
	} `json:"infos"`
	StoreAvailabilities                 []StoreAvailablity `json:"storeAvailabilities"`
	Categories                          []StoreCategory    `json:"categories"`
	Category                            *StoreCategory     `json:"selectedFutureCategory"`
	Localizations                       []Locale           `json:"localizations"`
	LatestBinary                        interface{}        `json:"latestBinary"`
	AutomaticBugfixVersionCompatibility bool               `json:"automaticBugfixVersionCompatibility"`
	ProductType                         *StoreProductType  `json:"productType"`
	Status                              struct {
		Name string `json:"name"`
	} `json:"status"`
	ReleaseDate                           interface{} `json:"releaseDate"`
	IconURL                               string      `json:"iconUrl"`
	IsCompatibleWithLatestOnlishopVersion bool        `json:"isCompatibleWithLatestOnlishopVersion"`
}

type CreateExtensionRequest struct {
	Name       string `json:"name,omitempty"`
	Generation struct {
		Name string `json:"name"`
	} `json:"generation"`
	ProducerID int `json:"producerId"`
}

func (e ProducerEndpoint) UpdateExtension(ctx context.Context, extension *Extension) error {
	requestBody, err := json.Marshal(extension)
	if err != nil {
		return err
	}
	fmt.Println(string(requestBody))
	// Patch the name
	r, err := e.c.NewAuthenticatedRequest(ctx, "PUT", fmt.Sprintf("%s/plugins/%s", ApiUrl, extension.Id), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e ProducerEndpoint) GetSoftwareVersions(ctx context.Context, generation string) (*SoftwareVersionList, error) {
	errorFormat := "onlishop_versions: %v"
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/pluginstatics/softwareVersions?filter=[{\"property\":\"pluginGeneration\",\"value\":\"%s\"},{\"property\":\"includeNonPublic\",\"value\":\"1\"}]", ApiUrl, generation), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var versions SoftwareVersionList

	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return &versions, nil
}

type SoftwareVersion struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Parent      interface{} `json:"parent"`
	Selectable  bool        `json:"selectable"`
	Major       string      `json:"major"`
	ReleaseDate string      `json:"releaseDate"`
	Status      string      `json:"status"`
}

type Locale struct {
	Name string `json:"name"`
}

type StoreAvailablity struct {
	Name string `json:"name"`
}

type StoreCategory struct {
	Name string `json:"name"`
}

type StoreTag struct {
	Name string `json:"name"`
}

type StoreVideo struct {
	URL string `json:"url"`
}

type StoreProductType struct {
	Name string `json:"name"`
}

type StoreFaq struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Position int    `json:"position"`
}

type ExtensionGeneralInformation struct {
	FutureCategories    []StoreCategory    `json:"futureCategories"`
	Locales             []string           `json:"locales"`
	StoreAvailabilities []StoreAvailablity `json:"storeAvailabilities"`
	ProductTypes        []StoreProductType `json:"productTypes"`
}

func (e ProducerEndpoint) GetExtensionGeneralInfo(ctx context.Context) (*ExtensionGeneralInformation, error) {
	r, err := e.c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/pluginstatics/all", ApiUrl), nil)
	if err != nil {
		return nil, fmt.Errorf("GetExtensionGeneralInfo: %v", err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf("GetExtensionGeneralInfo: %v", err)
	}

	var info *ExtensionGeneralInformation

	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, fmt.Errorf("onlishop_versions: %v", err)
	}

	return info, nil
}
