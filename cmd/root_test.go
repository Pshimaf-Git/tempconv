package cmd

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestDegreeType_String(t *testing.T) {
	testCases := []struct {
		desc  string
		dType DegreeType

		want string
	}{
		{
			desc:  "Celsius",
			dType: Celsius,

			want: "C",
		},
		{
			desc:  "Fahrenheit",
			dType: Forenheit,

			want: "F",
		},
		{
			desc:  "Kelvin",
			dType: Kelvin,

			want: "K",
		},

		{
			desc:  "Unknown",
			dType: 0,

			want: "unknown",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.dType.String()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want '%s' got '%s'", tt.want, got)
			}
		})
	}
}

func TestParseTempratureSys(t *testing.T) {
	testCases := []struct {
		desc  string
		input string

		wantErr string
		wantVal DegreeType
	}{
		{
			desc:  "one Kelvin in up register",
			input: "1K",

			wantErr: "",
			wantVal: Kelvin,
		},
		{
			desc:  "one Kelvin in low register",
			input: "1k",

			wantErr: "",
			wantVal: Kelvin,
		},
		{
			desc:  "one Kelvin in up register with °",
			input: "1°K",

			wantErr: "",
			wantVal: Kelvin,
		},
		{
			desc:  "one Kelvin in low register with °",
			input: "1°k",

			wantErr: "",
			wantVal: Kelvin,
		},

		{
			desc:  "10000 Kelvin in up register",
			input: "10000K",

			wantErr: "",
			wantVal: Kelvin,
		},

		{
			desc:  "10000 Kelvin in low register",
			input: "10000k",

			wantErr: "",
			wantVal: Kelvin,
		},

		{
			desc:  "one Celsius in up register",
			input: "1C",

			wantErr: "",
			wantVal: Celsius,
		},
		{
			desc:  "one Celsius in up register with °",
			input: "1°C",

			wantErr: "",
			wantVal: Celsius,
		},
		{
			desc:  "one Celsius in low register with °",
			input: "1°c",

			wantErr: "",
			wantVal: Celsius,
		},
		{
			desc:  "one Celsius in low register",
			input: "1c",

			wantErr: "",
			wantVal: Celsius,
		},

		{
			desc:  "10000 Celsius in up register",
			input: "10000C",

			wantErr: "",
			wantVal: Celsius,
		},

		{
			desc:  "10000 Celsius in low register",
			input: "10000c",

			wantErr: "",
			wantVal: Celsius,
		},

		{
			desc:  "one Forenheit in up register",
			input: "1F",

			wantErr: "",
			wantVal: Forenheit,
		},
		{
			desc:  "one Forenheit in low register",
			input: "1f",

			wantErr: "",
			wantVal: Forenheit,
		},
		{
			desc:  "one Forenheit in up register with °",
			input: "1°F",

			wantErr: "",
			wantVal: Forenheit,
		},
		{
			desc:  "one Forenheit in low register with °",
			input: "1°f",

			wantErr: "",
			wantVal: Forenheit,
		},

		{
			desc:  "10000 Forenheit in up register",
			input: "10000F",

			wantErr: "",
			wantVal: Forenheit,
		},

		{
			desc:  "10000 Forenheit in low register",
			input: "10000f",

			wantErr: "",
			wantVal: Forenheit,
		},

		// BAD CASES

		{
			desc:  "empty string",
			input: "",

			wantErr: "ParseTempratureSys: empty",
		},
		{
			desc:  "invalid temprature system",
			input: "10G",

			wantErr: "ParseTempratureSys: invalid",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			gotVal, gotErr := ParseTempratureSys(tt.input)
			checkErr(t, gotErr, tt.wantErr)

			if !reflect.DeepEqual(tt.wantVal, gotVal) {
				t.Errorf("expected '%v', got: %v", tt.wantVal, gotVal)
			}
		})
	}
}

