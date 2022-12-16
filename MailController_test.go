package main

import "testing"

func TestSendMail(t *testing.T) {
	type args struct {
		body    string
		to      []string
		subject string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				body:    "Test 1",
				to:      []string{"nicolasxd147@gmail.com"},
				subject: "Test 1",
			},
			wantErr: false,
		},
		//{
		//	name: "Test 2",
		//	args: args{
		//		body:    "Test 2",
		//		to:      "gabrielborgescast@gmail.com",
		//		subject: "Test 2",
		//	},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMail(tt.args.body, tt.args.subject, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
