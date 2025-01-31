package txbuilder

import (
	"encoding/binary"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sort"
	"strings"
)

type SignData struct {
	SignType common.DasAlgorithmId `json:"sign_type"`
	SignMsg  string                `json:"sign_msg"`
}

func (d *DasTxBuilder) GenerateMultiSignDigest(group []int, firstN uint8, signatures [][]byte, sortArgsList [][]byte) ([]byte, error) {
	if len(group) == 0 {
		return nil, fmt.Errorf("group is nil")
	}

	wa := GenerateMultiSignWitnessArgs(firstN, signatures, sortArgsList)
	data, err := wa.Serialize()
	if err != nil {
		return nil, err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := d.Transaction.ComputeHash()
	if err != nil {
		return nil, err
	}
	message := append(hash.Bytes(), length...)
	message = append(message, data...)

	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
		}
	}

	// hash witnesses which do not in any input group
	for _, wit := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}

	message, err = blake2b.Blake256(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (d *DasTxBuilder) GenerateDigestListFromTx(skipGroups []int) ([]SignData, error) {
	skipGroups = append(skipGroups, d.ServerSignGroup...)
	groups, err := d.getGroupsFromTx()
	if err != nil {
		return nil, fmt.Errorf("getGroupsFromTx err: %s", err.Error())
	}
	log.Info("groups:", len(groups), groups[0])
	var digestList []SignData
	for _, group := range groups {
		if digest, err := d.generateDigestByGroup(group, skipGroups); err != nil {
			return nil, fmt.Errorf("generateDigestByGroup err: %s", err.Error())
		} else {
			digestList = append(digestList, digest)
		}
	}
	return digestList, nil
}

func (d *DasTxBuilder) getGroupsFromTx() ([][]int, error) {
	var tmpMapForGroup = make(map[string][]int)
	var sortList []string
	for i, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return nil, fmt.Errorf("getInputCell err: %s", err.Error())
		}

		cellHash, err := item.Cell.Output.Lock.Hash()
		if err != nil {
			return nil, fmt.Errorf("inputs lock to hash err: %s", err.Error())
		}
		indexList, okTmp := tmpMapForGroup[cellHash.String()]
		if !okTmp {
			sortList = append(sortList, cellHash.String())
		}
		indexList = append(indexList, i)
		tmpMapForGroup[cellHash.String()] = indexList
	}
	sort.Strings(sortList)
	var list [][]int
	for _, v := range sortList {
		item, _ := tmpMapForGroup[v]
		list = append(list, item)
	}
	return list, nil
}

