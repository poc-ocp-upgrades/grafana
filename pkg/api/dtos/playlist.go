package dtos

type PlaylistDashboard struct {
	Id	int64	`json:"id"`
	Slug	string	`json:"slug"`
	Title	string	`json:"title"`
	Uri	string	`json:"uri"`
	Order	int	`json:"order"`
}
type PlaylistDashboardsSlice []PlaylistDashboard

func (slice PlaylistDashboardsSlice) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(slice)
}
func (slice PlaylistDashboardsSlice) Less(i, j int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return slice[i].Order < slice[j].Order
}
func (slice PlaylistDashboardsSlice) Swap(i, j int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	slice[i], slice[j] = slice[j], slice[i]
}
