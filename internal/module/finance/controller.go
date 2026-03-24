package finance

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	GetPGConfig(c *fiber.Ctx) error
	UpsertPGConfig(c *fiber.Ctx) error

	ListCampaigns(c *fiber.Ctx) error
	ListPublicCampaigns(c *fiber.Ctx) error
	CreateCampaign(c *fiber.Ctx) error
	GetCampaign(c *fiber.Ctx) error
	GetPublicCampaignBySlug(c *fiber.Ctx) error
	UpdateCampaign(c *fiber.Ctx) error

	ListTransactions(c *fiber.Ctx) error
	ListPublicDonors(c *fiber.Ctx) error
	// CreateDonation(c *fiber.Ctx) error // Nanti di Tahap 4
}
