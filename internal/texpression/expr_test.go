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
			"test1exp",
			args{expression: "(1..100)"},
			[]ExprField{
				ExprField{Val: "1..100", Type: exprtype.Expression},
			},
			false,
		},
		{
			"test1",
			args{expression: "hello"},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
			},
			false,
		},
		{
			"test1num",
			args{expression: "$hello.foo[42]"},
			[]ExprField{
				ExprField{Val: "$hello", Type: exprtype.VarName},
				ExprField{Val: "foo", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Index: 42, Type: exprtype.Index},
					},
				},
			},
			false,
		},
		{
			"test1v",
			args{expression: "$hello.foo[$okay]"},
			[]ExprField{
				ExprField{Val: "$hello", Type: exprtype.VarName},
				ExprField{Val: "foo", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "$okay", Type: exprtype.VarName},
					},
				},
			},
			false,
		},
		{
			"test1exprin",
			args{expression: "$hello.foo[(2-3)]"},
			[]ExprField{
				ExprField{Val: "$hello", Type: exprtype.VarName},
				ExprField{Val: "foo", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "2-3", Type: exprtype.Expression},
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
					Val:  "hello",
					Type: exprtype.Literal,
				},
				ExprField{
					Val:  "tourist",
					Type: exprtype.Literal,
				},
			},
			false,
		},
		{
			"test4",
			args{expression: "hello[tourist]"},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test5",
			args{expression: "hello[tourist.dubist]"},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist", Type: exprtype.Literal},
						ExprField{Val: "dubist", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test6",
			args{expression: "hello[tourist[dubist]]"},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{Val: "dubist", Type: exprtype.Literal},
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
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{Val: "dubist", Type: exprtype.Literal},
								ExprField{
									Type: exprtype.InnerRef,
									InnerRefs: []ExprField{
										ExprField{Val: "in", Type: exprtype.Literal},
										ExprField{
											Type: exprtype.InnerRef,
											InnerRefs: []ExprField{
												ExprField{Val: "budapest", Type: exprtype.Literal},
												ExprField{Val: "capitol", Type: exprtype.Literal},
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
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{Val: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
				ExprField{Val: "in", Type: exprtype.Literal},
				ExprField{Val: "budapest", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "capitol", Type: exprtype.Literal},
					},
				},
			},
			false,
		},
		{
			"test9",
			args{expression: "hello[tourist[dubist]].in"},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist", Type: exprtype.Literal},
						ExprField{
							Type: exprtype.InnerRef,
							InnerRefs: []ExprField{
								ExprField{Val: "dubist", Type: exprtype.Literal},
							},
						},
					},
				},
				ExprField{Val: "in", Type: exprtype.Literal},
			},
			false,
		},
		{
			"test10",
			args{expression: `hello["tourist.dubist"]`},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist.dubist", Type: exprtype.Literal},
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
		{
			"test12",
			args{expression: `hello["tourist.dubist"]`},
			[]ExprField{
				ExprField{Val: "hello", Type: exprtype.Literal},
				ExprField{
					Type: exprtype.InnerRef,
					InnerRefs: []ExprField{
						ExprField{Val: "tourist.dubist", Type: exprtype.Literal},
					},
				},
			},
			false,
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
				t.Errorf("NewExpression() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
