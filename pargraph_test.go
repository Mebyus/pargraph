package pargraph

import (
	"reflect"
	"testing"
)

func TestMakeReachableTree(t *testing.T) {
	type args struct {
		originID int64
		rows     []*SourceRow
	}
	tests := []struct {
		name     string
		args     args
		wantTree *TreeNode
		wantErr  bool
	}{
		{
			name: "single tree",
			args: args{
				originID: 4,
				rows: []*SourceRow{
					{
						ID:        10,
						HasParent: true,
						ParentID:  5,
					},
					{
						ID:        5,
						HasParent: true,
						ParentID:  3,
					},
					{
						ID:        4,
						HasParent: true,
						ParentID:  3,
					},
					{
						ID:        3,
						HasParent: false,
					},
				},
			},
			wantTree: &TreeNode{
				ID: 3,
				Children: []*TreeNode{
					{
						ID: 4,
					},
					{
						ID: 5,
						Children: []*TreeNode{
							{
								ID: 10,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "looped, no root",
			args: args{
				originID: 4,
				rows: []*SourceRow{
					{
						ID:        10,
						HasParent: true,
						ParentID:  4,
					},
					{
						ID:        5,
						HasParent: true,
						ParentID:  4,
					},
					{
						ID:        4,
						HasParent: true,
						ParentID:  3,
					},
					{
						ID:        3,
						HasParent: true,
						ParentID:  5,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "two trees, start at bottom",
			args: args{
				originID: 8,
				rows: []*SourceRow{
					{
						ID:        10,
						HasParent: true,
						ParentID:  5,
					},
					{
						ID:        5,
						HasParent: true,
						ParentID:  3,
					},
					{
						ID:        4,
						HasParent: true,
						ParentID:  3,
					},
					{
						ID:        3,
						HasParent: false,
					},
					{
						ID:        2,
						HasParent: false,
					},
					{
						ID:        1,
						HasParent: true,
						ParentID:  2,
					},
					{
						ID:        6,
						HasParent: true,
						ParentID:  2,
					},
					{
						ID:        7,
						HasParent: true,
						ParentID:  1,
					},
					{
						ID:        8,
						HasParent: true,
						ParentID:  1,
					},
					{
						ID:        9,
						HasParent: true,
						ParentID:  1,
					},
				},
			},
			wantTree: &TreeNode{
				ID: 2,
				Children: []*TreeNode{
					{
						ID: 1,
						Children: []*TreeNode{
							{
								ID: 8,
							},
							{
								ID: 7,
							},
							{
								ID: 9,
							},
						},
					},
					{
						ID: 6,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTree, err := MakeReachableTree(tt.args.originID, tt.args.rows)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeReachableTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTree, tt.wantTree) {
				t.Errorf("MakeReachableTree() = %v, want %v", gotTree, tt.wantTree)
			}
		})
	}
}
