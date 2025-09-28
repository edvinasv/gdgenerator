package main

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/text"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
	"strings"
)

func generateDashboard(ddata DashboardData) *dashboard.DashboardBuilder {

	marginPos := dashboard.NewGridPos()
	marginPos.X = 0
	marginPos.Y = 0

	builder := dashboard.NewDashboardBuilder(ddata.Name).
		Uid(ddata.Uid).
		Tags([]string{"generic-dashboard", "generated"}).
		Time(ddata.TimeRange.From, ddata.TimeRange.To).
		Timezone(common.TimeZoneBrowser).
		//Editable().
		Readonly().
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		Refresh(ddata.Refresh)

	if len(ddata.Banner) > 0 {
		generateBanner(builder, ddata.Banner, marginPos)
	}

	if ddata.DryRun {
		return builder
	}

	for _, row := range ddata.Rows {
		if len(row.DatasourceUid) == 0 {
			row.DatasourceUid = ddata.DatasourceUid
		}
		if len(row.Type) == 0 {
			row.Type = ddata.Type
		}
		if ddata.Disabled {
			row.Disabled = true
		}
		generateRow(builder, row, marginPos)
	}
	return builder
}

func generateBanner(builder *dashboard.DashboardBuilder, banner string, marginPos *dashboard.GridPos) {
	gridPos := dashboard.NewGridPos()
	gridPos.H = 2
	gridPos.W = 24
	gridPos.Y = marginPos.Y
	gridPos.X = marginPos.X

	builder.WithPanel(
		text.NewPanelBuilder().
			Content(banner).
			Mode("html").
			GridPos(*gridPos),
	)

	marginPos.Y += gridPos.H

}

func generateRow(
	builder *dashboard.DashboardBuilder,
	row DataRow,
	marginPos *dashboard.GridPos) {

	colRefPos := dashboard.NewGridPos()    // column reference position
	colMarginPos := dashboard.NewGridPos() // column reference position

	marginPos.X = 0
	colRefPos.Y = marginPos.Y + 1
	colRefPos.X = 0
	colMarginPos.Y = colRefPos.Y
	colMarginPos.X = colRefPos.X

	builder.WithRow(
		dashboard.NewRowBuilder(row.Title).GridPos(*marginPos))

	for _, column := range row.Columns {
		if len(column.DatasourceUid) == 0 {
			column.DatasourceUid = row.DatasourceUid
		}
		if len(column.Type) == 0 {
			column.Type = row.Type
		}
		if row.Disabled {
			column.Disabled = true
		}
		generateRowColumn(builder, column, colRefPos, colMarginPos)
		colRefPos.X = colMarginPos.X
		marginPos.Y = max(marginPos.Y, colMarginPos.Y)
		marginPos.X = colMarginPos.X
		colMarginPos.Y = colRefPos.Y
	}
}

func generateRowColumn(
	builder *dashboard.DashboardBuilder,
	column DataRowColumn,
	colRefPos *dashboard.GridPos,
	colMarginPos *dashboard.GridPos) {

	groupRefPos := dashboard.NewGridPos()

	groupRefPos.Y = colRefPos.Y
	groupRefPos.X = colRefPos.X

	for _, group := range column.Groups {
		if len(group.DatasourceUid) == 0 {
			group.DatasourceUid = column.DatasourceUid
		}
		if len(group.Type) == 0 {
			group.Type = column.Type
		}
		if column.Disabled {
			group.Disabled = true
		}
		generateRowColumnGroup(builder, group, groupRefPos, colMarginPos)
		groupRefPos.Y = colMarginPos.Y
	}
}

func generateRowColumnGroup(
	builder *dashboard.DashboardBuilder,
	group DataRowColumnGroup,
	groupRefPos *dashboard.GridPos,
	colMarginPos *dashboard.GridPos) {

	gridPos := dashboard.NewGridPos()
	gridPos.W = group.Width
	gridPos.H = group.Height
	colMarginPos.H = groupRefPos.H

	for index, panel := range group.Items {
		if len(panel.DatasourceUid) == 0 {
			panel.DatasourceUid = group.DatasourceUid
		}
		if len(panel.Type) == 0 {
			panel.Type = group.Type
		}
		if group.Disabled {
			panel.Disabled = true
		}
		gridPos.Y = groupRefPos.Y + uint32(index)/group.Columns*group.Height
		gridPos.X = groupRefPos.X + (uint32(index))%(group.Columns)*group.Width
		switch panel.Type {
		case "stat":
			generateItemStatPanel(builder, panel, gridPos)
		case "timeseries":
			generateItemTimeseriesPanel(builder, panel, gridPos)
		}
	}
	colMarginPos.Y = max(colMarginPos.Y, groupRefPos.Y+uint32(len(group.Items)+1)/group.Columns*group.Height)
	colMarginPos.X = max(colMarginPos.X, groupRefPos.X+min(group.Columns, uint32(len(group.Items)))*group.Width+group.Spacing)
}

