package spansql

import (
	"strings"
	"log"
)

func gologoo__IsKeyword_f9d5d723d19d971c78d4758fe1e3c508(id string) bool {
	return keywords[strings.ToUpper(id)]
}

var keywords = map[string]bool{"ALL": true, "AND": true, "ANY": true, "ARRAY": true, "AS": true, "ASC": true, "ASSERT_ROWS_MODIFIED": true, "AT": true, "BETWEEN": true, "BY": true, "CASE": true, "CAST": true, "COLLATE": true, "CONTAINS": true, "CREATE": true, "CROSS": true, "CUBE": true, "CURRENT": true, "DEFAULT": true, "DEFINE": true, "DESC": true, "DISTINCT": true, "ELSE": true, "END": true, "ENUM": true, "ESCAPE": true, "EXCEPT": true, "EXCLUDE": true, "EXISTS": true, "EXTRACT": true, "FALSE": true, "FETCH": true, "FOLLOWING": true, "FOR": true, "FROM": true, "FULL": true, "GROUP": true, "GROUPING": true, "GROUPS": true, "HASH": true, "HAVING": true, "IF": true, "IGNORE": true, "IN": true, "INNER": true, "INTERSECT": true, "INTERVAL": true, "INTO": true, "IS": true, "JOIN": true, "LATERAL": true, "LEFT": true, "LIKE": true, "LIMIT": true, "LOOKUP": true, "MERGE": true, "NATURAL": true, "NEW": true, "NO": true, "NOT": true, "NULL": true, "NULLS": true, "OF": true, "ON": true, "OR": true, "ORDER": true, "OUTER": true, "OVER": true, "PARTITION": true, "PRECEDING": true, "PROTO": true, "RANGE": true, "RECURSIVE": true, "RESPECT": true, "RIGHT": true, "ROLLUP": true, "ROWS": true, "SELECT": true, "SET": true, "SOME": true, "STRUCT": true, "TABLESAMPLE": true, "THEN": true, "TO": true, "TREAT": true, "TRUE": true, "UNBOUNDED": true, "UNION": true, "UNNEST": true, "USING": true, "WHEN": true, "WHERE": true, "WINDOW": true, "WITH": true, "WITHIN": true}
var funcs = make(map[string]bool)
var funcArgParsers = make(map[string]func(*parser) (Expr, *parseError))

func gologoo__init_f9d5d723d19d971c78d4758fe1e3c508() {
	for _, f := range allFuncs {
		funcs[f] = true
	}
	funcArgParsers["CAST"] = typedArgParser
	funcArgParsers["SAFE_CAST"] = typedArgParser
	funcArgParsers["EXTRACT"] = extractArgParser
}

var allFuncs = []string{"ANY_VALUE", "ARRAY_AGG", "AVG", "BIT_XOR", "COUNT", "MAX", "MIN", "SUM", "CAST", "SAFE_CAST", "ABS", "MOD", "FARM_FINGERPRINT", "SHA1", "SHA256", "SHA512", "BYTE_LENGTH", "CHAR_LENGTH", "CHARACTER_LENGTH", "CODE_POINTS_TO_BYTES", "CODE_POINTS_TO_STRING", "CONCAT", "ENDS_WITH", "FORMAT", "FROM_BASE32", "FROM_BASE64", "FROM_HEX", "LENGTH", "LOWER", "LPAD", "LTRIM", "REGEXP_CONTAINS", "REGEXP_EXTRACT", "REGEXP_EXTRACT_ALL", "REGEXP_REPLACE", "REPEAT", "REPLACE", "REVERSE", "RPAD", "RTRIM", "SAFE_CONVERT_BYTES_TO_STRING", "SPLIT", "STARTS_WITH", "STRPOS", "SUBSTR", "TO_BASE32", "TO_BASE64", "TO_CODE_POINTS", "TO_HEX", "TRIM", "UPPER", "ARRAY", "ARRAY_CONCAT", "ARRAY_LENGTH", "ARRAY_TO_STRING", "GENERATE_ARRAY", "GENERATE_DATE_ARRAY", "OFFSET", "ORDINAL", "ARRAY_REVERSE", "ARRAY_IS_DISTINCT", "SAFE_OFFSET", "SAFE_ORDINAL", "CURRENT_DATE", "EXTRACT", "DATE", "DATE_ADD", "DATE_SUB", "DATE_DIFF", "DATE_TRUNC", "DATE_FROM_UNIX_DATE", "FORMAT_DATE", "PARSE_DATE", "UNIX_DATE", "CURRENT_TIMESTAMP", "STRING", "TIMESTAMP", "TIMESTAMP_ADD", "TIMESTAMP_SUB", "TIMESTAMP_DIFF", "TIMESTAMP_TRUNC", "FORMAT_TIMESTAMP", "PARSE_TIMESTAMP", "TIMESTAMP_SECONDS", "TIMESTAMP_MILLIS", "TIMESTAMP_MICROS", "UNIX_SECONDS", "UNIX_MILLIS", "UNIX_MICROS", "PENDING_COMMIT_TIMESTAMP", "JSON_VALUE"}

func IsKeyword(id string) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsKeyword_f9d5d723d19d971c78d4758fe1e3c508")
	log.Printf("Input : %v\n", id)
	r0 := gologoo__IsKeyword_f9d5d723d19d971c78d4758fe1e3c508(id)
	log.Printf("Output: %v\n", r0)
	return r0
}
func init() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__init_f9d5d723d19d971c78d4758fe1e3c508")
	log.Printf("Input : (none)\n")
	gologoo__init_f9d5d723d19d971c78d4758fe1e3c508()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
