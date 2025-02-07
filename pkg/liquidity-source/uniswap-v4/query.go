package uniswapv4

import (
	"bytes"
	"text/template"
)

type PoolsListQueryParams struct {
	LastCreatedAtTimestamp int
	First                  int
	Skip                   int
}

type PoolTicksQueryParams struct {
	AllowSubgraphError bool
	PoolAddress        string
	LastTickIdx        string
}

func getPoolsListQuery(lastCreatedAtTimestamp int, first int) string {
	var tpl bytes.Buffer
	td := PoolsListQueryParams{
		lastCreatedAtTimestamp,
		first,
		0,
	}

	// Add subgraphError: allow
	t, err := template.New("poolsListQuery").Parse(`{
		pools(
			where: {
				blockTimestamp_gte: {{ .LastCreatedAtTimestamp }}
			},
			first: {{ .First }},
			skip: {{ .Skip }},
			orderBy: blockTimestamp,
			orderDirection: asc
		) {
			id
			poolId
			currency0
			currency1
			fee
			tickSpacing
			hooks
			blockTimestamp
		}
	}`)

	if err != nil {
		panic(err)
	}

	err = t.Execute(&tpl, td)

	if err != nil {
		panic(err)
	}

	return tpl.String()
}
