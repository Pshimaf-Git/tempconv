# tempconv

A command-line tool for converting temperatures between Celsius, Fahrenheit, and Kelvin.

## Installation

```bash
go install github.com/Pshimaf-Git/tempconv@latest

Or clone and build:

```bash
git clone https://github.com/Pshimaf-Git/tempconv.git
cd tempconv
go build -o tempconv
```

Also do not forget add the path to binary file in your `$PATH` env variable


## Usage

```bash
tempconv <temperature> [flags]
```

## Examples

```bash
# Convert 100°F to Celsius (default target)
tempconv 100F
# Result: 37.78°C

# Convert 0°C to Fahrenheit
tempconv 0C -t F
# Result: 32.00°F

# Convert 300K to Kelvin (no change)
tempconv 300K -t K
# Result: 300.00°K

# Scientific notation with 4-digit accuracy
tempconv 100F -t C -e -a 4
# Result: 3.7778E+01°C

# Spaces and degree symbols are allowed
tempconv "212 °F" -t C
# Result: 100.00°C
```

## Supported Units

For now `tempconv` supports conversion between:
  - Celsius
  - Fahrenheit
  - Kelvin

Recognised suffixes(any case available, for example `c`, `°c` and `CeLSius` are valid):
  - `C`, `°C`, `Celsius`
  - `F`, `°F`, `Fahrenheit`
  - `K`, `°K`, `Kelvin`

## License

This project available by [MIT](./LICENSE)
