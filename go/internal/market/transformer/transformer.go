package transformer

import (
	"github.com/HTM1000/homebroker/go/internal/market/dto"
	"github.com/HTM1000/homebroker/go/internal/market/entity"
)

func TransformInput(input dto.TradeInputDTO) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPoition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPoition)
	}

	return order
}

func TransformOutput(order *entity.Order) *dto.OrderOutputDTO {
	outputDTO := &dto.OrderOutputDTO{
		OrderID:    order.ID,
		InvestorID: order.Investor.ID,
		AssetID:    order.Asset.ID,
		OrderType:  order.OrderType,
		Status:     order.Status,
		Partial:    float64(order.PendingShares),
		Shares:     order.Shares,
	}

	var transactionsOutput []*dto.TransactionOutputDTO

	for _, t := range order.Transactions {
		transactionOutput := &dto.TransactionOutputDTO{
			TransactionID: t.ID,
			BuyerID:       t.BuyingOrder.ID,
			SellerID:      t.SellingOrder.ID,
			AssetID:       t.SellingOrder.Asset.ID,
			Price:         t.Price,
			Shares:        t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}
		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	outputDTO.TransactionOutputDTO = transactionsOutput

	return outputDTO
}
