package texpression

import (
	"reflect"
	"testing"

	"github.com/kassybas/tame/types/exprtype"
)

func TestNewExpression(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name    string
		args    args
		want    []ExprField
		wantErr bool
	}{
		{
			"test1",
			args{expression: "hello"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
			},
			false,
		},

		{
			"test1v",
			args{expression: "$hello.foo[$okay]"},
			[]ExprField{
				ExprField{FieldName: "$hello", Type: exprtype.VarName},
				ExprField{FieldName: "foo", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "$okay", Type: exprtype.VarName},
					},
				},
			},
			false,
		},
		{
			"test2",
			args{expression: "hello.tourist"},
			[]ExprField{
				ExprField{
					FieldName: "hello",
					Type:      exprtype.Literal,
				},
				ExprField{
					FieldName: "tourist",
					Type:      exprtype.Literal,
				},
			},
			false,
		},
		{
			"test4",
			args{expression: "hello[tourist]"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test5",
			args{expression: "hello[tourist.dubist]"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist", Type: exprtype.Literal},
						ExprField{FieldName: "dubist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test6",
			args{expression: "hello[tourist[dubist]]"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{FieldName: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
			},
			false,
		},
		{
			"test7",
			args{expression: "hello[tourist[dubist[in[budapest.capitol]]]]"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{FieldName: "dubist", Type: exprtype.Literal},
								ExprField{
									Type: exprtype.InnerRef,
									InnerRefs: []ExprField{
										ExprField{FieldName: "in", Type: exprtype.Literal},
										ExprField{
											Type: exprtype.InnerRef,
											InnerRefs: []ExprField{
												ExprField{FieldName: "budapest", Type: exprtype.Literal},
												ExprField{FieldName: "capitol", Type: exprtype.Literal},
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
			args{expression: "hello[tourist[dubist]].in.budapest[capitol]"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{FieldName: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
				ExprField{FieldName: "in", Type: exprtype.Literal},
				ExprField{FieldName: "budapest", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "capitol", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test9",
			args{expression: "hello[tourist[dubist]].in"},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{FieldName: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
				ExprField{FieldName: "in", Type: exprtype.Literal},
			},
			false,
		},
		{
			"test10",
			args{expression: `hello["tourist.dubist"]`},
			[]ExprField{
				ExprField{FieldName: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{FieldName: "tourist.dubist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test11",
			args{expression: `hello[tourist[dubist]`},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExpression(tt.args.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
