
## CRUD

### R

- 带有(多串条件)的查询。

```go
var (
    ret         []*entity.Progress
    err         error
    filterField = make([]string, 0, len(justCodeVerion))
    cols        = "user_id, application_code, application_version, module_code, module_version"
)

for k := range justCodeVerion {
    filterField = append(filterField, fmt.Sprintf(" ( user_id = %d AND application_code='%s' AND application_version='%s' AND module_code='%s' AND module_version='%s' ) ", justCodeVerion[k].UserId, justCodeVerion[k].ApplicationCode, justCodeVerion[k].ApplicationVersion, justCodeVerion[k].ModuleCode, justCodeVerion[k].ModuleVersion))
}

err = tlr.Xorm().Context(context.Background()).
    Table(entity.TableProgress).
    Select(cols).
    Where(strings.Join(filterField, " OR ")).
    Find(&ret)
if err != nil {
    return nil, err
}
return ret, err
```



- 使用SQL然后再用Find输出数据到ORM

```go
	ts := make([]*entity.Timeline, 0)
	var err error
	var (
		conditionField []string
	)

	if len(applicationVersion) > 0 {
		conditionField = append(conditionField, fmt.Sprintf(" application_version = '%s' ", applicationVersion))
	}
	if len(moduleCode) > 0 {
		conditionField = append(conditionField, fmt.Sprintf(" module_code = '%s' ", moduleCode))
	}
	if len(moduleVersion) > 0 {
		conditionField = append(conditionField, fmt.Sprintf(" module_version = '%s' ", moduleVersion))
	}

	conditionField = append(conditionField, fmt.Sprintf(" user_id = %d AND application_code = '%s' ", UserId, applicationCode))

	var sqlString = "select * from at_2090_timeline where id in ( select max(id) from at_2090_timeline as t where %s group by module_code )"

	err = tlr.Xorm().Context(ctx).
		SQL(fmt.Sprintf(sqlString, strings.Join(conditionField, " AND "))).
		Find(&ts)
	if err != nil {
		return ts, err
	}
	return ts, nil
```


























