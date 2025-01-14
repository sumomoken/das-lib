package txbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"sync"
)

var log = mylog.NewLogger("txbuilder", mylog.LevelDebug)

type DasTxBuilder struct {
	*DasTxBuilderBase                                  // for base
	*DasTxBuilderTransaction                           // for tx
	DasMMJson                                          // for mmjson
	mapCellDep               map[string]*types.CellDep // for memory
	notCheckInputs           bool
	otherWitnesses           [][]byte
}

func NewDasTxBuilderBase(ctx context.Context, dasCore *core.DasCore, handle sign.HandleSignCkbMessage, serverArgs string) *DasTxBuilderBase {
	var base DasTxBuilderBase
	base.ctx = ctx
	base.dasCore = dasCore
	base.handleServerSign = handle
	base.serverArgs = serverArgs
	return &base
}

func NewDasTxBuilderFromBase(base *DasTxBuilderBase, tx *DasTxBuilderTransaction) *DasTxBuilder {
	var b DasTxBuilder
	b.DasTxBuilderBase = base
	b.DasTxBuilderTransaction = tx
	if tx == nil {
		b.DasTxBuilderTransaction = &DasTxBuilderTransaction{}
		b.MapInputsCell = make(map[string]*types.CellWithStatus)
	}
	b.mapCellDep = make(map[string]*types.CellDep)
	return &b
}

type DasTxBuilderBase struct {
	ctx              context.Context
	dasCore          *core.DasCore
	handleServerSign sign.HandleSignCkbMessage
	serverArgs       string
}

type DasTxBuilderTransaction struct {
	Transaction     *types.Transaction               `json:"transaction"`
	MapInputsCell   map[string]*types.CellWithStatus `json:"map_inputs_cell"`
	ServerSignGroup []int                            `json:"server_sign_group"`
}

type DasMMJson struct {
	account   string
	salePrice uint64
	offers    int // cancel offer count
}

type BuildTransactionParams struct {
	CellDeps       []*types.CellDep    `json:"cell_deps"`
	HeadDeps       []types.Hash        `json:"head_deps"`
	Inputs         []*types.CellInput  `json:"inputs"`
	Outputs        []*types.CellOutput `json:"outputs"`
	OutputsData    [][]byte            `json:"outputs_data"`
	Witnesses      [][]byte            `json:"witnesses"`
	OtherWitnesses [][]byte            `json:"other_witnesses"`
}

func (d *DasTxBuilder) BuildTransactionWithCheckInputs(p *BuildTransactionParams, notCheckInputs bool) error {
	d.notCheckInputs = notCheckInputs
	err := d.newTx()
	if err != nil {
		return fmt.Errorf("newBaseTx err: %s", err.Error())
	}

	err = d.addInputsForTx(p.Inputs)
	if err != nil {
		return fmt.Errorf("addInputsForBaseTx err: %s", err.Error())
	}

	err = d.addOutputsForTx(p.Outputs, p.OutputsData)
	if err != nil {
		return fmt.Errorf("addOutputsForBaseTx err: %s", err.Error())
	}

	d.Transaction.Witnesses = append(d.Transaction.Witnesses, p.Witnesses...)
	d.otherWitnesses = append(d.otherWitnesses, p.OtherWitnesses...)

	if err := d.addMapCellDepWitnessForBaseTx(p.CellDeps); err != nil {
		return fmt.Errorf("addMapCellDepWitnessForBaseTx err: %s", err.Error())
	}

	for _, v := range p.HeadDeps {
		d.Transaction.HeaderDeps = append(d.Transaction.HeaderDeps, v)
	}

	return nil
}

func (d *DasTxBuilder) BuildTransaction(p *BuildTransactionParams) error {
	return d.BuildTransactionWithCheckInputs(p, false)
}

func (d *DasTxBuilder) TxString() string {
	txStr, _ := rpc.TransactionString(d.Transaction)
	return txStr
}

func (d *DasTxBuilder) GetDasTxBuilderTransactionString() string {
	bys, err := json.Marshal(d.DasTxBuilderTransaction)
	if err != nil {
		return ""
	}
	return string(bys)
}

func GenerateDigestListFromTx(cli rpc.Client, txJson string, skipGroups []int) ([]SignData, error) {
	Tx, err := rpc.TransactionFromString(txJson)
	if err != nil {
		return nil, err
	}
	hash, _ := Tx.ComputeHash()
	fmt.Println(hash.Hex())

	var dasTxBuilderTransaction DasTxBuilderTransaction
	var txBuilder DasTxBuilder
	var netType common.DasNetType
	wgServer := sync.WaitGroup{}
	ctxServer := context.Background()
	blockInfo, err := cli.GetBlockchainInfo(ctxServer)
	if err != nil {
		return nil, err
	}

	if blockInfo.Chain == "ckb" {
		netType = common.DasNetTypeMainNet
	} else if blockInfo.Chain == "ckb_testnet" {
		netType = common.DasNetTypeTestnet2
	} else {
		netType = common.DasNetTypeTestnet3
	}

	dasTxBuilderTransaction.Transaction = Tx
	dasTxBuilderTransaction.MapInputsCell = make(map[string]*types.CellWithStatus)

	txBuilder.DasTxBuilderTransaction = &dasTxBuilderTransaction
	ops := []core.DasCoreOption{
		core.WithClient(cli),
		core.WithDasNetType(netType),
	}
	dasCore := core.NewDasCore(ctxServer, &wgServer, ops...)

	txBuilderBase := NewDasTxBuilderBase(ctxServer, dasCore, nil, "")
	txBuilder.DasTxBuilderBase = txBuilderBase
	digestList, err := txBuilder.GenerateDigestListFromTx(skipGroups)
	return digestList, nil
}
