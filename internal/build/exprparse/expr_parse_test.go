package exprparse

import (
	"reflect"
	"testing"
)

func TestParseExpression(t *testing.T) {
	type args struct {
		fullName string
	}
	tests := []struct {
		name    string
		args    args
		want    ParseTree
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{fullName: "hello"},
			ParseTree{
				Nodes: []Node{
					Node{Val: "hello"},
				},
			},
			false,
		},
		{
			"test1v",
			args{fullName: "$hello.foo"},
			ParseTree{
				Nodes: []Node{
					Node{Val: "$hello"},
					Node{Val: "foo"},
				}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseExpression(tt.args.fullName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