func (d *DasTxBuilder) generateDigestByGroup(group []int, skipGroups []int) (SignData, error) {
	var signData = SignData{}
	if group == nil || len(group) < 1 {
		return signData, fmt.Errorf("invalid param")
	}
	// check AlgorithmId
	item, err := d.getInputCell(d.Transaction.Inputs[group[0]].PreviousOutput)
	if err != nil {
		return signData, fmt.Errorf("getInputCell err: %s", err.Error())
	}

	daf := core.DasAddressFormat{DasNetType: d.dasCore.NetType()}
	ownerHex, managerHex, _ := daf.ArgsToHex(item.Cell.Output.Lock.Args)
	ownerAlgorithmId, managerAlgorithmId := ownerHex.DasAlgorithmId, managerHex.DasAlgorithmId

	signData.SignType = ownerAlgorithmId

	actionBuilder, err := witness.ActionDataBuilderFromTx(d.Transaction)
	has712, action := false, ""
	if err != nil {
		if err != witness.ErrNotExistActionData {
			return signData, fmt.Errorf("witness.ActionDataBuilderFromTx err: %s", err.Error())
		}
	} else {
		action = actionBuilder.Action
		switch actionBuilder.Action {
		case common.DasActionEditRecords:
			signData.SignType = managerAlgorithmId
		case common.DasActionEnableSubAccount, common.DasActionCreateSubAccount,
			common.DasActionConfigSubAccountCustomScript:
			if signData.SignType == common.DasAlgorithmIdEth712 {
				signData.SignType = common.DasAlgorithmIdEth
			}
		}
		// 712
		switch actionBuilder.Action {
		case common.DasActionEditManager, common.DasActionEditRecords,
			common.DasActionTransferAccount, common.DasActionTransfer,
			common.DasActionWithdrawFromWallet, common.DasActionStartAccountSale,
			common.DasActionEditAccountSale, common.DasActionCancelAccountSale,
			common.DasActionBuyAccount, common.DasActionDeclareReverseRecord,
			common.DasActionRedeclareReverseRecord, common.DasActionRetractReverseRecord,
			common.DasActionMakeOffer, common.DasActionEditOffer, common.DasActionCancelOffer,
			common.DasActionAcceptOffer, common.DasActionLockAccountForCrossChain:
			has712 = true
		}
	}

	// gen digest
	log.Warn("generateDigestByGroup:", len(group), group, action, has712)
	digest := ""
	emptyWitnessArg := types.WitnessArgs{
		Lock:       make([]byte, 65),
		InputType:  nil,
		OutputType: nil,
	}
	if signData.SignType == common.DasAlgorithmIdDogeChain {
		emptyWitnessArg.Lock = make([]byte, 66)
	} else if signData.SignType == common.DasAlgorithmIdEth712 && has712 {
		emptyWitnessArg.Lock = make([]byte, 105)
	}
	data, err := emptyWitnessArg.Serialize()
	if err != nil {
		return signData, err
	}

	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := d.Transaction.ComputeHash()
	if err != nil {
		return signData, err
	}
	//fmt.Println("tx_hash:", hash.Hex())

	message := append(hash.Bytes(), length...)
	message = append(message, data...)
	//fmt.Println("init witness:", common.Bytes2Hex(message))
	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
			//fmt.Println("add group other witness:", common.Bytes2Hex(message))
		}
	}
	//fmt.Println("add group other witness:", common.Bytes2Hex(message))
	// hash witnesses which do not in any input group
	for _, wit := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}
	//fmt.Println("add other witness:", common.Bytes2Hex(message))

	message, err = blake2b.Blake256(message)
	if err != nil {
		return signData, err
	}
	digest = common.Bytes2Hex(message)
	signData.SignMsg = digest

	// fix tron digest
	if signData.SignType == common.DasAlgorithmIdTron {
		signData.SignMsg += "04"
	} else if signData.SignType == common.DasAlgorithmIdDogeChain {
		signData.SignMsg = strings.TrimPrefix(signData.SignMsg, common.HexPreFix)
	}
	log.Info("digest:", signData.SignMsg)

	// skip useless signature
	if len(skipGroups) != 0 {
		skip := false
		for i := range group {
			for j := range skipGroups {
				if group[i] == skipGroups[j] {
					skip = true
					break
				}
			}
			if skip {
				break
			}
		}
		if skip {
			signData.SignMsg = ""
		}
	}
	return signData, nil
}

func (d *DasTxBuilder) getInputCell(o *types.OutPoint) (*types.CellWithStatus, error) {
	if o == nil {
		return nil, fmt.Errorf("OutPoint is nil")
	}
	key := fmt.Sprintf("%s-%d", o.TxHash.Hex(), o.Index)
	if item, ok := d.MapInputsCell[key]; ok {
		if item.Cell != nil && item.Cell.Output != nil && item.Cell.Output.Lock != nil {
			return item, nil
		}
	}
	if cell, err := d.dasCore.Client().GetLiveCell(d.ctx, o, true); err != nil {
		return nil, fmt.Errorf("GetLiveCell err: %s", err.Error())
	} else if cell.Cell.Output.Lock != nil {
		d.MapInputsCell[key] = cell
		return cell, nil
	} else {
		log.Warn("GetLiveCell:", key, cell.Status)
		if !d.notCheckInputs {
			return nil, fmt.Errorf("cell [%s] not live", key)
		} else {
			d.MapInputsCell[key] = cell
			return cell, nil
		}
	}
}
