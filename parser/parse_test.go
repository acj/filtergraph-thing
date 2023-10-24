package parser

import (
	"reflect"
	"testing"
)

func TestParseFiltergraph(t *testing.T) {
	type args struct {
		rawFiltergraph string
	}
	tests := []struct {
		name    string
		args    args
		want    *Filtergraph
		wantErr bool
	}{
		{"empty", args{""}, &Filtergraph{}, true},
		{
			"scale filter only",
			args{"scale"},
			&Filtergraph{
				Filterchains: []*Filterchain{
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scale",
							},
						},
					},
				},
			},
			false,
		},
		{
			"scale filter with args",
			args{"scale=w=100:h=200"},
			&Filtergraph{
				Filterchains: []*Filterchain{
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scale",
								Args: []SingleArg{
									SingleArg{Name: ptr("w"), Value: ptr("100")},
									SingleArg{Name: ptr("h"), Value: ptr("200")},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"scale filter with args and input/output links",
			args{"[in]scale=w=100:h=200[out]"},
			&Filtergraph{
				Filterchains: []*Filterchain{
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scale",
								Args: []SingleArg{
									SingleArg{Name: ptr("w"), Value: ptr("100")},
									SingleArg{Name: ptr("h"), Value: ptr("200")},
								},
								InLinks:  []string{"in"},
								OutLinks: []string{"out"},
							},
						},
					},
				},
			},
			false,
		},
		{
			"scale filter with quoted args and input/output links",
			args{"[in]scale='w=100:h=200'[out]"},
			&Filtergraph{
				Filterchains: []*Filterchain{
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scale",
								Args: []SingleArg{
									SingleArg{Name: ptr("w"), Value: ptr("100")},
									SingleArg{Name: ptr("h"), Value: ptr("200")},
								},
								InLinks:  []string{"in"},
								OutLinks: []string{"out"},
							},
						},
					},
				},
			},
			false,
		},
		{
			"multiple filterchains",
			args{"scale; reverse; scroll"},
			&Filtergraph{
				Filterchains: []*Filterchain{
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scale",
							},
						},
					},
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "reverse",
							},
						},
					},
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scroll",
							},
						},
					},
				},
			},
			false,
		},
		{
			"multiple filterchains with whitespace",
			args{"scale;\nreverse;\n\tscroll"},
			&Filtergraph{
				Filterchains: []*Filterchain{
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scale",
							},
						},
					},
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "reverse",
							},
						},
					},
					&Filterchain{
						Filters: []*Filter{
							&Filter{
								Filter: "scroll",
							},
						},
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFiltergraph(tt.args.rawFiltergraph)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFiltergraph() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFiltergraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptr[s any](v s) *s { return &v }
