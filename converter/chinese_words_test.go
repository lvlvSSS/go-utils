package converter

import "testing"

func TestWordToPinyin(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{
			name:    "test",
			text:    "测试",
			want:    "ce shi",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WordToPinyin(tt.text)
			if (err != nil) != tt.wantErr {
				t.Logf(got)
				t.Errorf("WordToPinyin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WordToPinyin() got = %v, want %v", got, tt.want)
			}
		})
	}
}