func generateItemStatPanel(
	builder *dashboard.DashboardBuilder,
	item DataRowColumnGroupItem,
	gridPos *dashboard.GridPos) {

	noDataColor := "grey"
	disabledColor := "rgb(50,50,50)"

	vmResultGreen := dashboard.ValueMappingResult{
		Color: cog.ToPtr("green"),
		Index: cog.ToPtr(int32(0)),
	}
	vmResultRed := dashboard.ValueMappingResult{
		Color: cog.ToPtr("red"),
		Index: cog.ToPtr(int32(1)),
	}
	vmResultYellow := dashboard.ValueMappingResult{
		Color: cog.ToPtr("yellow"),
		Index: cog.ToPtr(int32(2)),
	}

	rangeMapOption1 := dashboard.DashboardRangeMapOptions{
		From:   cog.ToPtr(float64(0)),
		To:     cog.ToPtr(float64(1)),
		Result: vmResultYellow,
	}
	valueMap1 := dashboard.ValueMap{
		Type: "value",
		Options: map[string]dashboard.ValueMappingResult{
			"1": vmResultGreen,
			"0": vmResultRed,
		},
	}
	rangeMap1 := dashboard.RangeMap{
		Type:    "range",
		Options: rangeMapOption1,
	}

	vMappingValues := dashboard.ValueMapping{
		ValueMap: &valueMap1,
	}
	vMappingRange := dashboard.ValueMapping{
		RangeMap: &rangeMap1,
	}

	panel := stat.NewPanelBuilder()

	color := noDataColor

	if item.Disabled {
		color = disabledColor
		item.Title = " "
		item.Descr = []string{}
		item.Link = ""

	} else {
		panel = panel.Mappings([]dashboard.ValueMapping{vMappingValues, vMappingRange})
	}

	panel = panel.ColorScheme(dashboard.NewFieldColorBuilder().
		Mode("fixed").
		FixedColor(color))

	if len(item.Link) > 0 {
		link := dashboard.NewDashboardLinkBuilder("Link").
			Type("link").
			Url(item.Link).
			TargetBlank(true)

		panel = panel.Links([]cog.Builder[dashboard.DashboardLink]{link})
	}

	panel = panel.
		Description(strings.Join(item.Descr, "<br>\n")).
		WithTarget(
			prometheus.NewDataqueryBuilder().
				Expr(item.Expr).
				Range().
				Format(prometheus.PromQueryFormatTimeSeries).
				LegendFormat(item.Title),
		).
		TextMode("name").
		ColorMode("background_solid").
		GraphMode("none").
		Min(item.Min).
		Max(item.Max).
		Datasource(
			dashboard.DataSourceRef{
				Type: cog.ToPtr("prometheus"),
				Uid:  cog.ToPtr(item.DatasourceUid),
			},
		).
		GridPos(*gridPos)

	builder.WithPanel(panel)
}

func generateItemTimeseriesPanel(
	builder *dashboard.DashboardBuilder,
	item DataRowColumnGroupItem,
	gridPos *dashboard.GridPos) {

	panel := timeseries.NewPanelBuilder()

	if item.Disabled {
		item.Title = " "
		item.Descr = []string{}
		item.Link = ""

	}

	if len(item.Link) > 0 {
		link := dashboard.NewDashboardLinkBuilder("Link").
			Type("link").
			Url(item.Link)

		panel = panel.Links([]cog.Builder[dashboard.DashboardLink]{link})
	}

	panel = panel.
		Description(strings.Join(item.Descr, "<br>\n")).
		WithTarget(
			prometheus.NewDataqueryBuilder().
				Expr(item.Expr).
				Range().
				Format(prometheus.PromQueryFormatTimeSeries).
				LegendFormat(item.Title),
		).
		Min(item.Min).
		Max(item.Max).
		Datasource(
			dashboard.DataSourceRef{
				Type: cog.ToPtr("prometheus"),
				Uid:  cog.ToPtr(item.DatasourceUid),
			},
		).
		GridPos(*gridPos)

	builder.WithPanel(panel)
}
