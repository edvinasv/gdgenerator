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
		row.ItemSettings.inheritItemSettings(ddata.DefaultItemSettings)
		row.GroupSettings.inheritGroupSettings(ddata.DefaultGroupSettings)
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
		column.ItemSettings.inheritItemSettings(row.ItemSettings)
		column.GroupSettings.inheritGroupSettings(row.GroupSettings)

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
		group.ItemSettings.inheritItemSettings(column.ItemSettings)
		group.Settings.inheritGroupSettings(column.GroupSettings)

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
	gridPos.W = group.ItemSettings.Width
	gridPos.H = group.ItemSettings.Height
	colMarginPos.H = groupRefPos.H
	for index, panel := range group.Items {

		panel.configureItemSettings(group.ItemSettings)

		gridPos.Y = groupRefPos.Y + uint32(index)/group.Settings.Columns*group.ItemSettings.Height
		gridPos.X = groupRefPos.X + (uint32(index))%(group.Settings.Columns)*group.ItemSettings.Width
		switch panel.Type {
		case "stat":
			generateItemStatPanel(builder, panel, gridPos)
		case "timeseries":
			generateItemTimeseriesPanel(builder, panel, gridPos)
		}
	}
	colMarginPos.Y = max(colMarginPos.Y, groupRefPos.Y+uint32(len(group.Items)+1)/group.Settings.Columns*group.ItemSettings.Height)
	colMarginPos.X = max(colMarginPos.X, groupRefPos.X+min(group.Settings.Columns, uint32(len(group.Items)))*group.ItemSettings.Width+group.Settings.Spacing)
}

func generateItemStatPanel(
	builder *dashboard.DashboardBuilder,
	item DataRowColumnGroupItem,
	gridPos *dashboard.GridPos) {

	noDataColor := "grey"
	disabledColor := "rgb(50,50,50)"

	var vMappings []dashboard.ValueMapping

	for index, r := range item.Ranges {
		vMappings = append(vMappings, dashboard.ValueMapping{
			RangeMap: &dashboard.RangeMap{
				Type: "range",
				Options: dashboard.DashboardRangeMapOptions{
					From: cog.ToPtr(r.Min),
					To:   cog.ToPtr(r.Max),
					Result: dashboard.ValueMappingResult{
						Color: cog.ToPtr(r.Color),
						Index: cog.ToPtr(int32(index)),
					},
				},
			},
		})
	}

	panel := stat.NewPanelBuilder()

	color := noDataColor

	if item.Disabled {
		color = disabledColor
		item.Title = " "
		item.Descr = []string{}
		item.Link = ""

	} else {
		panel = panel.Mappings(vMappings)
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
		//TODO
		//Min(0).
		//Max(500).
		Datasource(
			dashboard.DataSourceRef{
				Type: cog.ToPtr("prometheus"),
				Uid:  cog.ToPtr(item.DatasourceUid),
			},
		).
		GridPos(*gridPos)

	builder.WithPanel(panel)
}
