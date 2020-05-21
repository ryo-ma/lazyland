package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
)

type Client struct {
	RepositoryURL *url.URL
	HTTPClient    *http.Client
}

type Item struct {
	Type  string `json:"type"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Desc  string `json:"desc"`
}

type Readme struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"`
}

//type Items map[string]Item
type Result struct {
	Items []Item
}

func (item *Item) GetRepositoryName() string {
	return item.Owner + "/" + item.Repo
}

func (item *Item) String() string {
	const templateText = `
	Repository : {{.GetRepositoryName}}
	Description: {{.Desc}}
	`
	template, err := template.New("Repository").Parse(templateText)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	if err := template.Execute(&doc, item); err != nil {
		panic(err)
	}
	return doc.String()
}

func (result *Result) Draw(writer io.Writer) error {
	for _, item := range result.Items {
		//starText := " ⭐️ " + strconv.Itoa(item.GetStars())
		fmt.Fprintf(writer, "\033[32m%s\033[0m\n", item.Repo)
	}
	return nil
}

func NewClient() (*Client, error) {
	repositoryURL, err := url.Parse("https://raw.githubusercontent.com/denoland/deno_website2/master/database.json")
	if err != nil {
		return nil, err
	}
	return &Client{
		RepositoryURL: repositoryURL,
		HTTPClient:    http.DefaultClient,
	}, nil
}

//func (client *Client) GetReadme(item Item) (*Readme, error) {
//	url := *client.OfficialURL
//	url.Path = path.Join(url.Path, "repos", item.GetRepositoryName(), "readme")
//	req, err := http.NewRequest("GET", url.String(), nil)
//	if err != nil {
//		panic(err)
//	}
//	req.Header.Add("Accept", "application/vnd.github.mercy-preview+json")
//	resp, err := client.HTTPClient.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	var readme *Readme
//	if err = json.Unmarshal(body, &readme); err != nil {
//		return nil, err
//	}
//	return readme, nil
//}
func (client *Client) GetThirdPartyRepositores() (*Result, error) {
	url := client.RepositoryURL
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var itemMap map[string]Item
	if err = json.Unmarshal(body, &itemMap); err != nil {
		return nil, err
	}
	items := []Item{}
	for _, v := range itemMap {
		items = append(items, v)
	}
	return &Result{
		Items: items,
	}, nil
}
