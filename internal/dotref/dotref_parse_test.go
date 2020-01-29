package dotref

import (
	"reflect"
	"testing"

	"github.com/kassybas/tame/types/exprtype"
)

func TestParseDotRef(t *testing.T) {
	type args struct {
		fullName string
	}
	tests := []struct {
		name    string
		args    args
		want    []RefField
		wantErr bool
	}{
		{
			"test1",
			args{fullName: "hello"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
			},
			false,
		},
		{
			"test1v",
			args{fullName: "$hello.foo[$okay]"},
			[]RefField{
				RefField{FieldName: "$hello", Type: exprtype.VarName},
				RefField{FieldName: "foo", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "$okay", Type: exprtype.VarName},
					},
				},
			},
			false,
		},
		{
			"test2",
			args{fullName: "hello.tourist"},
			[]RefField{
				RefField{
					FieldName: "hello",
					Type:      exprtype.Literal,
				},
				RefField{
					FieldName: "tourist",
					Type:      exprtype.Literal,
				},
			},
			false,
		},
		{
			"test4",
			args{fullName: "hello[tourist]"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test5",
			args{fullName: "hello[tourist.dubist]"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist", Type: exprtype.Literal},
						RefField{FieldName: "dubist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test6",
			args{fullName: "hello[tourist[dubist]]"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist", Type: exprtype.Literal},
						RefField{
							Type: exprtype.InnerRef,
							InnerRefs: []RefField{
								RefField{FieldName: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
			},
			false,
		},
		{
			"test7",
			args{fullName: "hello[tourist[dubist[in[budapest.capitol]]]]"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist", Type: exprtype.Literal},
						RefField{
							Type: exprtype.InnerRef,
							InnerRefs: []RefField{
								RefField{FieldName: "dubist", Type: exprtype.Literal},
								RefField{
									Type: exprtype.InnerRef,
									InnerRefs: []RefField{
										RefField{FieldName: "in", Type: exprtype.Literal},
										RefField{
											Type: exprtype.InnerRef,
											InnerRefs: []RefField{
												RefField{FieldName: "budapest", Type: exprtype.Literal},
												RefField{FieldName: "capitol", Type: exprtype.Literal},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"test8",
			args{fullName: "hello[tourist[dubist]].in.budapest[capitol]"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist", Type: exprtype.Literal},
						RefField{
							Type: exprtype.InnerRef,
							InnerRefs: []RefField{
								RefField{FieldName: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
				RefField{FieldName: "in", Type: exprtype.Literal},
				RefField{FieldName: "budapest", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "capitol", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test9",
			args{fullName: "hello[tourist[dubist]].in"},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist", Type: exprtype.Literal},
						RefField{
							Type: exprtype.InnerRef,
							InnerRefs: []RefField{
								RefField{FieldName: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
				RefField{FieldName: "in", Type: exprtype.Literal},
			},
			false,
		},
		{
			"test10",
			args{fullName: `hello["tourist.dubist"]`},
			[]RefField{
				RefField{FieldName: "hello", Type: exprtype.Literal},
				RefField{
					Type: exprtype.InnerRef,
					InnerRefs: []RefField{
						RefField{FieldName: "tourist.dubist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test11",
			args{fullName: `hello[tourist[dubist]`},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVarRef(tt.args.fullName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDotRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDotRef() = \n\t\t%+v, \nwant\t%+v", got, tt.want)
			}
		})
	}
}
