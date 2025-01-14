package common

type ActionDataType = string

const (
	ActionDataTypeActionData         ActionDataType = "0x00000000" // action data
	ActionDataTypeAccountCell        ActionDataType = "0x01000000" // account cell
	ActionDataTypeAccountSaleCell    ActionDataType = "0x02000000" // account sale cell
	ActionDataTypeAccountAuctionCell ActionDataType = "0x03000000" // account auction cell
	ActionDataTypeProposalCell       ActionDataType = "0x04000000" // proposal cell
	ActionDataTypePreAccountCell     ActionDataType = "0x05000000" // pre account cell
	ActionDataTypeIncomeCell         ActionDataType = "0x06000000" // income cell
	ActionDataTypeOfferCell          ActionDataType = "0x07000000" // offer cell
	ActionDataTypeSubAccount         ActionDataType = "0x08000000" // sub account
	ActionDataTypeSubAccountMintSign ActionDataType = "0x09000000"
)

const (
	WitnessDas                  = "das"
	WitnessDasCharLen           = 3
	WitnessDasTableTypeEndIndex = 7
)

type DataType = int

const (
	DataTypeNew          DataType = 0
	DataTypeOld          DataType = 1
	DataTypeDep          DataType = 2
	GoDataEntityVersion1 uint32   = 1
	GoDataEntityVersion2 uint32   = 2
	GoDataEntityVersion3 uint32   = 3
)

type EditKey = string

const (
	EditKeyOwner     EditKey = "owner"
	EditKeyManager   EditKey = "manager"
	EditKeyRecords   EditKey = "records"
	EditKeyExpiredAt EditKey = "expired_at"
)

type SubAction = string

const (
	SubActionCreate SubAction = "create"
	SubActionEdit   SubAction = "edit"
	SubActionRenew  SubAction = "renew"
)
