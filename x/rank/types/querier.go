package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParameters  	= "params"
	QueryRank		 	= "rank"
	QuerySearch      	= "search"
	QueryBacklinks   	= "backlinks"
	QueryTop		 	= "top"
	QueryIsLinkExist 	= "is_link_exist"
	QueryIsAnyLinkExist = "is_any_link_exist"
	QueryKarma 	 		= "karma"
	QueryEntropy 	 	= "entropy"
	QueryLuminosity 	= "luminosity"
	QueryKarmas 	 	= "karmas"
)


type QueryRankParams struct {
	Cid string
}

func NewQueryRankParams(cid string) QueryRankParams {
	return QueryRankParams{cid}
}

type QuerySearchParams struct {
	Cid string
	Page, PerPage uint32
}

func NewQuerySearchParams(cid string, page, perPage uint32) QuerySearchParams {
	return QuerySearchParams{cid, page, perPage}
}

type QueryTopParams struct {
	Page, PerPage uint32
}

func NewQueryTopParams(page, perPage uint32) QueryTopParams {
	return QueryTopParams{page, perPage}
}


type QueryIsLinkExistParams struct {
	From, To string
	Address sdk.AccAddress
}

func NewQueryIsLinkExistParams(from, to string, account sdk.AccAddress) QueryIsLinkExistParams {
	return QueryIsLinkExistParams{from, to, account}
}

type QueryIsAnyLinkExistParams struct {
	From, To string
}

func NewQueryIsAnyLinkExistParams(from, to string) QueryIsAnyLinkExistParams {
	return QueryIsAnyLinkExistParams{from, to}
}

type QueryEntropyParams struct {
	Cid string
}

func NewQueryEntropyParams(cid string) QueryEntropyParams {
	return QueryEntropyParams{cid}
}

type QueryLuminosityParams struct {
	Cid string
}

func NewQueryLuminosityParams(cid string) QueryLuminosityParams {
	return QueryLuminosityParams{cid}
}

type QueryKarmaParams struct {
	Address sdk.AccAddress
}

func NewQueryKarmaParams(account sdk.AccAddress) QueryKarmaParams {
	return QueryKarmaParams{account}
}