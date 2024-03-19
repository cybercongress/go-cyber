package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"

	graphtypes "github.com/cybercongress/go-cyber/v3/x/graph/types"

	"github.com/cybercongress/go-cyber/v3/types/query"
	"github.com/cybercongress/go-cyber/v3/x/rank/types"
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
		GetCmdQueryNegentropyParticle(),
		GetCmdQueryNegentropy(),
		GetCmdQueryKarma(),
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

func GetCmdQueryRank() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rank [particle]",
		Short: "Query the current rank of given particle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err := cid.Decode(args[0]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

			res, err := queryClient.Rank(
				context.Background(),
				&types.QueryRankRequest{Particle: args[0]},
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

func GetCmdQuerySearch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [particle] [page] [limit]",
		Short: "Query search of given particle",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err := cid.Decode(args[0]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

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
				&types.QuerySearchRequest{Particle: args[0], Pagination: &query.PageRequest{Page: page, PerPage: limit}},
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

func GetCmdQueryBacklinks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlinks [particle] [page] [limit]",
		Short: "Query backlinks of given particle",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err := cid.Decode(args[0]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

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
				&types.QuerySearchRequest{Particle: args[0], Pagination: &query.PageRequest{Page: page, PerPage: limit}},
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

func GetCmdQueryTop() *cobra.Command {
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

func GetCmdQueryIsLinkExist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-exist [from] [to] [account]",
		Short: "Query is link exist between particles for given account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err := cid.Decode(args[0]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

			if _, err := cid.Decode(args[1]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

			address, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			res, err := queryClient.IsLinkExist(
				context.Background(),
				&types.QueryIsLinkExistRequest{From: args[0], To: args[1], Address: address.String()},
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

func GetCmdQueryIsAnyLinkExist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-exist-any [from] [to]",
		Short: "Query is any link exist between particles",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err := cid.Decode(args[0]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

			if _, err := cid.Decode(args[1]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

			res, err := queryClient.IsAnyLinkExist(
				context.Background(),
				&types.QueryIsAnyLinkExistRequest{From: args[0], To: args[1]},
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

func GetCmdQueryNegentropyParticle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "negentropy [particle]",
		Short: "Query the current negentropy of given particle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err := cid.Decode(args[0]); err != nil {
				return graphtypes.ErrInvalidParticle
			}

			res, err := queryClient.ParticleNegentropy(
				context.Background(),
				&types.QueryNegentropyPartilceRequest{Particle: args[0]},
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

func GetCmdQueryNegentropy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "negentropy",
		Short: "Query the current negentropy of whole graph",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Negentropy(
				context.Background(),
				&types.QueryNegentropyRequest{},
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

func GetCmdQueryKarma() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "karma [neuron]",
		Short: "Query the current karma of given neuron",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Karma(
				context.Background(),
				&types.QueryKarmaRequest{Neuron: address.String()},
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
