package types


const (
	QueryLinks				= "links"
	QueryLinksAmount		= "amount"
	QueryInLinks			= "in"
	QueryOutLinks			= "out"
)

type QueryLinksParams struct {
	Cid Cid
}

func NewQueryLinksPrams(cid Cid) QueryLinksParams {
	return QueryLinksParams{
		Cid: cid,
	}
}

type ResultLinks struct {
	Cids []Cid `json:"cids" yaml:"cids"`
}

func NewResultLinks(cids []Cid) ResultLinks {
	return ResultLinks{
		Cids:cids,
	}
}
