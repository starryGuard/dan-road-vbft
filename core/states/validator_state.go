package states

import (
	"io"

	"dan-road-vbft/common/serialization"
	"dan-road-vbft/crypto/keypair"
	"dan-road-vbft/errors"
)

type ValidatorState struct {
	StateBase
	PublicKey keypair.PublicKey
}

func (this *ValidatorState) Serialize(w io.Writer) error {
	this.StateBase.Serialize(w)
	buf := keypair.SerializePublicKey(this.PublicKey)
	if err := serialization.WriteVarBytes(w, buf); err != nil {
		return err
	}
	return nil
}

func (this *ValidatorState) Deserialize(r io.Reader) error {
	err := this.StateBase.Deserialize(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ValidatorState], StateBase Deserialize failed.")
	}
	buf, err := serialization.ReadVarBytes(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	pk, err := keypair.DeserializePublicKey(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	this.PublicKey = pk
	return nil
}
