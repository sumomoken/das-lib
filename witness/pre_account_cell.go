package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type PreAccountCellDataBuilder struct {
	Index           uint32
	Version         uint32
	AccountChars    *molecule.AccountChars
	Account         string
	RefundLock      *molecule.Script
	OwnerLockArgs   string
	InviterId       string
	InviterLock     *molecule.Script
	ChannelLock     *molecule.Script
	Price           *molecule.PriceConfig
	Quote           *molecule.Uint64
	InvitedDiscount *molecule.Uint32
	CreatedAt       *molecule.Uint64
	InitialRecords  *molecule.Records

	preAccountCellDataV1 *molecule.PreAccountCellDataV1
	preAccountCellData   *molecule.PreAccountCellData
	dataEntityOpt        *molecule.DataEntityOpt
}

func PreAccountCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*PreAccountCellDataBuilder, error) {
	respMap, err := PreAccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist pre account cell")
}
func PreAccountIdCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*PreAccountCellDataBuilder, error) {
	respMap, err := PreAccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*PreAccountCellDataBuilder)
	for k, v := range respMap {
		k1 := common.Bytes2Hex(common.GetAccountIdByAccount(k))
		retMap[k1] = v
	}
	return retMap, nil
}
func PreAccountCellDataBuilderMapFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*PreAccountCellDataBuilder, error) {
	var respMap = make(map[string]*PreAccountCellDataBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypePreAccountCell:
			var resp PreAccountCellDataBuilder
			dataEntityOpt, dataEntity, err := getDataEntityOpt(dataBys, dataType)
			if err != nil {
				return false, fmt.Errorf("getDataEntityOpt err: %s", err.Error())
			}
			resp.dataEntityOpt = dataEntityOpt

			version, err := molecule.Bytes2GoU32(dataEntity.Version().RawData())
			if err != nil {
				return false, fmt.Errorf("get version err: %s", err.Error())
			}
			resp.Version = version

			index, err := molecule.Bytes2GoU32(dataEntity.Index().RawData())
			if err != nil {
				return false, fmt.Errorf("get index err: %s", err.Error())
			}
			resp.Index = index

			switch version {
			case common.GoDataEntityVersion1:
				if err := resp.PreAccountCellDataV1FromSlice(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("PreAccountCellDataV1FromSlice err: %s", err.Error())
				}
			default:
				if err := resp.PreAccountCellDataFromSlice(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("PreAccountCellDataFromSlice err: %s", err.Error())
				}
			}

			respMap[resp.Account] = &resp
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist pre account cell")
	}
	return respMap, nil
}

func (p *PreAccountCellDataBuilder) PreAccountCellDataV1FromSlice(bys []byte) error {
	data, err := molecule.PreAccountCellDataV1FromSlice(bys, true)
	if err != nil {
		return fmt.Errorf("PreAccountCellDataV1FromSlice err: %s", err.Error())
	}
	p.preAccountCellDataV1 = data

	p.AccountChars = data.Account()
	p.Account = common.AccountCharsToAccount(data.Account())
	p.RefundLock = data.RefundLock()
	p.OwnerLockArgs = common.Bytes2Hex(data.OwnerLockArgs().RawData())
	p.InviterId = common.Bytes2Hex(data.InviterId().RawData())
	if !data.InviterLock().IsNone() {
		p.InviterLock, _ = data.InviterLock().IntoScript()
	}
	if !data.ChannelLock().IsNone() {
		p.ChannelLock, _ = data.ChannelLock().IntoScript()
	}
	p.Price = data.Price()
	p.Quote = data.Quote()
	p.InvitedDiscount = data.InvitedDiscount()
	p.CreatedAt = data.CreatedAt()

	return nil
}

