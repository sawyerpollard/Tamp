package encoders

type Code string
type Source string

// Encoder defines simple encoding methods.
type Encoder interface {
	Encode(source Source) Code
	Decode(code Code) Source
}
