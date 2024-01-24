package impl

import (
	"github.com/stretchr/testify/require"
	"gophermart/internal/entity"
	"testing"
)

func TestPwdImpl(t *testing.T) {
	type args struct {
		pwd      entity.UserPassword
		pwdCheck entity.UserPassword
	}
	const (
		checkHash = iota
		checkIsCorrect
	)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Test #1", args{"123", "123"}, true},
		{"Test #2", args{"123", "12"}, false},
		{"Test #2", args{"SuperComplexPassword!23u348u3#*#@#", "SuperComplexPassword!23u348u3#*#@#"}, true},
		{"Test #2", args{"SuperComplexPassword!23u348u3#*#@#", "SuperComplexPassword!23u348u3#*#@"}, false},
		{"Test #2", args{"", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ah := &PwdImpl{}
			gotHash, err := ah.Hash(tt.args.pwd)
			require.NoError(t, err)
			require.NotEqual(t, gotHash, "")

			got := ah.IsCorrect(tt.args.pwdCheck, gotHash)
			require.Equal(t, got, tt.want)
		})
	}
}