func (p *PreAccountCellDataBuilder) PreAccountCellDataFromSlice(bys []byte) error {
	data, err := molecule.PreAccountCellDataFromSlice(bys, true)
	if err != nil {
		return fmt.Errorf("PreAccountCellDataV1FromSlice err: %s", err.Error())
	}
	p.preAccountCellData = data

	p.AccountChars = data.Account()
	p.Account = common.AccountCharsToAccount(data.Account())
	p.RefundLock = data.RefundLock()
	p.OwnerLockArgs = common.Bytes2Hex(data.OwnerLockArgs().RawData())
	p.InviterId = common.Bytes2Hex(data.InviterId().RawData())
	if !data.InviterLock().IsNone() {
		p.InviterLock, _ = data.InviterLock().IntoScript()
	}
	if !data.ChannelLock().IsNone() {
		p.ChannelLock, _ = data.ChannelLock().IntoScript()
	}
	p.Price = data.Price()
	p.Quote = data.Quote()
	p.InvitedDiscount = data.InvitedDiscount()
	p.CreatedAt = data.CreatedAt()

	p.InitialRecords = data.InitialRecords()

	return nil
}

type PreAccountCellParam struct {
	OldIndex uint32
	NewIndex uint32
	Status   uint8
	Action   string

	CreatedAt       int64
	InvitedDiscount uint32
	Quote           uint64
	InviterScript   *types.Script
	ChannelScript   *types.Script
	InviterId       []byte
	OwnerLockArgs   []byte
	RefundLock      *types.Script
	Price           molecule.PriceConfig
	AccountChars    molecule.AccountChars

	InitialRecords []Record
}

func (p *PreAccountCellDataBuilder) getOldDataEntityOpt(param *PreAccountCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	switch p.Version {
	case common.GoDataEntityVersion1:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(p.preAccountCellDataV1.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(param.OldIndex)).Build()
	default:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(p.preAccountCellData.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(param.OldIndex)).Build()
	}

	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}
func (p *PreAccountCellDataBuilder) GenWitness(param *PreAccountCellParam) ([]byte, []byte, error) {

	switch param.Action {
	case common.DasActionPreRegister:
		createdAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(param.CreatedAt)).Build()
		invitedDiscount := molecule.GoU32ToMoleculeU32(param.InvitedDiscount)
		quote := molecule.GoU64ToMoleculeU64(param.Quote)
		var iScript, cScript molecule.ScriptOpt
		if param.InviterScript != nil {
			iScript = molecule.NewScriptOptBuilder().Set(molecule.CkbScript2MoleculeScript(param.InviterScript)).Build()
		} else {
			iScript = *molecule.ScriptOptFromSliceUnchecked([]byte{})
		}
		if param.ChannelScript != nil {
			cScript = molecule.NewScriptOptBuilder().Set(molecule.CkbScript2MoleculeScript(param.ChannelScript)).Build()
		} else {
			cScript = *molecule.ScriptOptFromSliceUnchecked([]byte{})
		}
		inviterId := molecule.GoBytes2MoleculeBytes(param.InviterId)
		ownerLockArgs := molecule.GoBytes2MoleculeBytes(param.OwnerLockArgs)
		refundLock := molecule.CkbScript2MoleculeScript(param.RefundLock)

		initialRecords := molecule.RecordsDefault()
		if len(param.InitialRecords) > 0 {
			records := ConvertToCellRecords(param.InitialRecords)
			initialRecords = *records
		}

		preAccountCellData := molecule.NewPreAccountCellDataBuilder().
			Account(param.AccountChars).
			RefundLock(refundLock).
			OwnerLockArgs(ownerLockArgs).
			InviterId(inviterId).
			InviterLock(iScript).
			ChannelLock(cScript).
			Price(param.Price).
			Quote(quote).
			InvitedDiscount(invitedDiscount).
			CreatedAt(createdAt).
			InitialRecords(initialRecords).Build()
		newDataBytes := molecule.GoBytes2MoleculeBytes(preAccountCellData.AsSlice())
		newDataEntity := molecule.NewDataEntityBuilder().Entity(newDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(param.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypePreAccountCell, &tmp)
		return witness, common.Blake2b(preAccountCellData.AsSlice()), nil
	case common.DasActionPropose:
		oldDataEntityOpt := p.getOldDataEntityOpt(param)
		tmp := molecule.NewDataBuilder().Dep(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypePreAccountCell, &tmp)
		return witness, nil, nil
	case common.DasActionConfirmProposal, common.DasActionRefundPreRegister:
		oldDataEntityOpt := p.getOldDataEntityOpt(param)
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypePreAccountCell, &tmp)
		return witness, nil, nil
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", param.Action)
}
