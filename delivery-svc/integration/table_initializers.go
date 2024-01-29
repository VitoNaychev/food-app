package integration

import (
	"testing"

	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
)

func initLocationsTable(t testing.TB, locationStore *models.PgLocationStore) {
	location := testdata.VolenLocation

	testutil.AssertNoErr(t, locationStore.CreateLocation(&location))
}

func initCouriersTable(t testing.TB, courierStore *models.PgCourierStore) {
	courier := testdata.VolenCourier

	testutil.AssertNoErr(t, courierStore.CreateCourier(&courier))
}

func initAddressesTable(t testing.TB, addressStore *models.PgAddressStore) {
	pickupAddress := testdata.VolenPickupAddress
	deliveryAddress := testdata.VolenDeliveryAddress

	testutil.AssertNoErr(t, addressStore.CreateAddress(&pickupAddress))
	testutil.AssertNoErr(t, addressStore.CreateAddress(&deliveryAddress))
}

func initDeliveriesTable(t testing.TB, deliveryStore *models.PgDeliveryStore) {
	delivery := testdata.VolenActiveDelivery

	testutil.AssertNoErr(t, deliveryStore.CreateDelivery(&delivery))
}
