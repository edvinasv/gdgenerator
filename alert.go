package main

import (
	"strconv"
)

func generateAlerts(ddata DashboardData) *AlertGroups {

	alertGroups := new(AlertGroups)
	alertGroup := new(AlertGroup)

	if ddata.DryRun { // || ddata.Disabled || ddata.AlertingDisabled {
		return alertGroups
	}

	alertGroup.Name = ddata.Name
	alertGroup.Interval = "1m"

	for _, row := range ddata.Rows {
		row.ItemSettings.inheritItemSettings(ddata.DefaultItemSettings)
		if row.ItemSettings.Disabled || row.ItemSettings.AlertingDisabled {
			continue
		}
		readRow(alertGroup, row)
	}
	alertGroups.Groups = append(alertGroups.Groups, *alertGroup)
	return alertGroups
}

func readRow(
	alertGroup *AlertGroup,
	row DataRow) {

	for _, column := range row.Columns {
		column.ItemSettings.inheritItemSettings(row.ItemSettings)
		if column.ItemSettings.Disabled || column.ItemSettings.AlertingDisabled {
			continue
		}
		readRowColumn(alertGroup, column)
	}
}

func readRowColumn(
	alertGroup *AlertGroup,
	column DataRowColumn) {

	for _, group := range column.Groups {
		group.ItemSettings.inheritItemSettings(column.ItemSettings)
		if group.ItemSettings.Disabled || group.ItemSettings.AlertingDisabled {
			continue
		}
		readRowColumnGroup(alertGroup, group)
	}
}

func readRowColumnGroup(
	alertGroup *AlertGroup,
	group DataRowColumnGroup) {

	for _, panel := range group.Items {
		panel.configureItemSettings(group.ItemSettings)
		if panel.Disabled || panel.AlertingDisabled {
			continue
		}
		generateAlert(alertGroup, panel)
	}
}

func generateAlert(
	alertGroup *AlertGroup,
	panel DataRowColumnGroupItem) {
	for _, r := range panel.Ranges {
		if len(r.Severity) == 0 ||
			(panel.CriticalDisabled && r.Severity == "critical") ||
			(panel.WarningDisabled && r.Severity == "warning") {
			continue
		}
		alertGroup.Rules = append(alertGroup.Rules, renderAlert(panel, r))
	}
}

func renderAlert(panel DataRowColumnGroupItem, r Range) AlertRule {
	var expr string
	if r.Min == r.Max {
		expr = "(" + panel.Expr + ") == 0"
	} else {
		expr = "(" + panel.Expr + ") >= " +
			strconv.FormatFloat(r.Min, 'f', -1, 64) +
			" and (" + panel.Expr + ") <= " +
			strconv.FormatFloat(r.Max, 'f', -1, 64)
	}
	return AlertRule{
		Alert:       panel.Title,
		Expr:        expr,
		ForDuration: "5m",
		Labels: AlertLabels{
			Severity: r.Severity,
			Service:  panel.Service},
		Annotations: AlertAnnotations{
			Summary: panel.Title + " became " + r.Severity},
	}
}
