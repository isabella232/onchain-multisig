package vault

import (
	"github.com/bjartek/go-with-the-flow/gwtf"
	util "github.com/flow-hydraulics/onchain-multisig"
	"github.com/onflow/cadence"
)

func AddVaultToAccount(
	g *gwtf.GoWithTheFlow,
	vaultAcct string,
) (events []*gwtf.FormatedEvent, err error) {
	txFilename := "../../../transactions/create_vault.cdc"
	txScript := util.ParseCadenceTemplate(txFilename)

	pk1000 := g.Accounts["w-1000"].PrivateKey.PublicKey().String()
	pk500_1 := g.Accounts["w-500-1"].PrivateKey.PublicKey().String()
	pk500_2 := g.Accounts["w-500-2"].PrivateKey.PublicKey().String()
	pk250_1 := g.Accounts["w-250-1"].PrivateKey.PublicKey().String()
	pk250_2 := g.Accounts["w-250-2"].PrivateKey.PublicKey().String()
	w1000, _ := cadence.NewUFix64("1000.0")
	w500, _ := cadence.NewUFix64("500.0")
	w250, _ := cadence.NewUFix64("250.0")

	multiSigPubKeys := []cadence.Value{
		cadence.String(pk1000[2:]),
		cadence.String(pk500_1[2:]),
		cadence.String(pk500_2[2:]),
		cadence.String(pk250_1[2:]),
		cadence.String(pk250_2[2:]),
	}
	multiSigKeyWeights := []cadence.Value{w1000, w500, w500, w250, w250}

	e, err := g.TransactionFromFile(txFilename, txScript).
		SignProposeAndPayAs(vaultAcct).
		Argument(cadence.NewArray(multiSigPubKeys)).
		Argument(cadence.NewArray(multiSigKeyWeights)).
		Run()
	events = util.ParseTestEvents(e)
	return
}

func AccountSignerTransferTokens(
	g *gwtf.GoWithTheFlow,
	amount string,
	fromAcct string,
	toAcct string,
) (events []*gwtf.FormatedEvent, err error) {
	txFilename := "../../../transactions/account_signer_token_transfer.cdc"
	txScript := util.ParseCadenceTemplate(txFilename)

	e, err := g.TransactionFromFile(txFilename, txScript).
		SignProposeAndPayAs(fromAcct).
		UFix64Argument(amount).
		AccountArgument(toAcct).
		Run()
	events = util.ParseTestEvents(e)
	return
}

func MultiSig_NewPendingTransferPayload(
	g *gwtf.GoWithTheFlow,
	amount string,
	publicKey string,
	signerAcct string,
	vaultAcct string,
) (events []*gwtf.FormatedEvent, err error) {
	txFilename := "../../../transactions/new_pending_transfer.cdc"
	txScript := util.ParseCadenceTemplate(txFilename)

	method := "transfer"

	signable, err := util.GetSignableDataFromScript(g, method, amount)
	if err != nil {
		return
	}

	sig, err := util.SignPayloadOffline(g, signable, signerAcct)
	if err != nil {
		return
	}

	// TODO add to in the signature
	//Argument(cadence.NewArray(sigArray)).
	e, err := g.TransactionFromFile(txFilename, txScript).
		SignProposeAndPayAs(signerAcct).
		StringArgument(publicKey).
		StringArgument(sig).
		AccountArgument(vaultAcct).
		StringArgument(method).
		UFix64Argument(amount).
		Run()
	events = util.ParseTestEvents(e)
	return
}

//func MultiSig_MasterMinterExecuteTx(
//	g *gwtf.GoWithTheFlow,
//	index uint64,
//	ownerAcct string,
//) (events []*gwtf.FormatedEvent, err error) {
//	txFilename := "../../../transactions/owner/multisig/executeTx.cdc"
//	txScript := util.ParseCadenceTemplate(txFilename)
//
//	e, err := g.TransactionFromFile(txFilename, txScript).
//		SignProposeAndPayAs(ownerAcct).
//		AccountArgument("owner").
//		UInt64Argument(index).
//		Run()
//	events = util.ParseTestEvents(e)
//	return
//
//}
