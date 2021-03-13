package cli

import (
	//"fmt"
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	//"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	//"github.com/cosmos/cosmos-sdk/codec"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cybercongress/go-cyber/types/query"
	"github.com/cybercongress/go-cyber/x/rank/types"
)

func GetQueryCmd() *cobra.Command {
	rankingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the rank module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	rankingQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryRank(),
		GetCmdQuerySearch(),
		GetCmdQueryBacklinks(),
		GetCmdQueryTop(),
		GetCmdQueryIsLinkExist(),
		GetCmdQueryIsAnyLinkExist(),
	)

	return rankingQueryCmd
}


func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current rank parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryRank() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "rank [cid]",
		Short: "Query the current rank of given CID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Rank(
				context.Background(),
				&types.QueryRankRequest{Cid: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQuerySearch() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "search [cid] [page] [limit]",
		Short: "Query search of given CID",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var page, limit uint32
			if len(args) == 3 {
				p, err := strconv.ParseUint(args[1], 10, 32)
				if err != nil {
					return err
				}
				page = uint32(p)
				l, err := strconv.ParseUint(args[2], 10, 32)
				if err != nil {
					return err
				}
				limit = uint32(l)
			} else {
				page = 0
				limit = 10
			}

			res, err := queryClient.Search(
				context.Background(),
				&types.QuerySearchRequest{Cid: args[0], Pagination: &query.PageRequest{Page: page, PerPage: limit}},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryBacklinks() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "backlinks [cid] [page] [limit]",
		Short: "Query backlinks of given CID",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var page, limit uint32
			if len(args) == 3 {
				p, err := strconv.ParseUint(args[1], 10, 32)
				if err != nil {
					return err
				}
				page = uint32(p)
				l, err := strconv.ParseUint(args[2], 10, 32)
				if err != nil {
					return err
				}
				limit = uint32(l)
			} else {
				page = 0
				limit = 10
			}

			res, err := queryClient.Backlinks(
				context.Background(),
				&types.QuerySearchRequest{Cid: args[0], Pagination: &query.PageRequest{Page: page, PerPage: limit}},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryTop() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Query top",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			var page, limit uint32
			if len(args) == 2 {
				p, err := strconv.ParseUint(args[1], 10, 32)
				if err != nil {
					return err
				}
				page = uint32(p)
				l, err := strconv.ParseUint(args[2], 10, 32)
				if err != nil {
					return err
				}
				limit = uint32(l)
			} else {
				page = 0
				limit = 10
			}

			res, err := queryClient.Top(
				context.Background(),
				&query.PageRequest{Page: page, PerPage: limit},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryIsLinkExist() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "is-exist [from] [to] [account]",
		Short: "Query is link exist between cids for given account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IsLinkExist(
				context.Background(),
				&types.QueryIsLinkExistRequest{args[0], args[1], args[2]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryIsAnyLinkExist() *cobra.Command{
	cmd := &cobra.Command{
		Use:   "is-exist-any [from] [to]",
		Short: "Query is any link exist between cids",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IsAnyLinkExist(
				context.Background(),
				&types.QueryIsAnyLinkExistRequest{args[0], args[1]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}