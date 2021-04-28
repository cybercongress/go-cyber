package types

const (
	QueryLinks			= "links"
	QueryLinksAmount	= "amount_links"
	QueryCidsAmount		= "amount_cids"
	QueryGraphStats		= "graph_stats"
	QueryInLinks		= "in"
	QueryOutLinks		= "out"
)

type QueryLinksParams struct {
	Cid   string
}

func NewQueryLinksParams(cid string) QueryLinksParams {
	return QueryLinksParams{cid}
}