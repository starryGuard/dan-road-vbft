package validation

import (
	"errors"
	"fmt"

	"dan-road-vbft/common"
	"dan-road-vbft/common/constants"
	"dan-road-vbft/common/log"
	"dan-road-vbft/core/ledger"
	"dan-road-vbft/core/payload"
	"dan-road-vbft/core/signature"
	"dan-road-vbft/core/types"
	ontErrors "dan-road-vbft/errors"
)

// VerifyTransaction verifys received single transaction
func VerifyTransaction(tx *types.Transaction) ontErrors.ErrCode {
	if err := checkTransactionSignatures(tx); err != nil {
		log.Info("transaction verify error:", err)
		return ontErrors.ErrVerifySignature
	}

	if err := checkTransactionPayload(tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ontErrors.ErrTransactionPayload
	}

	return ontErrors.ErrNoError
}

func VerifyTransactionWithLedger(tx *types.Transaction, ledger *ledger.Ledger) ontErrors.ErrCode {
	//TODO: replay check
	return ontErrors.ErrNoError
}

func checkTransactionSignatures(tx *types.Transaction) error {
	hash := tx.Hash()

	lensig := len(tx.Sigs)
	if lensig > constants.TX_MAX_SIG_SIZE {
		return fmt.Errorf("transaction signature number %d execced %d", lensig, constants.TX_MAX_SIG_SIZE)
	}

	address := make(map[common.Address]bool, len(tx.Sigs))
	for _, sigdata := range tx.Sigs {
		sig, err := sigdata.GetSig()
		if err != nil {
			return err
		}

		m := int(sig.M)
		kn := len(sig.PubKeys)
		sn := len(sig.SigData)

		if kn > constants.MULTI_SIG_MAX_PUBKEY_SIZE || sn < m || m > kn || m <= 0 {
			return errors.New("wrong tx sig param length")
		}

		if kn == 1 {
			err := signature.Verify(sig.PubKeys[0], hash[:], sig.SigData[0])
			if err != nil {
				return errors.New("signature verification failed")
			}

			address[types.AddressFromPubKey(sig.PubKeys[0])] = true
		} else {
			if err := signature.VerifyMultiSignature(hash[:], sig.PubKeys, m, sig.SigData); err != nil {
				return err
			}

			addr, err := types.AddressFromMultiPubKeys(sig.PubKeys, m)
			if err != nil {
				return err
			}
			address[addr] = true
		}
	}

	// check payer in address
	if address[tx.Payer] == false {
		return errors.New("signature missing for payer: " + tx.Payer.ToBase58())
	}

	addrList := make([]common.Address, 0, len(address))
	for addr := range address {
		addrList = append(addrList, addr)
	}

	tx.SignedAddr = addrList

	return nil
}

func checkTransactionPayload(tx *types.Transaction) error {

	switch pld := tx.Payload.(type) {
	case *payload.DeployCode:
		return nil
	case *payload.InvokeCode:
		return nil
	default:
		return errors.New(fmt.Sprint("[txValidator], unimplemented transaction payload type.", pld))
	}
	return nil
}
