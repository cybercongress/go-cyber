package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/iavl"
	"github.com/spf13/cobra"
	goleveldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	dbm "github.com/tendermint/tm-db"
)

const (
	DefaultCacheSize int = 16000
)

var (
	DefaultHome = os.ExpandEnv("$HOME/") + ".cyber/data"
	rootCmd     = &cobra.Command{Use: "iavltool"}
	home        string
)

// TODO autoconf stores
var appStores = []string{
	"acc",
	"bank",
	"staking",
	"mint",
	"distribution",
	"slashing",
	"gov",
	"params",
	"ibc",
	"upgrade",
	"evidence",
	"transfer",
	"capability",
	"wasm",
	"graph",
	"bandwidth",
	"grid",
	"rank",
	"dmn",
	"resources",
	"liquidity",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&home, "home", DefaultHome, "path to cyber data")
	rootCmd.AddCommand(dataCmd)
	rootCmd.AddCommand(shapeCmd)
	rootCmd.AddCommand(versionsCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(pruneCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data [store] [version] [kv] [hash]",
	Short: "Print data of given stores at given block",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := OpenDB(home)
		if err != nil {
			fmt.Println("ERROR DB OPEN:", err)
		}

		stores := appStores
		version := int64(0)
		keysOpt, hashingOpt := false, false
		switch len(args) {
		case 4:
			hashingOpt, _ = strconv.ParseBool(args[3])
			fallthrough
		case 3:
			keysOpt, _ = strconv.ParseBool(args[2])
			fallthrough
		case 2:
			version, _ = strconv.ParseInt(args[1], 10, 64)
			fallthrough
		case 1:
			var a []string
			if args[0] != "all" {
				stores = append(a, args[0])
			}
		}

		for _, name := range stores {
			tree, err := ReadTree(db, version, name)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
				os.Exit(1)
			}
			if keysOpt {
				PrintKeys(tree, hashingOpt)
			}
			hash, err := tree.Hash()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Hash: %X\n", hash)
			fmt.Printf("Size: %X\n", tree.Size())
		}
	},
}

var shapeCmd = &cobra.Command{
	Use:   "shape [store] [version]",
	Short: "Print shape of given stores at given block",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := OpenDB(home)
		if err != nil {
			fmt.Println("ERROR DB OPEN:", err)
		}

		store := "rank"
		version := int64(0)
		switch len(args) {
		case 2:
			version, _ = strconv.ParseInt(args[1], 10, 64)
			fallthrough
		case 1:
			store = args[0]
		}

		tree, err := ReadTree(db, version, store)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
			os.Exit(1)
		}
		PrintShape(tree)
	},
}

var versionsCmd = &cobra.Command{
	Use:   "versions [store]",
	Short: "Print shape of given stores at given block",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := OpenDB(home)
		if err != nil {
			fmt.Println("ERROR DB OPEN:", err)
		}

		tree, err := ReadTree(db, 0, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
			os.Exit(1)
		}
		PrintVersions(tree)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [store] [from] [to]",
	Short: "Delete versions range for given stores",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := OpenDB(home)
		if err != nil {
			fmt.Println("ERROR DB OPEN:", err)
		}

		stores := appStores
		from, to := int64(0), int64(0)
		switch len(args) {
		case 3:
			to, _ = strconv.ParseInt(args[2], 10, 64)
			fallthrough
		case 2:
			from, _ = strconv.ParseInt(args[1], 10, 64)
			fallthrough
		case 1:
			var a []string
			if args[0] != "all" {
				stores = append(a, args[0])
			}
		}

		for _, name := range stores {
			tree, err := GetTree(db, name)
			fmt.Println("processing", name)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
				continue
			}
			err = tree.DeleteVersionsRange(from, to)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error deleting data: %s\n", err)
				continue
			}
		}
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats [store] [version]",
	Short: "Print shape of given stores at given block",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := OpenDB(home)
		if err != nil {
			fmt.Println("ERROR DB OPEN:", err)
		}

		PrintDBStats(db)
	},
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune leveldb",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		db, _ := goleveldb.OpenFile(home+"/application.db", nil)
		defer db.Close()
		_ = db.CompactRange(util.Range{Start: nil, Limit: nil})
	},
}

