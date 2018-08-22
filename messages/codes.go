package messages

// constant Request Codes
const (
	C_JOIN_REQUEST     = 1
	C_LTPK_REQUEST     = 2
	C_KEY_EXCHANGE     = 3
	C_EXP_DC_VECTOR    = 4
	C_SIMPLE_DC_VECTOR = 5
	C_TX_CONFIRMATION  = 6
	C_KESK_RESPONSE    = 7
)

// constant Response Codes
const (
	S_JOIN_RESPONSE    = 101
	S_START_DICEMIX    = 102
	S_KEY_EXCHANGE     = 103
	S_EXP_DC_VECTOR    = 104
	S_SIMPLE_DC_VECTOR = 105
	S_TX_SUCCESSFUL    = 106
	S_KESK_REQUEST     = 107
)
