package witness

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

var (
	ErrNotExistNewIncomeCell = errors.New("not exist new income cell")
)

const (
	IncomeCellCurrentVersion = common.GoDataEntityVersion1
)

type IncomeCellDataBuilder struct {
	Index          uint32
	Version        uint32
	IncomeCellData *molecule.IncomeCellData
	DataEntityOpt  *molecule.DataEntityOpt
}

type IncomeCellRecord struct {
	BelongTo *molecule.Script
	Capacity uint64
}

func IncomeCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*IncomeCellDataBuilder, error) {
	respList, err := IncomeCellDataBuilderListFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	return respList[0], nil
}

func IncomeCellDataBuilderListFromTx(tx *types.Transaction, dataType common.DataType) ([]*IncomeCellDataBuilder, error) {
	var respList = make([]*IncomeCellDataBuilder, 0)
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeIncomeCell:
			var resp IncomeCellDataBuilder
			dataEntityOpt, dataEntity, err := getDataEntityOpt(dataBys, dataType)
			if err != nil {
				if err == ErrDataEntityOptIsNil {
					//log.Warn("getDataEntityOpt err:", err.Error(), tx.Hash)
					return true, nil
				}
				return false, fmt.Errorf("getDataEntityOpt err: %s", err.Error())
			}
			resp.DataEntityOpt = dataEntityOpt
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

			incomeCellData, err := molecule.IncomeCellDataFromSlice(dataEntity.Entity().RawData(), true)
			if err != nil {
				return false, fmt.Errorf("IncomeCellDataFromSlice err: %s", err.Error())
			}
			resp.IncomeCellData = incomeCellData
			respList = append(respList, &resp)
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respList) == 0 {
		return nil, ErrNotExistNewIncomeCell
	}
	return respList, nil
}

func (i *IncomeCellDataBuilder) Creator() *molecule.Script {
	if i.IncomeCellData != nil {
		return i.IncomeCellData.Creator()
	}
	return nil
}

func (i *IncomeCellDataBuilder) Records() []IncomeCellRecord {
	var list []IncomeCellRecord
	if i.IncomeCellData != nil {
		for index, lenRecords := uint(0), i.IncomeCellData.Records().Len(); index < lenRecords; index++ {
			record := i.IncomeCellData.Records().Get(index)
			capacity, _ := molecule.Bytes2GoU64(record.Capacity().RawData())
			list = append(list, IncomeCellRecord{
				BelongTo: record.BelongTo(),
				Capacity: capacity,
			})
		}
	}
	return list
}

type IncomeCellParam struct {
	OldRecordsDataList []*molecule.IncomeCellData
	Creators           []*molecule.Script
	NewRecords         []*molecule.IncomeRecords
	Action             string
}

func GenBatchIncomeWitnessData(p *IncomeCellParam) ([][]byte, [][]byte, error) {
	if p == nil {
		return nil, nil, fmt.Errorf("param shouldn't be nil")
	}
	if len(p.Creators) != len(p.NewRecords) {
		return nil, nil, fmt.Errorf("param invalid, length error")
	}
	var incomeCellWitness [][]byte
	var newWitHash [][]byte
	switch p.Action {
	case common.DasActionConsolidateIncome:
		for index, oldData := range p.OldRecordsDataList {
			oldCellDataBytes := molecule.GoBytes2MoleculeBytes(oldData.AsSlice())
			oldDataEntity := molecule.NewDataEntityBuilder().Entity(oldCellDataBytes).
				Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(uint32(index))).Build()
			oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
			tmp := molecule.NewDataBuilder().Old(oldDataEntityOpt).Build()

			incomeCellWitness = append(incomeCellWitness, GenDasDataWitness(common.ActionDataTypeIncomeCell, &tmp))
		}
		for index, creator := range p.Creators {
			newIncomeWitData := molecule.NewIncomeCellDataBuilder().
				Creator(*creator).Records(*p.NewRecords[index]).Build()
			newIncomeWitDataBytes := molecule.GoBytes2MoleculeBytes(newIncomeWitData.AsSlice())
			newDataEntity := molecule.NewDataEntityBuilder().Entity(newIncomeWitDataBytes).
				Version(DataEntityVersion1).Index(molecule.GoU32ToMoleculeU32(uint32(index))).Build()
			newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
			newIncomeWit := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
			incomeCellWitness = append(incomeCellWitness, GenDasDataWitness(common.ActionDataTypeIncomeCell, &newIncomeWit))
			newWitHash = append(newWitHash, common.Blake2b(newIncomeWitData.AsSlice()))
		}
	}
	return incomeCellWitness, newWitHash, nil
}

