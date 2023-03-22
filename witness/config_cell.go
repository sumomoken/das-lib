package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

type ConfigCellDataBuilder struct {
	ConfigCellAccount                *molecule.ConfigCellAccount
	ConfigCellPrice                  *molecule.ConfigCellPrice
	PriceConfigMap                   map[uint8]*molecule.PriceConfig
	PriceMaxLength                   uint8
	ConfigCellSecondaryMarket        *molecule.ConfigCellSecondaryMarket
	ConfigCellIncome                 *molecule.ConfigCellIncome
	ConfigCellProfitRate             *molecule.ConfigCellProfitRate
	ConfigCellMain                   *molecule.ConfigCellMain
	ConfigCellReverseResolution      *molecule.ConfigCellReverseResolution
	ConfigCellProposal               *molecule.ConfigCellProposal
	ConfigCellApply                  *molecule.ConfigCellApply
	ConfigCellRelease                *molecule.ConfigCellRelease
	ConfigCellSubAccount             *molecule.ConfigCellSubAccount
	ConfigCellSystemStatus           *molecule.ConfigCellSystemStatus
	ConfigCellRecordKeys             []string
	ConfigCellEmojis                 []string
	ConfigCellCharSetDigit           []string
	ConfigCellCharSetEn              []string
	ConfigCellCharSetHanS            []string
	ConfigCellCharSetHanT            []string
	ConfigCellCharSetJa              []string
	ConfigCellCharSetKo              []string
	ConfigCellCharSetRu              []string
	ConfigCellCharSetTr              []string
	ConfigCellCharSetTh              []string
	ConfigCellCharSetVi              []string
	ConfigCellUnavailableAccountMap  map[string]struct{}
	ConfigCellPreservedAccountMap    map[string]struct{}
	ConfigCellSubAccountWhiteListMap map[string]struct{}
}

