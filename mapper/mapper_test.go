package mapper

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-renderer/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{}
	mdl := model.Page{}

	Convey("test filter flex overview maps correctly", t, func() {
		m := CreateFilterFlexOverview(ctx, mdl, cfg)
		So(m.BetaBannerEnabled, ShouldBeTrue)
		So(m.Type, ShouldEqual, "filter-flex-overview")
		So(m.Metadata.Title, ShouldEqual, "Filter flex overview")
	})
}
