package valueobject

import "testing"

func TestNewBookStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    BookStatus
		wantErr bool
	}{
		{"unread", "unread", Unread, false},
		{"reading", "reading", Reading, false},
		{"finished", "finished", Finished, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBookStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBookStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewBookStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookStatus_String(t *testing.T) {
	if Unread.String() != "unread" {
		t.Errorf("Unread.String() = %v, want unread", Unread.String())
	}
	if Reading.String() != "reading" {
		t.Errorf("Reading.String() = %v, want reading", Reading.String())
	}
	if Finished.String() != "finished" {
		t.Errorf("Finished.String() = %v, want finished", Finished.String())
	}
}

func TestBookStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from BookStatus
		to   BookStatus
		want bool
	}{
		{"unread to reading", Unread, Reading, true},
		{"unread to finished", Unread, Finished, false},
		{"unread to unread", Unread, Unread, false},
		{"reading to finished", Reading, Finished, true},
		{"reading to unread", Reading, Unread, false},
		{"reading to reading", Reading, Reading, false},
		{"finished to unread", Finished, Unread, false},
		{"finished to reading", Finished, Reading, false},
		{"finished to finished", Finished, Finished, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
