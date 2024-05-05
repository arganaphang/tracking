package pagination

import "testing"

func TestToLimitOffset(t *testing.T) {
	type args struct {
		page    uint
		perPage uint
	}
	tests := []struct {
		name       string
		args       args
		wantLimit  uint
		wantOffset uint
	}{
		{
			name: "first",
			args: args{
				page:    1,
				perPage: 10,
			},
			wantLimit:  10,
			wantOffset: 0,
		},
		{
			name: "second",
			args: args{
				page:    5,
				perPage: 5,
			},
			wantLimit:  5,
			wantOffset: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit, gotOffset := ToLimitOffset(tt.args.page, tt.args.perPage)
			if gotLimit != tt.wantLimit {
				t.Errorf("ToLimitOffset() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("ToLimitOffset() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}
