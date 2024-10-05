package frontend

import (
	"reflect"
	"testing"
)

func Test_skipWhitespace(t *testing.T) {
	type args struct {
		code     string
		startPos Position
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 Position
	}{
		{"no whitespace",
			args{
				"int main()",
				Position{1, 1},
			},
			"int main()",
			Position{1, 1},
		},
		{"with whitespace",
			args{
				"    \tint main()",
				Position{1, 1},
			},
			"int main()",
			Position{1, 6},
		},
		{"with newline",
			args{
				"    \n int main()",
				Position{1, 1},
			},
			"int main()",
			Position{2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := skipWhitespace(tt.args.code, tt.args.startPos)
			if got != tt.want {
				t.Errorf("skipWhitespace() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("skipWhitespace() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    []TokenType
		wantErr bool
	}{
		{
			"it works",
			args{
				`int main() {
						return 2;
			 		}`,
			},
			[]TokenType{
				TokTypeInt,
				TokTypeIdentifier,
				TokTypeLeftParen,
				TokTypeRightParen,
				TokTypeLeftBrace,
				TokTypeReturn,
				TokTypeIntConstant,
				TokTypeSemicolon,
				TokTypeRightBrace,
			},
			false,
		},
		{
			"end before expr",
			args{
				`int main() {
						return`,
			},
			[]TokenType{
				TokTypeInt,
				TokTypeIdentifier,
				TokTypeLeftParen,
				TokTypeRightParen,
				TokTypeLeftBrace,
				TokTypeReturn,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotTypes []TokenType
			for _, token := range got {
				gotTypes = append(gotTypes, token.tokenType)
			}

			if !reflect.DeepEqual(gotTypes, tt.want) {
				t.Errorf("Tokenize() got = %v, want %v", gotTypes, tt.want)
			}
		})
	}
}