func TestParseTemprature(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		wantErr string
		wantVal Degree
	}{
		// Kelvin
		{
			desc:    "zero Kelvin",
			input:   "0K",
			wantVal: Degree{Type: Kelvin, Value: 0.0},
		},
		{
			desc:    "negative Kelvin (physically invalid but parseable)",
			input:   "-10K",
			wantVal: Degree{Type: Kelvin, Value: -10.0},
		},
		{
			desc:    "decimal Kelvin",
			input:   "273.15K",
			wantVal: Degree{Type: Kelvin, Value: 273.15},
		},
		{
			desc:    "Kelvin with degree symbol and lower case k",
			input:   "100°k",
			wantVal: Degree{Type: Kelvin, Value: 100.0},
		},
		// Celsius
		{
			desc:    "zero Celsius",
			input:   "0C",
			wantVal: Degree{Type: Celsius, Value: 0.0},
		},
		{
			desc:    "negative Celsius",
			input:   "-40C",
			wantVal: Degree{Type: Celsius, Value: -40.0},
		},
		{
			desc:    "decimal Celsius",
			input:   "36.6C",
			wantVal: Degree{Type: Celsius, Value: 36.6},
		},
		{
			desc:    "Celsius with degree symbol and lower case",
			input:   "25°c",
			wantVal: Degree{Type: Celsius, Value: 25.0},
		},
		// Fahrenheit
		{
			desc:    "zero Fahrenheit",
			input:   "0F",
			wantVal: Degree{Type: Forenheit, Value: 0.0},
		},
		{
			desc:    "negative Fahrenheit",
			input:   "-4F",
			wantVal: Degree{Type: Forenheit, Value: -4.0},
		},
		{
			desc:    "decimal Fahrenheit",
			input:   "98.6F",
			wantVal: Degree{Type: Forenheit, Value: 98.6},
		},
		{
			desc:    "Fahrenheit with degree symbol and lower case",
			input:   "212°f",
			wantVal: Degree{Type: Forenheit, Value: 212.0},
		},
		// Edge cases with spaces
		{
			desc:    "leading spaces",
			input:   "  100C",
			wantVal: Degree{Type: Celsius, Value: 100.0},
		},
		{
			desc:    "trailing spaces",
			input:   "100F  ",
			wantVal: Degree{Type: Forenheit, Value: 100.0},
		},
		{
			desc:    "space between number and unit",
			input:   "100 K",
			wantVal: Degree{Type: Kelvin, Value: 100.0},
		},
		// Invalid cases
		{
			desc:    "empty string",
			input:   "",
			wantErr: "empty string",
		},
		{
			desc:    "only unit without number",
			input:   "C",
			wantErr: "failed to parse number",
		},
		{
			desc:    "invalid unit",
			input:   "100X",
			wantErr: "invalid temperature specifier",
		},
		{
			desc:    "multiple units",
			input:   "10CF",
			wantErr: "",
			wantVal: Degree{Type: Forenheit, Value: 10.0},
		},
		{
			desc:    "non-numeric value",
			input:   "abcC",
			wantErr: "failed to parse number",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			gotVal, gotErr := ParseTemprature(tt.input)
			checkErr(t, gotErr, tt.wantErr)

			if tt.wantErr == "" && !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("expected %v, got %v", tt.wantVal, gotVal)
			}
		})
	}
}

