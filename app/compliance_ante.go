package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	compliancemodulekeeper "febelchain/x/compliance/keeper"
)

type ComplianceDecorator struct {
	keeper compliancemodulekeeper.Keeper
}

func NewComplianceDecorator(k compliancemodulekeeper.Keeper) ComplianceDecorator {
	return ComplianceDecorator{keeper: k}
}

func (cd ComplianceDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {

	msgs := tx.GetMsgs()

	for _, msg := range msgs {

		if sendMsg, ok := msg.(*banktypes.MsgSend); ok {

			fromAddr := sendMsg.FromAddress

			whitelisted, err := cd.keeper.IsWhitelisted(ctx, fromAddr)
			if err != nil || !whitelisted {
				return ctx, sdkerrors.ErrUnauthorized.Wrap(
					fmt.Sprintf("address %s is not whitelisted", fromAddr),
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}
