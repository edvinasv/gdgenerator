package main

func (s *ItemSettings) inheritItemSettings(parentSettings ItemSettings) *ItemSettings {
	if len(s.Type) == 0 {
		s.Type = parentSettings.Type
	}
	if parentSettings.Disabled {
		s.Disabled = true
	}
	if parentSettings.AlertingDisabled {
		s.AlertingDisabled = true
	}
	if parentSettings.CriticalDisabled {
		s.CriticalDisabled = true
	}
	if parentSettings.WarningDisabled {
		s.WarningDisabled = true
	}
	if len(s.Service) == 0 {
		s.Service = parentSettings.Service
	}
	if len(s.DatasourceUid) == 0 {
		s.DatasourceUid = parentSettings.DatasourceUid
	}
	if len(s.Ranges) == 0 {
		s.Ranges = parentSettings.Ranges
	}
	if !(s.Width > 0) {
		s.Width = parentSettings.Width
	}
	if !(s.Height > 0) {
		s.Height = parentSettings.Height
	}
	return s
}

func (s *DataRowColumnGroupItem) configureItemSettings(parentSettings ItemSettings) *DataRowColumnGroupItem {
	if len(s.Type) == 0 {
		s.Type = parentSettings.Type
	}
	if parentSettings.Disabled {
		s.Disabled = true
	}
	if parentSettings.AlertingDisabled {
		s.AlertingDisabled = true
	}
	if parentSettings.CriticalDisabled {
		s.CriticalDisabled = true
	}
	if parentSettings.WarningDisabled {
		s.WarningDisabled = true
	}
	if len(s.Service) == 0 {
		s.Service = parentSettings.Service
	}
	if len(s.DatasourceUid) == 0 {
		s.DatasourceUid = parentSettings.DatasourceUid
	}
	if len(s.Ranges) == 0 {
		s.Ranges = parentSettings.Ranges
	}
	return s
}

func (g *GroupSettings) inheritGroupSettings(parentSettings GroupSettings) *GroupSettings {
	if !(g.Columns > 0) {
		g.Columns = parentSettings.Columns
	}
	//TODO: overwrite spacing > 0
	if !(g.Spacing > 0) {
		g.Spacing = parentSettings.Spacing
	}
	return g
}