type NewIncomeCellParam struct {
	Creator     *molecule.Script
	BelongTos   []*types.Script
	Capacities  []uint64
	OutputIndex uint32
}

// CreateIncomeCellWitness the first element in belongTos and capacities must be the one who create the income cell
func CreateIncomeCellWitness(p *NewIncomeCellParam) ([]byte, []byte, error) {
	if len(p.BelongTos) != len(p.Capacities) {
		return nil, nil, fmt.Errorf("param invalid")
	}
	distinctBelongTos := []*types.Script{p.BelongTos[0]}
	distinctCapacities := []uint64{p.Capacities[0]}
	var distinctIncomeBelongTos []*types.Script
	var distinctIncomeCapacities []uint64
	for i, v0 := range p.BelongTos[1:] {
		found := false
		for j, v1 := range distinctIncomeBelongTos {
			h1, _ := v1.Hash()
			h2, _ := v0.Hash()
			if h1 == h2 {
				distinctIncomeCapacities[j] = distinctIncomeCapacities[j] + p.Capacities[i+1]
				found = true
			}
		}
		if !found {
			distinctIncomeBelongTos = append(distinctIncomeBelongTos, v0)
			distinctIncomeCapacities = append(distinctIncomeCapacities, p.Capacities[i+1])
		}
	}
	distinctBelongTos = append(distinctBelongTos, distinctIncomeBelongTos...)
	distinctCapacities = append(distinctCapacities, distinctIncomeCapacities...)

	incomeRecords := molecule.NewIncomeRecordsBuilder()
	for index, v := range distinctBelongTos {
		if distinctCapacities[index] == 0 {
			continue
		}
		incomeRecord := molecule.NewIncomeRecordBuilder().
			Capacity(molecule.GoU64ToMoleculeU64(distinctCapacities[index])).
			BelongTo(molecule.CkbScript2MoleculeScript(v)).
			Build()
		incomeRecords.Push(incomeRecord)
	}
	incomeCellData := molecule.NewIncomeCellDataBuilder().
		Creator(*p.Creator).
		Records(incomeRecords.Build()).
		Build()
	version := molecule.GoU32ToMoleculeU32(1)
	newBytes := molecule.GoBytes2MoleculeBytes(incomeCellData.AsSlice())
	newEntity := molecule.NewDataEntityBuilder().Entity(newBytes).Version(version).
		Index(molecule.GoU32ToMoleculeU32(p.OutputIndex)).Build()
	newOpt := molecule.NewDataEntityOptBuilder().Set(newEntity).Build()
	witnessData := molecule.NewDataBuilder().New(newOpt).Build()

	witness := GenDasDataWitness(common.ActionDataTypeIncomeCell, &witnessData)
	return witness, common.Blake2b(incomeCellData.AsSlice()), nil
}

func GenIncomeCellWitness(p *NewIncomeCellParam) ([]byte, []byte, error) {
	if len(p.BelongTos) != len(p.Capacities) {
		return nil, nil, fmt.Errorf("param invalid")
	}
	var distinctBelongTos []*types.Script
	var distinctCapacities []uint64
	var distinctIncomeBelongTos []*types.Script
	var distinctIncomeCapacities []uint64
	for i, v0 := range p.BelongTos {
		found := false
		for j, v1 := range distinctIncomeBelongTos {
			h1, _ := v1.Hash()
			h2, _ := v0.Hash()
			if h1 == h2 {
				distinctIncomeCapacities[j] = distinctIncomeCapacities[j] + p.Capacities[i]
				found = true
			}
		}
		if !found {
			distinctIncomeBelongTos = append(distinctIncomeBelongTos, v0)
			distinctIncomeCapacities = append(distinctIncomeCapacities, p.Capacities[i])
		}
	}
	distinctBelongTos = append(distinctBelongTos, distinctIncomeBelongTos...)
	distinctCapacities = append(distinctCapacities, distinctIncomeCapacities...)

	incomeRecords := molecule.NewIncomeRecordsBuilder()
	for index, v := range distinctBelongTos {
		if distinctCapacities[index] == 0 {
			continue
		}
		incomeRecord := molecule.NewIncomeRecordBuilder().
			Capacity(molecule.GoU64ToMoleculeU64(distinctCapacities[index])).
			BelongTo(molecule.CkbScript2MoleculeScript(v)).
			Build()
		incomeRecords.Push(incomeRecord)
	}
	incomeCellData := molecule.NewIncomeCellDataBuilder().
		Creator(*p.Creator).
		Records(incomeRecords.Build()).
		Build()
	version := molecule.GoU32ToMoleculeU32(1)
	newBytes := molecule.GoBytes2MoleculeBytes(incomeCellData.AsSlice())
	newEntity := molecule.NewDataEntityBuilder().Entity(newBytes).Version(version).
		Index(molecule.GoU32ToMoleculeU32(p.OutputIndex)).Build()
	newOpt := molecule.NewDataEntityOptBuilder().Set(newEntity).Build()
	witnessData := molecule.NewDataBuilder().New(newOpt).Build()

	witness := GenDasDataWitness(common.ActionDataTypeIncomeCell, &witnessData)
	return witness, common.Blake2b(incomeCellData.AsSlice()), nil
}

