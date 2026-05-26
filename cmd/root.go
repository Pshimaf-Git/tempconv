// Package cmd
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type DegreeType int

func (t DegreeType) String() string {
	switch t {
	case Celsius:
		return "C"
	case Forenheit:
		return "F"
	case Kelvin:
		return "K"
	}
	return "unknown"
}

type Degree struct {
	Type  DegreeType
	Value float64
}

const (
	_ DegreeType = iota
	Celsius
	Forenheit
	Kelvin
)

type convertorFn func(target DegreeType, from Degree) (Degree, error)

// convertFromCelsius converts from Celsius to any target type
func convertFromCelsius(target DegreeType, from Degree) (Degree, error) {
	if from.Type != Celsius {
		return Degree{}, fmt.Errorf("convertFromCelsius: expected Celsius, got %v", from.Type)
	}
	switch target {
	case Celsius:
		return from, nil
	case Forenheit:
		return Degree{Type: Forenheit, Value: from.Value*1.8 + 32}, nil
	case Kelvin:
		return Degree{Type: Kelvin, Value: from.Value + 273.15}, nil
	default:
		return Degree{}, fmt.Errorf("unsupported target type: %v", target)
	}
}

// convertFromForenheit converts from Fahrenheit to any target type
func convertFromForenheit(target DegreeType, from Degree) (Degree, error) {
	if from.Type != Forenheit {
		return Degree{}, fmt.Errorf("convertFromForenheit: expected Forenheit, got %v", from.Type)
	}
	switch target {
	case Forenheit:
		return from, nil
	case Celsius:
		return Degree{Type: Celsius, Value: (from.Value - 32) / 1.8}, nil
	case Kelvin:
		celsius := (from.Value - 32) / 1.8
		return Degree{Type: Kelvin, Value: celsius + 273.15}, nil
	default:
		return Degree{}, fmt.Errorf("unsupported target type: %v", target)
	}
}

// convertFromKelvin converts from Kelvin to any target type
func convertFromKelvin(target DegreeType, from Degree) (Degree, error) {
	if from.Type != Kelvin {
		return Degree{}, fmt.Errorf("convertFromKelvin: expected Kelvin, got %v", from.Type)
	}
	switch target {
	case Kelvin:
		return from, nil
	case Celsius:
		return Degree{Type: Celsius, Value: from.Value - 273.15}, nil
	case Forenheit:
		celsius := from.Value - 273.15
		return Degree{Type: Forenheit, Value: celsius*1.8 + 32}, nil
	default:
		return Degree{}, fmt.Errorf("unsupported target type: %v", target)
	}
}

var convertorMap = map[DegreeType]convertorFn{
	Celsius:   convertFromCelsius,
	Forenheit: convertFromForenheit,
	Kelvin:    convertFromKelvin,
}

var rootCmd = &cobra.Command{
	Use:   "tempconv",
	Short: "Application for converting temperatures from different systems",
	Long:  `Example: tempconv 100F -t C`,
	RunE:  rootCmdRun,
}

// ParseTempratureSys extracts degree type from a string like "C", "°F", "K"
func ParseTempratureSys(s string) (DegreeType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return 0, fmt.Errorf("ParseTempratureSys: empty temperature system")
	}

	switch {
	case strings.HasSuffix(s, "CELSIUS"),
		strings.HasSuffix(s, "C"), strings.HasSuffix(s, "°C"):
		return Celsius, nil
	case strings.HasSuffix(s, "FORENHEIT"), strings.HasSuffix(s, "F"), strings.HasSuffix(s, "°F"):
		return Forenheit, nil
	case strings.HasSuffix(s, "KELVIN"), strings.HasSuffix(s, "K"), strings.HasSuffix(s, "°K"):
		return Kelvin, nil
	default:
		return 0, fmt.Errorf("ParseTempratureSys: invalid temperature specifier '%s'", s)
	}
}

// ParseTemprature parses a string like "100C", "32.5°F", "273.15K" into a Degree
func ParseTemprature(s string) (Degree, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return Degree{}, fmt.Errorf("empty string")
	}

	// Determine type
	dType, err := ParseTempratureSys(s)
	if err != nil {
		return Degree{}, err
	}

	replacer := strings.NewReplacer("C", "", "F", "", "K", "", " ", "", "°", "")

	s = replacer.Replace(s)

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Degree{}, fmt.Errorf("failed to parse number '%s': %w", s, err)
	}

	return Degree{Type: dType, Value: val}, nil
}

func rootCmdRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no temperature provided")
	}
	temperatureStr := args[0]

	temperature, err := ParseTemprature(temperatureStr)
	if err != nil {
		return err
	}

	targetSysStr, err := cmd.Flags().GetString("target")
	if err != nil {
		return fmt.Errorf("failed to get target flag: %w", err)
	}

	targetSys, err := ParseTempratureSys(targetSysStr)
	if err != nil {
		return err
	}

	convertor, exists := convertorMap[temperature.Type]
	if !exists {
		return fmt.Errorf("conversion from %v not supported", temperature.Type)
	}

	result, err := convertor(targetSys, temperature)
	if err != nil {
		return err
	}

	accuracy, err := cmd.Flags().GetUint("accuracy")
	if err != nil {
		return fmt.Errorf("failed to get accuracy flag: %w", err)
	}
	useE, err := cmd.Flags().GetBool("exponent")
	if err != nil {
		return fmt.Errorf("failed to get exponent flag: %w", err)
	}

	var format byte = 'f'
	if useE {
		format = 'E'
	}

	val := strconv.FormatFloat(result.Value, format, int(accuracy), 64)

	fmt.Printf("Result: %s°%s\n", val, result.Type.String())
	return nil
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("target", "t", "C", "Target temperature system (C, F, K)")
	rootCmd.Flags().BoolP("exponent", "e", false, "Use exponent format for print result")
	rootCmd.Flags().UintP("accuracy", "a", 2, "The number of digits after the decimal point in a floating-point number.")
}
