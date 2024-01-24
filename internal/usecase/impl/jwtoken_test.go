package impl

import (
	"gophermart/internal/entity"
	"testing"
)

func TestJwtImpl_Get(t *testing.T) {
	type fields struct {
		jwtSecret       string
		jwtExpiresHours int
	}
	type args struct {
		j entity.TokenData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ah := &JwtImpl{
				jwtSecret:       tt.fields.jwtSecret,
				jwtExpiresHours: tt.fields.jwtExpiresHours,
			}
			got, err := ah.Get(tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