func ConfigCellDataBuilderRefByTypeArgs(builder *ConfigCellDataBuilder, tx *types.Transaction, configCellTypeArgs common.ConfigCellTypeArgs) error {
	var configCellDataBys []byte
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		if actionDataType == configCellTypeArgs {
			configCellDataBys = dataBys
			return false, nil
		}
		return true, nil
	})
	fmt.Println(configCellDataBys)
	if err != nil {
		return fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}

	switch configCellTypeArgs {
	case common.ConfigCellTypeArgsAccount:
		ConfigCellAccount, err := molecule.ConfigCellAccountFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellAccountFromSlice err: %s", err.Error())
		}
		builder.ConfigCellAccount = ConfigCellAccount
	case common.ConfigCellTypeArgsPrice:
		ConfigCellPrice, err := molecule.ConfigCellPriceFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellPriceFromSlice err: %s", err.Error())
		}
		builder.ConfigCellPrice = ConfigCellPrice
		builder.PriceConfigMap = make(map[uint8]*molecule.PriceConfig)
		prices := builder.ConfigCellPrice.Prices()
		for i, count := uint(0), prices.Len(); i < count; i++ {
			price, err := molecule.PriceConfigFromSlice(prices.Get(i).AsSlice(), true)
			if err != nil {
				return fmt.Errorf("PriceConfigFromSlice err: %s", err.Error())
			}
			length, err := molecule.Bytes2GoU8(price.Length().RawData())
			if err != nil {
				return fmt.Errorf("price.Length() err: %s", err.Error())
			}
			if builder.PriceMaxLength < length {
				builder.PriceMaxLength = length
			}
			builder.PriceConfigMap[length] = price
		}
	case common.ConfigCellTypeArgsApply:
		ConfigCellApply, err := molecule.ConfigCellApplyFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
		}
		builder.ConfigCellApply = ConfigCellApply
	case common.ConfigCellTypeArgsRelease:
		ConfigCellRelease, err := molecule.ConfigCellReleaseFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
		}
		builder.ConfigCellRelease = ConfigCellRelease
	case common.ConfigCellTypeArgsSecondaryMarket:
		ConfigCellSecondaryMarket, err := molecule.ConfigCellSecondaryMarketFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellSecondaryMarketFromSlice err: %s", err.Error())
		}
		builder.ConfigCellSecondaryMarket = ConfigCellSecondaryMarket
	case common.ConfigCellTypeArgsIncome:
		ConfigCellIncome, err := molecule.ConfigCellIncomeFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellIncomeFromSlice err: %s", err.Error())
		}
		builder.ConfigCellIncome = ConfigCellIncome
	case common.ConfigCellTypeArgsProfitRate:
		ConfigCellProfitRate, err := molecule.ConfigCellProfitRateFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
		}
		builder.ConfigCellProfitRate = ConfigCellProfitRate
	case common.ConfigCellTypeArgsMain:
		ConfigCellMain, err := molecule.ConfigCellMainFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellMainFromSlice err: %s", err.Error())
		}
		builder.ConfigCellMain = ConfigCellMain
	case common.ConfigCellTypeArgsReverseRecord:
		ConfigCellReverseResolution, err := molecule.ConfigCellReverseResolutionFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellReverseResolutionFromSlice err: %s", err.Error())
		}
		builder.ConfigCellReverseResolution = ConfigCellReverseResolution
	case common.ConfigCellTypeArgsSubAccount:
		ConfigCellSubAccount, err := molecule.ConfigCellSubAccountFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellSubAccountFromSlice err: %s", err.Error())
		}
		builder.ConfigCellSubAccount = ConfigCellSubAccount
	case common.ConfigCellTypeArgsProposal:
		ConfigCellProposal, err := molecule.ConfigCellProposalFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellProposalFromSlice err: %s", err.Error())
		}
		builder.ConfigCellProposal = ConfigCellProposal
	case common.ConfigCellTypeArgsRecordNamespace:
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("key name space len err: %s", err.Error())
		}
		builder.ConfigCellRecordKeys = strings.Split(string(configCellDataBys[4:dataLength]), string([]byte{0x00}))

	case common.ConfigCellTypeArgsCharSetEmoji:
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("char set emoji err: %s", err.Error())
		}
		builder.ConfigCellEmojis = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
	case common.ConfigCellTypeArgsCharSetDigit:
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("char set digit err: %s", err.Error())
		}
		builder.ConfigCellCharSetDigit = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
	case common.ConfigCellTypeArgsCharSetEn:
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("char set en err: %s", err.Error())
		}
		builder.ConfigCellCharSetEn = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
	case common.ConfigCellTypeArgsCharSetHanS:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set hans err: %s", err.Error())
			}
			builder.ConfigCellCharSetHanS = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetHanT:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set hant err: %s", err.Error())
			}
			builder.ConfigCellCharSetHanT = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetJa:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set jp err: %s", err.Error())
			}
			builder.ConfigCellCharSetJa = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetKo:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set kr err: %s", err.Error())
			}
			builder.ConfigCellCharSetKo = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetRu:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set ru err: %s", err.Error())
			}
			builder.ConfigCellCharSetRu = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetTr:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set tr err: %s", err.Error())
			}
			builder.ConfigCellCharSetTr = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetTh:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set th err: %s", err.Error())
			}
			builder.ConfigCellCharSetTh = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsCharSetVi:
		if len(configCellDataBys) != 0 {
			dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
			if err != nil {
				return fmt.Errorf("char set vn err: %s", err.Error())
			}
			builder.ConfigCellCharSetVi = strings.Split(string(configCellDataBys[5:dataLength]), string([]byte{0x00}))
		}
	case common.ConfigCellTypeArgsUnavailable:
		if builder.ConfigCellUnavailableAccountMap == nil {
			builder.ConfigCellUnavailableAccountMap = make(map[string]struct{})
		}
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("unavailable account err: %s", err.Error())
		}
		for i := 20; i <= len(configCellDataBys[4:dataLength]); i += 20 {
			tmp := common.Bytes2Hex(configCellDataBys[4:dataLength][i-20 : i])
			if _, ok := builder.ConfigCellUnavailableAccountMap[tmp]; ok {
				fmt.Println(tmp, "ok")
			}
			builder.ConfigCellUnavailableAccountMap[tmp] = struct{}{}
		}
	case common.ConfigCellTypeArgsSubAccountWhiteList:
		if builder.ConfigCellSubAccountWhiteListMap == nil {
			builder.ConfigCellSubAccountWhiteListMap = make(map[string]struct{})
		}
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("preserved account err: %s", err.Error())
		}
		for i := 20; i <= len(configCellDataBys[4:dataLength]); i += 20 {
			tmp := common.Bytes2Hex(configCellDataBys[4:dataLength][i-20 : i])
			builder.ConfigCellSubAccountWhiteListMap[tmp] = struct{}{}
		}
	case common.ConfigCellTypeArgsSystemStatus:
		configCellSystemStatus, err := molecule.ConfigCellSystemStatusFromSlice(configCellDataBys, true)
		if err != nil {
			return fmt.Errorf("ConfigCellSystemStatusFromSlice err: %s", err.Error())
		}
		builder.ConfigCellSystemStatus = configCellSystemStatus
	case common.ConfigCellTypeArgsPreservedAccount00,
		common.ConfigCellTypeArgsPreservedAccount01,
		common.ConfigCellTypeArgsPreservedAccount02,
		common.ConfigCellTypeArgsPreservedAccount03,
		common.ConfigCellTypeArgsPreservedAccount04,
		common.ConfigCellTypeArgsPreservedAccount05,
		common.ConfigCellTypeArgsPreservedAccount06,
		common.ConfigCellTypeArgsPreservedAccount07,
		common.ConfigCellTypeArgsPreservedAccount08,
		common.ConfigCellTypeArgsPreservedAccount09,
		common.ConfigCellTypeArgsPreservedAccount10,
		common.ConfigCellTypeArgsPreservedAccount11,
		common.ConfigCellTypeArgsPreservedAccount12,
		common.ConfigCellTypeArgsPreservedAccount13,
		common.ConfigCellTypeArgsPreservedAccount14,
		common.ConfigCellTypeArgsPreservedAccount15,
		common.ConfigCellTypeArgsPreservedAccount16,
		common.ConfigCellTypeArgsPreservedAccount17,
		common.ConfigCellTypeArgsPreservedAccount18,
		common.ConfigCellTypeArgsPreservedAccount19:
		if builder.ConfigCellPreservedAccountMap == nil {
			builder.ConfigCellPreservedAccountMap = make(map[string]struct{})
		}
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("preserved account err: %s", err.Error())
		}
		for i := 20; i <= len(configCellDataBys[4:dataLength]); i += 20 {
			tmp := common.Bytes2Hex(configCellDataBys[4:dataLength][i-20 : i])
			builder.ConfigCellPreservedAccountMap[tmp] = struct{}{}
		}
	}
	return nil
}

