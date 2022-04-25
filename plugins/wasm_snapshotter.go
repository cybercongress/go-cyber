package plugins

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/snapshots/types"
	snapshot "github.com/cosmos/cosmos-sdk/snapshots/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protoio "github.com/gogo/protobuf/io"
)

type WasmSnapshotter struct {
	wasmDirectory string
}

func NewWasmSnapshotter(wasmDirectory string) *WasmSnapshotter {
	return &WasmSnapshotter{
		wasmDirectory,
	}
}

func (ws *WasmSnapshotter) SnapshotName() string {
	return "WASM Files Snapshot"
}

func (ws *WasmSnapshotter) SnapshotFormat() uint32 {
	return 1
}

func (ws *WasmSnapshotter) SupportedFormats() []uint32 {
	return []uint32{1}
}

var wasmFileNameRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)

func (ws *WasmSnapshotter) Snapshot(height uint64, protoWriter protoio.Writer) error {
	wasmFiles, err := ioutil.ReadDir(ws.wasmDirectory)
	if err != nil {
		return err
	}

	// In case snapshotting needs to be deterministic
	sort.SliceStable(wasmFiles, func(i, j int) bool {
		return strings.Compare(wasmFiles[i].Name(), wasmFiles[j].Name()) < 0
	})

	for _, wasmFile := range wasmFiles {
		if !wasmFileNameRegex.MatchString(wasmFile.Name()) {
			continue
		}

		wasmFilePath := path.Join(ws.wasmDirectory, wasmFile.Name())
		wasmBytes, err := ioutil.ReadFile(wasmFilePath)
		if err != nil {
			return err
		}

		// snapshotItem is 64 bytes of the file name, then the actual WASM bytes
		snapshotItem := append([]byte(wasmFile.Name()), wasmBytes...)

		snapshot.WriteExtensionItem(protoWriter, snapshotItem)
	}

	return nil
}

func (ws *WasmSnapshotter) Restore(
	height uint64, format uint32, protoReader protoio.Reader,
) (snapshot.SnapshotItem, error) {
	if format != 1 {
		return snapshot.SnapshotItem{}, types.ErrUnknownFormat
	}

	err := os.MkdirAll(ws.wasmDirectory, os.ModePerm)
	if err != nil {
		return snapshot.SnapshotItem{}, sdkerrors.Wrapf(err, "failed to create directory '%s'", ws.wasmDirectory)
	}

	for {
		item := &snapshot.SnapshotItem{}
		err = protoReader.ReadMsg(item)
		if err == io.EOF {
			break
		} else if err != nil {
			return snapshot.SnapshotItem{}, sdkerrors.Wrap(err, "invalid protobuf message")
		}

		payload := item.GetExtensionPayload()
		if payload == nil {
			return snapshot.SnapshotItem{}, sdkerrors.Wrap(err, "invalid protobuf message")
		}

		// snapshotItem is 64 bytes of the file name, then the actual WASM bytes
		if len(payload.Payload) < 64 {
			return snapshot.SnapshotItem{}, sdkerrors.Wrapf(err, "wasm snapshot must be at least 64 bytes, got %v bytes", len(payload.Payload))
		}

		wasmFileName := string(payload.Payload[0:64])
		wasmBytes := payload.Payload[64:]

		wasmFilePath := path.Join(ws.wasmDirectory, wasmFileName)

		err = ioutil.WriteFile(wasmFilePath, wasmBytes, 0664 /* -rw-rw-r-- */)
		if err != nil {
			return snapshot.SnapshotItem{}, sdkerrors.Wrapf(err, "failed to write wasm file '%v' to disk", wasmFilePath)
		}
	}

	return snapshot.SnapshotItem{}, nil
}