type ParamNewIncomeCellWitness struct {
	Creator      *types.Script
	CapacityList []uint64
	BelongTo     []*types.Script
	OldIndex     uint32
	NewIndex     uint32
}

func (i *IncomeCellDataBuilder) getOldDataEntityOpt(p *ParamNewIncomeCellWitness) *molecule.DataEntityOpt {
	oldDataEntityOpt := molecule.DataEntityOptDefault()
	if i.IncomeCellData != nil {
		if i.Version == 0 {
			i.Version = IncomeCellCurrentVersion
		}
		version := molecule.GoU32ToMoleculeU32(i.Version)
		index := molecule.GoU32ToMoleculeU32(p.OldIndex)

		oldIncomeCellDataBytes := molecule.GoBytes2MoleculeBytes(i.IncomeCellData.AsSlice())
		oldDataEntity := molecule.NewDataEntityBuilder().Entity(oldIncomeCellDataBytes).
			Version(version).Index(index).Build()

		oldDataEntityOpt = molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	}
	return &oldDataEntityOpt
}

func (i *IncomeCellDataBuilder) getNewDataEntityOpt(p *ParamNewIncomeCellWitness) (*molecule.DataEntityOpt, *molecule.IncomeCellData) {
	var incomeCellData molecule.IncomeCellData

	if i.IncomeCellData != nil {
		incomeCellDataBuilder := i.IncomeCellData.AsBuilder()
		incomeRecordsBuilder := i.IncomeCellData.Records().AsBuilder()
		for index, v := range p.BelongTo {
			incomeRecord := molecule.NewIncomeRecordBuilder().
				Capacity(molecule.GoU64ToMoleculeU64(p.CapacityList[index])).
				BelongTo(molecule.CkbScript2MoleculeScript(v)).
				Build()
			incomeRecordsBuilder.Push(incomeRecord)
		}
		incomeCellData = incomeCellDataBuilder.Records(incomeRecordsBuilder.Build()).Build()
	} else {
		incomeRecords := molecule.NewIncomeRecordsBuilder()
		for index, v := range p.BelongTo {
			incomeRecord := molecule.NewIncomeRecordBuilder().
				Capacity(molecule.GoU64ToMoleculeU64(p.CapacityList[index])).
				BelongTo(molecule.CkbScript2MoleculeScript(v)).
				Build()
			incomeRecords.Push(incomeRecord)
		}
		creator := molecule.ScriptDefault()
		if p.Creator != nil {
			creator = molecule.CkbScript2MoleculeScript(p.Creator)
		}
		incomeCellData = molecule.NewIncomeCellDataBuilder().
			Creator(creator).
			Records(incomeRecords.Build()).
			Build()
	}

	version := molecule.GoU32ToMoleculeU32(IncomeCellCurrentVersion)
	newBytes := molecule.GoBytes2MoleculeBytes(incomeCellData.AsSlice())
	newEntity := molecule.NewDataEntityBuilder().Entity(newBytes).Version(version).
		Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()

	newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newEntity).Build()

	return &newDataEntityOpt, &incomeCellData
}

func (i *IncomeCellDataBuilder) NewIncomeCellWitness(p *ParamNewIncomeCellWitness) (*IncomeCellDataBuilder, []byte, []byte, error) {
	if p == nil || len(p.CapacityList) == 0 || len(p.CapacityList) != len(p.BelongTo) {
		return nil, nil, nil, fmt.Errorf("param invaild")
	}

	oldDataEntityOpt := i.getOldDataEntityOpt(p)
	newDataEntityOpt, incomeCellData := i.getNewDataEntityOpt(p)

	witnessData := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(*newDataEntityOpt).Build()

	witness := GenDasDataWitness(common.ActionDataTypeIncomeCell, &witnessData)
	return &IncomeCellDataBuilder{
		Index:          p.NewIndex,
		Version:        IncomeCellCurrentVersion,
		IncomeCellData: incomeCellData,
		DataEntityOpt:  newDataEntityOpt,
	}, witness, common.Blake2b(incomeCellData.AsSlice()), nil
}