func OpenDB(dir string) (dbm.DB, error) {
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, dir)
	return db, err
}

func PrintDBStats(db dbm.DB) {
	count := 0
	prefix := map[string]int{}
	iter, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	for ; iter.Valid(); iter.Next() {
		key := string(iter.Key()[:1])
		prefix[key]++
		count++
	}
	iter.Close()
	fmt.Printf("DB contains %d entries\n", count)
	for k, v := range prefix {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

// ReadTree loads an iavl tree from the directory
// If version is 0, load latest, otherwise, load named version
func ReadTree(db dbm.DB, version int64, name string) (*iavl.MutableTree, error) {
	fmt.Println("--------------[", name, "]--------------")
	tree, err := iavl.NewMutableTree(dbm.NewPrefixDB(db, []byte("s/k:"+name+"/")), DefaultCacheSize, false)
	if err != nil {
		return nil, err
	}
	ver, err := tree.LoadVersion(version)
	fmt.Printf("Got version: %d\n", ver)
	return tree, err
}

func GetTree(db dbm.DB, name string) (*iavl.MutableTree, error) {
	tree, err := iavl.NewMutableTree(dbm.NewPrefixDB(db, []byte("s/k:"+name+"/")), DefaultCacheSize, false)
	if err != nil {
		return nil, err
	}
	_, err = tree.LoadVersion(int64(0))
	if err != nil {
		return nil, err
	}
	return tree, err
}

func PrintKeys(tree *iavl.MutableTree, hashing bool) {
	fmt.Println("Printing all keys with hashed values (to detect diff)")
	_, err := tree.Iterate(func(key, value []byte) bool {
		if hashing {
			printKey := parseWeaveKey(key)
			digest := sha256.Sum256(value)
			fmt.Printf("  %s\n    %X\n", printKey, digest)
		} else {
			fmt.Printf("  %s\n    %X\n", key, value)
		}
		return false
	})
	if err != nil {
		panic(err)
	}
}

// parseWeaveKey assumes a separating : where all in front should be ascii,
// and all afterwards may be ascii or binary
func parseWeaveKey(key []byte) string {
	cut := bytes.IndexRune(key, ':')
	if cut == -1 {
		return encodeID(key)
	}
	prefix := key[:cut]
	id := key[cut+1:]
	return fmt.Sprintf("%s:%s", encodeID(prefix), encodeID(id))
}

// casts to a string if it is printable ascii, hex-encodes otherwise
func encodeID(id []byte) string {
	for _, b := range id {
		if b < 0x20 || b >= 0x80 {
			return strings.ToUpper(hex.EncodeToString(id))
		}
	}
	return string(id)
}

func PrintShape(tree *iavl.MutableTree) {
	// shape := tree.RenderShape("  ", nil)
	shape, _ := tree.RenderShape("  ", nodeEncoder)
	fmt.Println(strings.Join(shape, "\n"))
}

func nodeEncoder(id []byte, depth int, isLeaf bool) string {
	prefix := fmt.Sprintf("-%d ", depth)
	if isLeaf {
		prefix = fmt.Sprintf("*%d ", depth)
	}
	if len(id) == 0 {
		return fmt.Sprintf("%s<nil>", prefix)
	}
	return fmt.Sprintf("%s%s", prefix, parseWeaveKey(id))
}

func PrintVersions(tree *iavl.MutableTree) {
	versions := tree.AvailableVersions()
	fmt.Println("Available versions:")
	for _, v := range versions {
		fmt.Printf("  %d\n", v)
	}
}
