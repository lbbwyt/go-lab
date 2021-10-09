package goinception

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
)

func (self *GoInception) parseSQLs(blob string) ([]string, []sqlparser.Statement, error) {
	queries, err := sqlparser.SplitStatementToPieces(blob)
	if err != nil {
		return nil, nil, err
	}
	stmts := make([]sqlparser.Statement, 0, len(queries))
	for _, query := range queries {
		stmt, err := sqlparser.Parse(query)
		if err != nil {
			if len(queries) == 1 {
				return nil, nil, fmt.Errorf("parseSQLs err : %s", blob)
			} else {
				return nil, nil, fmt.Errorf("parseSQLs err : %s", blob)
			}
		}
		stmts = append(stmts, stmt)
	}
	return queries, stmts, nil
}
