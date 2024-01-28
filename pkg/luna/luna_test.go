package luna

import "testing"

func TestCheck(t *testing.T) {

	tests := []struct {
		name          string
		num           string
		wantIsCorrect bool
		wantSum       int
		wantErr       bool
		wantErrSum    bool
	}{
		{"Test_01",
			"4561261212345467",
			true,
			53,
			false,
			false}, {"Test_02",
			"456126121234548",
			true,
			52,
			false,
			false},
		{"Test_03",
			"456126121234546d",
			false,
			53,
			true,
			false},
		{"Test_05",
			"d456126121234546",
			false,
			0,
			true,
			true},
		{"Test_05",
			"1",
			false,
			0,
			true,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := Sum(tt.num[:len(tt.num)-1])
			if (err != nil) != tt.wantErrSum {
				t.Errorf("Sum() error = %v, wantErrSum %v", err, tt.wantErrSum)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("Sum() got = %v, want %v", gotSum, tt.wantSum)
			}

			gotIsCorrect, err := Check(tt.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsCorrect != tt.wantIsCorrect {
				t.Errorf("Check() gotIsCorrect = %v, want %v", gotIsCorrect, tt.wantIsCorrect)
			}
		})
	}
}

//
//func TestSum(t *testing.T) {
//	type args struct {
//		num string
//	}
//	tests := []struct {
//		name    string
//		num     string
//		want    int
//		wantErr bool
//	}{
//		{"Test 01",
//			"456126121234546",
//			53,
//			false},
//		//{"Test 01",
//		//	"456126121234546d",
//		//	0,
//		//	true},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Sum(tt.num)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Sum() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("Sum() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