func ConfigCellDataBuilderByTypeArgs(tx *types.Transaction, configCellTypeArgs common.ConfigCellTypeArgs) (*ConfigCellDataBuilder, error) {
	var resp ConfigCellDataBuilder

	err := ConfigCellDataBuilderRefByTypeArgs(&resp, tx, configCellTypeArgs)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderRefByTypeArgs err: %s", err.Error())
	}

	return &resp, nil
}

func (c *ConfigCellDataBuilder) PriceInvitedDiscount() (uint32, error) {
	if c.ConfigCellPrice != nil {
		return molecule.Bytes2GoU32(c.ConfigCellPrice.Discount().InvitedDiscount().RawData())
	}
	return 0, fmt.Errorf("ConfigCellPrice is nil")
}

func (c *ConfigCellDataBuilder) RecordBasicCapacity() (uint64, error) {
	if c.ConfigCellReverseResolution != nil {
		return molecule.Bytes2GoU64(c.ConfigCellReverseResolution.RecordBasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellReverseResolution is nil")
}

func (c *ConfigCellDataBuilder) RecordPreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellReverseResolution != nil {
		return molecule.Bytes2GoU64(c.ConfigCellReverseResolution.RecordPreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellReverseResolution is nil")
}

func (c *ConfigCellDataBuilder) RecordCommonFee() (uint64, error) {
	if c.ConfigCellReverseResolution != nil {
		return molecule.Bytes2GoU64(c.ConfigCellReverseResolution.CommonFee().RawData())
	}
	return 0, fmt.Errorf("ConfigCellReverseResolution is nil")
}

func (c *ConfigCellDataBuilder) BasicCapacity() (uint64, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellAccount.BasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) BasicCapacityFromOwnerDasAlgorithmId(args string) (uint64, error) {
	if args == "" {
		return 0, fmt.Errorf("args is nil")
	}
	argsByte := common.Hex2Bytes(args)
	algorithmId := common.DasAlgorithmId(argsByte[0])
	switch algorithmId {
	case common.DasAlgorithmIdEd25519:
		return 230 * common.OneCkb, nil
	default:
		if c.ConfigCellAccount != nil {
			return molecule.Bytes2GoU64(c.ConfigCellAccount.BasicCapacity().RawData())
		}
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) TransferAccountThrottle() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.TransferAccountThrottle().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) EditRecordsThrottle() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.EditRecordsThrottle().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) RecordMinTtl() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.RecordMinTtl().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) ExpirationGracePeriod() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.ExpirationGracePeriod().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) EditManagerThrottle() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.EditManagerThrottle().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) MaxLength() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.MaxLength().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) AccountCommonFee() (uint64, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellAccount.CommonFee().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) PreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellAccount.PreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) AccountPrice(length uint8) (uint64, uint64, error) {
	if length > 5 {
		length = 5
	}
	if c.PriceConfigMap != nil {
		price, ok := c.PriceConfigMap[length]
		if ok {
			newPrice, err := molecule.Bytes2GoU64(price.New().RawData())
			if err != nil {
				return 0, 0, fmt.Errorf("price.New() err: %s", err.Error())
			}
			renewPrice, err := molecule.Bytes2GoU64(price.Renew().RawData())
			if err != nil {
				return 0, 0, fmt.Errorf("price.Renew() err: %s", err.Error())
			}
			return newPrice, renewPrice, nil
		}
	}
	return 0, 0, fmt.Errorf("not exist price of length[%d]", length)
}

