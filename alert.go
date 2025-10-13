package main

func generateAlerts(ddata DashboardData) *AlertGroups {

	alertGroups := new(AlertGroups)
	alertGroup := new(AlertGroup)

	if ddata.DryRun { // || ddata.Disabled || ddata.AlertingDisabled {
		return alertGroups
	}

	alertGroup.Name = ddata.Name

	alertGroup.Interval = "1m"

	for _, row := range ddata.Rows {
//		if row.Disabled || row.AlertingDisabled {
//			continue
//		}
//		if len(row.Service) == 0 {
//			row.Service = ddata.Service
//		}
		readRow(alertGroup, row)
	}
	alertGroups.Groups = append(alertGroups.Groups, *alertGroup)
	return alertGroups
}

func readRow(
	alertGroup *AlertGroup,
	row DataRow) {

	for _, column := range row.Columns {
//		if column.Disabled || column.AlertingDisabled {
//			continue
//		}
//		if len(column.Service) == 0 {
//			column.Service = row.Service
//		}
		readRowColumn(alertGroup, column)
	}
}

func readRowColumn(
	alertGroup *AlertGroup,
	column DataRowColumn) {

	for _, group := range column.Groups {
//		if group.Disabled || group.AlertingDisabled {
//			continue
//		}
//		if len(group.Service) == 0 {
//			group.Service = column.Service
//		}
		readRowColumnGroup(alertGroup, group)
	}
}

func readRowColumnGroup(
	alertGroup *AlertGroup,
	group DataRowColumnGroup) {

	for _, panel := range group.Items {
//		if panel.Disabled || panel.AlertingDisabled {
//			continue
//		}
//		if len(panel.Service) == 0 {
//			panel.Service = group.Service
//		}
		generateAlert(alertGroup, panel)
	}
}

func generateAlert(
	alertGroup *AlertGroup,
	panel DataRowColumnGroupItem) {

	alertGroup.Rules = append(alertGroup.Rules, AlertRule{
		Alert: panel.Title,
		Expr:  "(" + panel.Expr + ") == 0",
		//ForDuration: "0m",
		//KeepFiringForDuration: "0m",
		Labels: AlertLabels{
			Severity: "critical",
			Service:  panel.Service},
		Annotations: AlertAnnotations{
			Summary: panel.Title + " became critical"},
	})

	alertGroup.Rules = append(alertGroup.Rules, AlertRule{
		Alert: panel.Title,
		Expr:  "(" + panel.Expr + ") > 0 and (" + panel.Expr + ") < 1",
		//ForDuration: "0m",
		//KeepFiringForDuration: "0m",
		Labels: AlertLabels{
			Severity: "warning",
			Service:  panel.Service},
		Annotations: AlertAnnotations{
			Summary: panel.Title + " became warning"},
	})
}
