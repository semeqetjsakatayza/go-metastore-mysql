# Meta Store with MySQL-Backend

This package implements meta store support functions with table in MySQL database as backend.

# Development Notes

## Generate SQL statement code

```sh
go-literal-code-gen -in sqlstmt.md -out sqlstmt.go
go-literal-code-gen -in sqlstruct.md -out sqlstruct.go -sqlschema
```
