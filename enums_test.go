package silcomms

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    Status
		want string
	}{
		{
			name: "success",
			e:    StatusSuccess,
			want: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Status
		want bool
	}{
		{
			name: "valid type",
			e:    StatusSuccess,
			want: true,
		},
		{
			name: "invalid type",
			e:    Status("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Status.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