func (c *ConfigCellDataBuilder) PriceConfig(length uint8) *molecule.PriceConfig {
	if length > c.PriceMaxLength {
		length = c.PriceMaxLength
	}
	if c.PriceConfigMap != nil {
		if price, ok := c.PriceConfigMap[length]; ok {
			return price
		}
	}
	return nil
}

func (c *ConfigCellDataBuilder) SaleCellBasicCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.SaleCellBasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) SaleMinPrice() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.SaleMinPrice().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) SaleCellPreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.SaleCellPreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) CommonFee() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.CommonFee().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferCellBasicCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.OfferCellBasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferMinPrice() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.OfferMinPrice().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferCellPreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.OfferCellPreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferMessageBytesLimit() (uint32, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU32(c.ConfigCellSecondaryMarket.OfferMessageBytesLimit().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) IncomeBasicCapacity() (uint64, error) {
	if c.ConfigCellIncome != nil {
		return molecule.Bytes2GoU64(c.ConfigCellIncome.BasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellIncome is nil")
}

func (c *ConfigCellDataBuilder) IncomeMinTransferCapacity() (uint64, error) {
	if c.ConfigCellIncome != nil {
		return molecule.Bytes2GoU64(c.ConfigCellIncome.MinTransferCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellIncome is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateChannel() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.Channel().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateProposalCreate() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.ProposalCreate().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateProposalConfirm() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.ProposalConfirm().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateIncomeConsolidate() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.IncomeConsolidate().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateSaleBuyerInviter() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.SaleBuyerInviter().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateSaleBuyerChannel() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.SaleBuyerChannel().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateSaleDas() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.SaleDas().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateInviter() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.Inviter().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) Status() (uint8, error) {
	if c.ConfigCellMain != nil {
		return molecule.Bytes2GoU8(c.ConfigCellMain.Status().RawData())
	}
	return 0, fmt.Errorf("ConfigCellMain is nil")
}

func (c *ConfigCellDataBuilder) LuckyNumber() (uint32, error) {
	if c.ConfigCellRelease != nil {
		return molecule.Bytes2GoU32(c.ConfigCellRelease.LuckyNumber().RawData())
	}
	return 0, fmt.Errorf("ConfigCellRelease is nil")
}

func (c *ConfigCellDataBuilder) NewSubAccountPrice() (uint64, error) {
	if c.ConfigCellSubAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSubAccount.NewSubAccountPrice().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSubAccount is nil")
}

func (c *ConfigCellDataBuilder) RenewSubAccountPrice() (uint64, error) {
	if c.ConfigCellSubAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSubAccount.RenewSubAccountPrice().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSubAccount is nil")
}

func (c *ConfigCellDataBuilder) GetContractStatus(contractName common.DasContractName) (res common.ContractStatus, err error) {
	if c.ConfigCellSystemStatus == nil {
		err = fmt.Errorf("ConfigCellSystemStatus is nil")
		return
	}
	switch contractName {
	case common.DasContractNameApplyRegisterCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.ApplyRegisterCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.ApplyRegisterCellType().Version().RawData())
	case common.DasContractNamePreAccountCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.PreAccountCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.PreAccountCellType().Version().RawData())
	case common.DasContractNameProposalCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.ProposalCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.ProposalCellType().Version().RawData())
	case common.DasContractNameConfigCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.ConfigCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.ConfigCellType().Version().RawData())
	case common.DasContractNameAccountCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.AccountCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.AccountCellType().Version().RawData())
	case common.DasContractNameAccountSaleCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.AccountSaleCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.AccountSaleCellType().Version().RawData())
	case common.DASContractNameSubAccountCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.SubAccountCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.SubAccountCellType().Version().RawData())
	case common.DASContractNameOfferCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.OfferCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.OfferCellType().Version().RawData())
	case common.DasContractNameBalanceCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.BalanceCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.BalanceCellType().Version().RawData())
	case common.DasContractNameIncomeCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.IncomeCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.IncomeCellType().Version().RawData())
	case common.DasContractNameReverseRecordCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.ReverseRecordCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.ReverseRecordCellType().Version().RawData())
	case common.DASContractNameEip712LibCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.Eip712Lib().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.Eip712Lib().Version().RawData())
	case common.DasContractNameReverseRecordRootCellType:
		res.Status, _ = molecule.Bytes2GoU8(c.ConfigCellSystemStatus.ReverseRecordRootCellType().Status().AsSlice())
		res.Version = string(c.ConfigCellSystemStatus.ReverseRecordRootCellType().Version().RawData())
	default:
		err = fmt.Errorf("unknow contract-name[%s]", contractName)
		return
	}
	return
}
