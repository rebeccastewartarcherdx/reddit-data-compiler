package poller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"redditDataCompiler/models"
	"reflect"
	"testing"
)

func Test_poller_fetchRedditData(t *testing.T) {
	goodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "kind": "Listing",
  "data": {
    "after": "t3_1i22j1p",
    "children": [
      {
        "kind": "t3",
        "data": {
          "author_fullname": "t2_aa1ng",
          "title": "My cab driver tonight was so excited to share with me that he’d made the cover of the calendar. I told him I’d help let the world see",
          "name": "t3_7mjw12",
          "author": "the_Diva",
          "ups": 123
        }
      }
    ],
    "before": null
  }
}`)
	}))
	defer goodServer.Close()
	goodRequest, _ := http.NewRequest("GET", goodServer.URL, nil)

	badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error": "too many requests"}`, http.StatusTooManyRequests)
	}))
	defer badServer.Close()
	badRequest, _ := http.NewRequest("GET", badServer.URL, nil)

	type fields struct {
		httpClient *http.Client
		config     envConfig
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.APIResponse
		wantErr bool
	}{
		{
			name: "good resp",
			fields: fields{
				httpClient: goodServer.Client(),
			},
			args: args{
				request: goodRequest,
			},
			want: &models.APIResponse{
				Kind: "Listing",
				Data: models.APIData{
					After: "t3_1i22j1p",
					Children: []models.APIChildren{
						{
							Kind: "t3",
							Data: models.Listing{
								Name:           "t3_7mjw12",
								Ups:            123,
								AuthorFullName: "t2_aa1ng",
								Author:         "the_Diva",
								Title:          "My cab driver tonight was so excited to share with me that he’d made the cover of the calendar. I told him I’d help let the world see",
							},
						},
					},
				},
			},
		},
		{
			name: "bad resp",
			fields: fields{
				httpClient: badServer.Client(),
			},
			args: args{
				request: badRequest,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &poller{
				httpClient: tt.fields.httpClient,
				config:     tt.fields.config,
			}
			got, err := p.fetchRedditData(tt.args.request, tt.args.request.URL.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchRedditData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchRedditData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
