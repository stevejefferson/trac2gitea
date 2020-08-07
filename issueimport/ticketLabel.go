package issueimport

func (importer *Importer) importTicketLabels(issueID int64, component string, severity string, priority string, version string, resolution string, typ string) {
	var lbl string
	if component != "" {
		lbl = "Component / " + component
		importer.giteaAccessor.AddIssueLabel(issueID, lbl)
	}

	if severity != "" {
		lbl = "Severity / " + severity
		importer.giteaAccessor.AddIssueLabel(issueID, lbl)
	}

	if priority != "" {
		lbl = "Priority / " + priority
		importer.giteaAccessor.AddIssueLabel(issueID, lbl)
	}

	if version != "" {
		lbl = "Version / " + version
		importer.giteaAccessor.AddIssueLabel(issueID, lbl)
	}

	if resolution != "" {
		lbl = "Resolution / " + resolution
		importer.giteaAccessor.AddIssueLabel(issueID, lbl)
	}

	if typ != "" {
		lbl = "Type / " + typ
		importer.giteaAccessor.AddIssueLabel(issueID, lbl)
	}
}
