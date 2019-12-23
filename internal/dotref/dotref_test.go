package dotref

import (
	"reflect"
	"testing"
)

func TestNewReference(t *testing.T) {
	type args struct {
		dotName string
	}
	tests := []struct {
		name    string
		args    args
		want    []RefField
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				dotName: "hello.me",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
				RefField{FieldName: "me"},
			},
		},
		{
			name: "test2",
			args: args{
				dotName: "hello.me.more",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
				RefField{FieldName: "me"},
				RefField{FieldName: "more"},
			},
		},
		{
			name: "test3",
			args: args{
				dotName: "hello.me[2]",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
				RefField{FieldName: "me"},
				RefField{Index: 2},
			},
		},
		{
			name: "test4",
			args: args{
				dotName: "hello",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
			},
		},
		{
			name: "test5",
			args: args{
				dotName: "hello[2]",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
				RefField{Index: 2},
			},
		},
		{
			name: "test6",
			args: args{
				dotName: "hello.yolo[2].me",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
				RefField{FieldName: "yolo"},
				RefField{Index: 2},
				RefField{FieldName: "me"},
			},
		},
		{
			name: "test7",
			args: args{
				dotName: "hello.yolo[2].me[3].okay.okay[2]",
			},
			want: []RefField{
				RefField{FieldName: "hello"},
				RefField{FieldName: "yolo"},
				RefField{Index: 2},
				RefField{FieldName: "me"},
				RefField{Index: 3},
				RefField{FieldName: "okay"},
				RefField{FieldName: "okay"},
				RefField{Index: 2},
			},
		},
		{
			name: "test-err1",
			args: args{
				dotName: "hello]",
			},
			want:    []RefField{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFields(tt.args.dotName)
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
