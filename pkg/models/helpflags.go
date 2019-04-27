package models

type HelpFlags1 uint64

const (
	HelpFlagGettingStartedPanelDismissed	HelpFlags1	= 1 << iota
	HelpFlagDashboardHelp1
)

func (f HelpFlags1) HasFlag(flag HelpFlags1) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f&flag != 0
}
func (f *HelpFlags1) AddFlag(flag HelpFlags1) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*f |= flag
}
func (f *HelpFlags1) ClearFlag(flag HelpFlags1) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*f &= ^flag
}
func (f *HelpFlags1) ToggleFlag(flag HelpFlags1) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*f ^= flag
}

type SetUserHelpFlagCommand struct {
	HelpFlags1	HelpFlags1
	UserId		int64
}
