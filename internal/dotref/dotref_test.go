package dotref

import (
	"reflect"
	"testing"
)

func TestNewReference(t *testing.T) {
	type args struct {
		dotName string
		value   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    DotRef
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				dotName: "hello.me",
				value:   nil,
			},
			want: DotRef{
				Name: "hello",
				Fields: []Field{
					Field{FieldName: "me"},
				},
			},
		},
		{
			name: "test2",
			args: args{
				dotName: "hello.me.more",
				value:   nil,
			},
			want: DotRef{
				Name: "hello",
				Fields: []Field{
					Field{FieldName: "me"},
					Field{FieldName: "more"},
				},
			},
		},
		{
			name: "test3",
			args: args{
				dotName: "hello.me[2]",
				value:   nil,
			},
			want: DotRef{
				Name: "hello",
				Fields: []Field{
					Field{FieldName: "me"},
					Field{Index: 2},
				},
			},
		},
		{
			name: "test4",
			args: args{
				dotName: "hello",
				value:   nil,
			},
			want: DotRef{
				Name:   "hello",
				Fields: []Field{},
			},
		},
		{
			name: "test5",
			args: args{
				dotName: "hello[2]",
				value:   nil,
			},
			want: DotRef{
				Name: "hello",
				Fields: []Field{
					Field{Index: 2},
				},
			},
		},
		{
			name: "test6",
			args: args{
				dotName: "hello.yolo[2].me",
				value:   nil,
			},
			want: DotRef{
				Name: "hello",
				Fields: []Field{
					Field{FieldName: "yolo"},
					Field{Index: 2},
					Field{FieldName: "me"},
				},
			},
		},
		{
			name: "test7",
			args: args{
				dotName: "hello.yolo[2].me[3].okay.okay[2]",
				value:   nil,
			},
			want: DotRef{
				Name: "hello",
				Fields: []Field{
					Field{FieldName: "yolo"},
					Field{Index: 2},
					Field{FieldName: "me"},
					Field{Index: 3},
					Field{FieldName: "okay"},
					Field{FieldName: "okay"},
					Field{Index: 2},
				},
			},
		},
		{
			name: "test-err1",
			args: args{
				dotName: "hello]",
				value:   nil,
			},
			want: DotRef{
				Name:   "hello]",
				Fields: []Field{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReference(tt.args.dotName, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReference() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
