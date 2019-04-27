package sqlstore

import (
	"fmt"
	"time"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bus.AddHandler("sql", GetOrgQuotaByTarget)
	bus.AddHandler("sql", GetOrgQuotas)
	bus.AddHandler("sql", UpdateOrgQuota)
	bus.AddHandler("sql", GetUserQuotaByTarget)
	bus.AddHandler("sql", GetUserQuotas)
	bus.AddHandler("sql", UpdateUserQuota)
	bus.AddHandler("sql", GetGlobalQuotaByTarget)
}

type targetCount struct{ Count int64 }

func GetOrgQuotaByTarget(query *m.GetOrgQuotaByTargetQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quota := m.Quota{Target: query.Target, OrgId: query.OrgId}
	has, err := x.Get(&quota)
	if err != nil {
		return err
	} else if !has {
		quota.Limit = query.Default
	}
	rawSql := fmt.Sprintf("SELECT COUNT(*) as count from %s where org_id=?", dialect.Quote(query.Target))
	resp := make([]*targetCount, 0)
	if err := x.SQL(rawSql, query.OrgId).Find(&resp); err != nil {
		return err
	}
	query.Result = &m.OrgQuotaDTO{Target: query.Target, Limit: quota.Limit, OrgId: query.OrgId, Used: resp[0].Count}
	return nil
}
func GetOrgQuotas(query *m.GetOrgQuotasQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quotas := make([]*m.Quota, 0)
	sess := x.Table("quota")
	if err := sess.Where("org_id=? AND user_id=0", query.OrgId).Find(&quotas); err != nil {
		return err
	}
	defaultQuotas := setting.Quota.Org.ToMap()
	seenTargets := make(map[string]bool)
	for _, q := range quotas {
		seenTargets[q.Target] = true
	}
	for t, v := range defaultQuotas {
		if _, ok := seenTargets[t]; !ok {
			quotas = append(quotas, &m.Quota{OrgId: query.OrgId, Target: t, Limit: v})
		}
	}
	result := make([]*m.OrgQuotaDTO, len(quotas))
	for i, q := range quotas {
		rawSql := fmt.Sprintf("SELECT COUNT(*) as count from %s where org_id=?", dialect.Quote(q.Target))
		resp := make([]*targetCount, 0)
		if err := x.SQL(rawSql, q.OrgId).Find(&resp); err != nil {
			return err
		}
		result[i] = &m.OrgQuotaDTO{Target: q.Target, Limit: q.Limit, OrgId: q.OrgId, Used: resp[0].Count}
	}
	query.Result = result
	return nil
}
func UpdateOrgQuota(cmd *m.UpdateOrgQuotaCmd) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return inTransaction(func(sess *DBSession) error {
		quota := m.Quota{Target: cmd.Target, OrgId: cmd.OrgId}
		has, err := sess.Get(&quota)
		if err != nil {
			return err
		}
		quota.Updated = time.Now()
		quota.Limit = cmd.Limit
		if !has {
			quota.Created = time.Now()
			if _, err := sess.Insert(&quota); err != nil {
				return err
			}
		} else {
			if _, err := sess.ID(quota.Id).Update(&quota); err != nil {
				return err
			}
		}
		return nil
	})
}
func GetUserQuotaByTarget(query *m.GetUserQuotaByTargetQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quota := m.Quota{Target: query.Target, UserId: query.UserId}
	has, err := x.Get(&quota)
	if err != nil {
		return err
	} else if !has {
		quota.Limit = query.Default
	}
	rawSql := fmt.Sprintf("SELECT COUNT(*) as count from %s where user_id=?", dialect.Quote(query.Target))
	resp := make([]*targetCount, 0)
	if err := x.SQL(rawSql, query.UserId).Find(&resp); err != nil {
		return err
	}
	query.Result = &m.UserQuotaDTO{Target: query.Target, Limit: quota.Limit, UserId: query.UserId, Used: resp[0].Count}
	return nil
}
func GetUserQuotas(query *m.GetUserQuotasQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	quotas := make([]*m.Quota, 0)
	sess := x.Table("quota")
	if err := sess.Where("user_id=? AND org_id=0", query.UserId).Find(&quotas); err != nil {
		return err
	}
	defaultQuotas := setting.Quota.User.ToMap()
	seenTargets := make(map[string]bool)
	for _, q := range quotas {
		seenTargets[q.Target] = true
	}
	for t, v := range defaultQuotas {
		if _, ok := seenTargets[t]; !ok {
			quotas = append(quotas, &m.Quota{UserId: query.UserId, Target: t, Limit: v})
		}
	}
	result := make([]*m.UserQuotaDTO, len(quotas))
	for i, q := range quotas {
		rawSql := fmt.Sprintf("SELECT COUNT(*) as count from %s where user_id=?", dialect.Quote(q.Target))
		resp := make([]*targetCount, 0)
		if err := x.SQL(rawSql, q.UserId).Find(&resp); err != nil {
			return err
		}
		result[i] = &m.UserQuotaDTO{Target: q.Target, Limit: q.Limit, UserId: q.UserId, Used: resp[0].Count}
	}
	query.Result = result
	return nil
}
func UpdateUserQuota(cmd *m.UpdateUserQuotaCmd) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return inTransaction(func(sess *DBSession) error {
		quota := m.Quota{Target: cmd.Target, UserId: cmd.UserId}
		has, err := sess.Get(&quota)
		if err != nil {
			return err
		}
		quota.Updated = time.Now()
		quota.Limit = cmd.Limit
		if !has {
			quota.Created = time.Now()
			if _, err := sess.Insert(&quota); err != nil {
				return err
			}
		} else {
			if _, err := sess.ID(quota.Id).Update(&quota); err != nil {
				return err
			}
		}
		return nil
	})
}
func GetGlobalQuotaByTarget(query *m.GetGlobalQuotaByTargetQuery) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rawSql := fmt.Sprintf("SELECT COUNT(*) as count from %s", dialect.Quote(query.Target))
	resp := make([]*targetCount, 0)
	if err := x.SQL(rawSql).Find(&resp); err != nil {
		return err
	}
	query.Result = &m.GlobalQuotaDTO{Target: query.Target, Limit: query.Default, Used: resp[0].Count}
	return nil
}