// New test function for convertFromCelsius
func TestConvertFromCelsius(t *testing.T) {
	testCases := []struct {
		desc    string
		target  DegreeType
		from    Degree
		want    Degree
		wantErr string
	}{
		{
			desc:   "Celsius to Celsius",
			target: Celsius,
			from:   Degree{Type: Celsius, Value: 25.0},
			want:   Degree{Type: Celsius, Value: 25.0},
		},
		{
			desc:   "Celsius to Fahrenheit",
			target: Forenheit,
			from:   Degree{Type: Celsius, Value: 0.0},
			want:   Degree{Type: Forenheit, Value: 32.0},
		},
		{
			desc:   "Celsius to Fahrenheit (negative)",
			target: Forenheit,
			from:   Degree{Type: Celsius, Value: -40.0},
			want:   Degree{Type: Forenheit, Value: -40.0},
		},
		{
			desc:   "Celsius to Kelvin",
			target: Kelvin,
			from:   Degree{Type: Celsius, Value: 100.0},
			want:   Degree{Type: Kelvin, Value: 373.15},
		},
		{
			desc:   "Celsius to Kelvin (zero)",
			target: Kelvin,
			from:   Degree{Type: Celsius, Value: 0.0},
			want:   Degree{Type: Kelvin, Value: 273.15},
		},
		{
			desc:    "wrong source type (Fahrenheit)",
			target:  Celsius,
			from:    Degree{Type: Forenheit, Value: 32.0},
			wantErr: "expected Celsius",
		},
		{
			desc:    "wrong source type (Kelvin)",
			target:  Celsius,
			from:    Degree{Type: Kelvin, Value: 300.0},
			wantErr: "expected Celsius",
		},
		{
			desc:    "unsaported target",
			target:  0,
			from:    Degree{Type: Celsius, Value: -1.0},
			wantErr: "unsupported target type",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := convertFromCelsius(tt.target, tt.from)
			checkErr(t, err, tt.wantErr)
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

// New test function for convertFromForenheit
func TestConvertFromForenheit(t *testing.T) {
	testCases := []struct {
		desc    string
		target  DegreeType
		from    Degree
		want    Degree
		wantErr string
	}{
		{
			desc:   "Fahrenheit to Fahrenheit",
			target: Forenheit,
			from:   Degree{Type: Forenheit, Value: 212.0},
			want:   Degree{Type: Forenheit, Value: 212.0},
		},
		{
			desc:   "Fahrenheit to Celsius",
			target: Celsius,
			from:   Degree{Type: Forenheit, Value: 32.0},
			want:   Degree{Type: Celsius, Value: 0.0},
		},
		{
			desc:   "Fahrenheit to Celsius (negative)",
			target: Celsius,
			from:   Degree{Type: Forenheit, Value: -40.0},
			want:   Degree{Type: Celsius, Value: -40.0},
		},
		{
			desc:   "Fahrenheit to Kelvin",
			target: Kelvin,
			from:   Degree{Type: Forenheit, Value: 32.0},
			want:   Degree{Type: Kelvin, Value: 273.15},
		},
		{
			desc:   "Fahrenheit to Kelvin (boiling point)",
			target: Kelvin,
			from:   Degree{Type: Forenheit, Value: 212.0},
			want:   Degree{Type: Kelvin, Value: 373.15},
		},
		{
			desc:    "wrong source type (Celsius)",
			target:  Forenheit,
			from:    Degree{Type: Celsius, Value: 0.0},
			wantErr: "expected Forenheit",
		},
		{
			desc:    "unsaported target",
			target:  0,
			from:    Degree{Type: Forenheit, Value: -1.0},
			wantErr: "unsupported target type",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := convertFromForenheit(tt.target, tt.from)
			checkErr(t, err, tt.wantErr)
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

// New test function for convertFromKelvin
func TestConvertFromKelvin(t *testing.T) {
	testCases := []struct {
		desc    string
		target  DegreeType
		from    Degree
		want    Degree
		wantErr string
	}{
		{
			desc:   "Kelvin to Kelvin",
			target: Kelvin,
			from:   Degree{Type: Kelvin, Value: 300.0},
			want:   Degree{Type: Kelvin, Value: 300.0},
		},
		{
			desc:   "Kelvin to Celsius",
			target: Celsius,
			from:   Degree{Type: Kelvin, Value: 273.15},
			want:   Degree{Type: Celsius, Value: 0.0},
		},
		{
			desc:   "Kelvin to Celsius (negative not physically possible but parseable)",
			target: Celsius,
			from:   Degree{Type: Kelvin, Value: 0.0},
			want:   Degree{Type: Celsius, Value: -273.15},
		},
		{
			desc:   "Kelvin to Fahrenheit",
			target: Forenheit,
			from:   Degree{Type: Kelvin, Value: 273.15},
			want:   Degree{Type: Forenheit, Value: 32.0},
		},
		{
			desc:   "Kelvin to Fahrenheit (boiling point)",
			target: Forenheit,
			from:   Degree{Type: Kelvin, Value: 373.15},
			want:   Degree{Type: Forenheit, Value: 212.0},
		},
		{
			desc:    "wrong source type (Celsius)",
			target:  Kelvin,
			from:    Degree{Type: Celsius, Value: 0.0},
			wantErr: "expected Kelvin",
		},
		{
			desc:    "unsaported target",
			target:  0,
			from:    Degree{Type: Kelvin, Value: -1.0},
			wantErr: "unsupported target type",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := convertFromKelvin(tt.target, tt.from)
			checkErr(t, err, tt.wantErr)
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestConvertorMap(t *testing.T) {
	testCases := []struct {
		desc     string
		fromType DegreeType
		fromVal  float64
		toType   DegreeType
		expected float64
	}{
		// Celsius conversions
		{desc: "Celsius to Fahrenheit", fromType: Celsius, fromVal: 0.0, toType: Forenheit, expected: 32.0},
		{desc: "Celsius to Kelvin", fromType: Celsius, fromVal: 100.0, toType: Kelvin, expected: 373.15},
		// Fahrenheit conversions
		{desc: "Fahrenheit to Celsius", fromType: Forenheit, fromVal: 32.0, toType: Celsius, expected: 0.0},
		{desc: "Fahrenheit to Kelvin", fromType: Forenheit, fromVal: 212.0, toType: Kelvin, expected: 373.15},
		// Kelvin conversions
		{desc: "Kelvin to Celsius", fromType: Kelvin, fromVal: 273.15, toType: Celsius, expected: 0.0},
		{desc: "Kelvin to Fahrenheit", fromType: Kelvin, fromVal: 373.15, toType: Forenheit, expected: 212.0},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			conv, ok := convertorMap[tt.fromType]
			if !ok {
				t.Fatalf("converter for %v not found", tt.fromType)
			}
			from := Degree{Type: tt.fromType, Value: tt.fromVal}
			result, err := conv(tt.toType, from)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			const eps = 1e-9
			if diff := result.Value - tt.expected; diff < -eps || diff > eps {
				t.Errorf("expected %f, got %f", tt.expected, result.Value)
			}
			if result.Type != tt.toType {
				t.Errorf("expected result type %v, got %v", tt.toType, result.Type)
			}
		})
	}
}

func TestRootCmdRun(t *testing.T) {
	tests := []struct {
		desc        string
		args        []string
		setFlags    func(cmd *cobra.Command)
		expectedOut string
		expectedErr string
	}{
		{
			desc: "success: default flags, Celsius to Fahrenheit",
			args: []string{"100C"},
			setFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("target", "F")
				cmd.Flags().Set("accuracy", "2")
				cmd.Flags().Set("exponent", "false")
			},
			expectedOut: "Result: 212.00°F\n",
		},
		{
			desc: "success: Fahrenheit to Celsius with accuracy 1",
			args: []string{"32F"},
			setFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("target", "C")
				cmd.Flags().Set("accuracy", "1")
				cmd.Flags().Set("exponent", "false")
			},
			expectedOut: "Result: 0.0°C\n",
		},
		{
			desc: "success: Kelvin to Celsius with exponent format",
			args: []string{"300K"},
			setFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("target", "C")
				cmd.Flags().Set("accuracy", "3")
				cmd.Flags().Set("exponent", "true")
			},
			expectedOut: "Result: 2.685E+01°C\n",
		},
		{
			desc: "success: from Kelvin to same type",
			args: []string{"273.15K"},
			setFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("target", "K")
				cmd.Flags().Set("accuracy", "2")
				cmd.Flags().Set("exponent", "false")
			},
			expectedOut: "Result: 273.15°K\n",
		},
		{
			desc: "success: spaces and degree symbol allowed",
			args: []string{"-40 °C"},
			setFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("target", "F")
				cmd.Flags().Set("accuracy", "0")
				cmd.Flags().Set("exponent", "false")
			},
			expectedOut: "Result: -40°F\n",
		},
		// Error cases
		{
			desc:        "error: no arguments",
			args:        []string{},
			setFlags:    func(cmd *cobra.Command) {},
			expectedErr: "no temperature provided",
		},
		{
			desc:        "error: invalid temperature format",
			args:        []string{"100X"},
			setFlags:    func(cmd *cobra.Command) {},
			expectedErr: "invalid temperature specifier",
		},
		{
			desc: "error: invalid target system",
			args: []string{"100C"},
			setFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("target", "X")
			},
			expectedErr: "invalid temperature specifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().StringP("target", "t", "C", "")
			cmd.Flags().UintP("accuracy", "a", 2, "")
			cmd.Flags().BoolP("exponent", "e", false, "")

			if tt.setFlags != nil {
				tt.setFlags(cmd)
			}

			getOutput, restore := captureStdout()
			defer restore()

			err := rootCmdRun(cmd, tt.args)

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectedErr)
				} else if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error %q, got %q", tt.expectedErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := getOutput()
			if output != tt.expectedOut {
				t.Errorf("output mismatch:\nexpected: %q\ngot:      %q", tt.expectedOut, output)
			}
		})
	}
}

func checkErr(t *testing.T, gotErr error, txt string) bool {
	t.Helper()

	if txt == "" && gotErr != nil {
		t.Fatalf("unexpected error: %v", gotErr)
		return false
	}
	if txt != "" && gotErr == nil {
		t.Fatalf("expected error that contains '%s', got nil", shortText(txt, -1))
		return false
	}
	if txt != "" && gotErr != nil {
		if !strings.Contains(gotErr.Error(), txt) {
			t.Errorf("error '%s' must contains '%s'", gotErr.Error(), shortText(txt, -1))
			return false
		}
	}

	return true
}

func shortText(txt string, n int) string {
	if n < 0 {
		n = 42
	}
	if len([]rune(txt)) <= n {
		return txt
	}

	short := txt[:n] + "..."
	return short
}

func captureStdout() (func() string, func()) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	restore := func() {
		w.Close()
		os.Stdout = old
	}

	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		io.Copy(&buf, r)
		close(done)
	}()

	// Return the output string and the restore function that also waits for copying.
	return func() string {
		restore()
		<-done
		return buf.String()
	}, restore
}
