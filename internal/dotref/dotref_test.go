package dotref

// import (
// 	"reflect"
// 	"testing"
// )

// func TestNewReference(t *testing.T) {
// 	type args struct {
// 		dotName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []RefField
// 		wantErr bool
// 	}{
// 		{
// 			name: "test1",
// 			args: args{
// 				dotName: "hello.me",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{FieldName: "me"},
// 			},
// 		},
// 		{
// 			name: "test2",
// 			args: args{
// 				dotName: "hello.me.more",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{FieldName: "me"},
// 				RefField{FieldName: "more"},
// 			},
// 		},
// 		{
// 			name: "test3",
// 			args: args{
// 				dotName: "hello.me[2]",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{FieldName: "me"},
// 				RefField{Index: 2},
// 			},
// 		},
// 		{
// 			name: "test4",
// 			args: args{
// 				dotName: "hello",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 			},
// 		},
// 		{
// 			name: "test5",
// 			args: args{
// 				dotName: "hello[2]",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{Index: 2},
// 			},
// 		},
// 		{
// 			name: "test6",
// 			args: args{
// 				dotName: "hello.yolo[2].me",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{FieldName: "yolo"},
// 				RefField{Index: 2},
// 				RefField{FieldName: "me"},
// 			},
// 		},
// 		{
// 			name: "test7",
// 			args: args{
// 				dotName: "hello.yolo[2].me[3].okay.okay[2]",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{FieldName: "yolo"},
// 				RefField{Index: 2},
// 				RefField{FieldName: "me"},
// 				RefField{Index: 3},
// 				RefField{FieldName: "okay"},
// 				RefField{FieldName: "okay"},
// 				RefField{Index: 2},
// 			},
// 		},
// 		{
// 			name: "test8",
// 			args: args{
// 				dotName: "hello.yolo[2].me[3].okay.okay[2].ok[good]",
// 			},
// 			want: []RefField{
// 				RefField{FieldName: "hello"},
// 				RefField{FieldName: "yolo"},
// 				RefField{Index: 2},
// 				RefField{FieldName: "me"},
// 				RefField{Index: 3},
// 				RefField{FieldName: "okay"},
// 				RefField{FieldName: "okay"},
// 				RefField{Index: 2},
// 				RefField{FieldName: "ok"},
// 				RefField{FieldName: "good"},
// 			},
// 		},
// 		{
// 			name: "test-err1",
// 			args: args{
// 				dotName: "hello]",
// 			},
// 			want:    []RefField{},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseFields(tt.args.dotName)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("NewReference() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewReference() = %+v, want %+v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestParseFieldsRec(t *testing.T) {
// 	type args struct {
// 		dotName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []RefField
// 		wantErr bool
// 	}{
// 		{
// 			"test1",
// 			args{dotName: "hello"},
// 			[]RefField{
// 				RefField{
// 					FieldName: "hello",
// 					Type:      Literal,
// 				},
// 			},
// 			false,
// 		},
// 		{
// 			"test2",
// 			args{dotName: "hello.tourist"},
// 			[]RefField{
// 				RefField{FieldName: "hello", Type: Literal},
// 				RefField{FieldName: "tourist", Type: Literal},
// 			},
// 			false,
// 		},
// 		{
// 			"test3",
// 			args{dotName: "hello[tourist]"},
// 			[]RefField{
// 				RefField{FieldName: "hello", Type: Literal},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{
// 							FieldName: "tourist", Type: Literal,
// 						},
// 					},
//					Type: InnerRef,
// 				},
// 			},
// 			false,
// 		},
// 		{
// 			"test4",
// 			args{dotName: "hello[tourist.dubist]"},
// 			[]RefField{
// 				RefField{FieldName: "hello", Type: Literal},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "tourist", Type: Literal},
// 						RefField{FieldName: "dubist", Type: Literal},
// 					}, Type: InnerRef,
// 				},
// 			},
// 			false,
// 		},
// 		{
// 			"test5",
// 			args{dotName: "hello[tourist][dubist]"},
// 			[]RefField{
// 				RefField{FieldName: "hello", Type: Literal},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "tourist", Type: Literal},
// 					}, Type: InnerRef,
// 				},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "dubist", Type: Literal},
// 					}, Type: InnerRef,
// 				},
// 			},
// 			false,
// 		},
// 		{
// 			"test6",
// 			args{dotName: "hello[tourist][dubist].in.budapest"},
// 			[]RefField{
// 				RefField{FieldName: "hello", Type: Literal},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "tourist", Type: Literal},
// 					}, Type: InnerRef,
// 				},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "dubist", Type: Literal},
// 					}, Type: InnerRef,
// 				},
// 				RefField{FieldName: "in", Type: Literal},
// 				RefField{FieldName: "budapest", Type: Literal},
// 			},
// 			false,
// 		},
// 		{
// 			"test7",
// 			args{dotName: "hello[tourist][dubist.oder[nicht]].in.budapest"},
// 			[]RefField{
// 				RefField{FieldName: "hello", Type: Literal},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "tourist", Type: Literal},
// 					}, Type: InnerRef,
// 				},
// 				RefField{
// 					InnerRef: []RefField{
// 						RefField{FieldName: "dubist", Type: Literal},
// 						RefField{FieldName: "oder", Type: Literal},
// 						RefField{
// 							InnerRef: []RefField{
// 								RefField{FieldName: "nicht", Type: Literal},
// 							}, Type: InnerRef,
// 						},
// 					}, Type: InnerRef,
// 				},
// 				RefField{FieldName: "in", Type: Literal},
// 				RefField{FieldName: "budapest", Type: Literal},
// 			},
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseFieldsRec(tt.args.dotName)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ParseFieldsRec() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ParseFieldsRec() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
