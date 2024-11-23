package hserver

import "testing"

var html = `
<!DOCTYPE html>
<html>
<service mame="human" description="人员信息" url="/api/v1.0/tenants/{tenantId}/human/" >
    <runtime name="init">
        let ctx = context.background();
        $self.schema = newSchema(humanSchemas.default);
        $self.dao = db.model('human');
        $self.dao.table(ctx, $self.schema).create(ctx);
    </runtime>

    <request type="post" name="create" url="" description="创建人员信息">
        <params url="create.json"></params>
        <call ref="dao" method="create" params="params.body.list"></call>
    </request>

    <request type="put" name="update" url="" description="创建人员信息">
        <params url="update.json"></params>
        <call ref="dao" method="update" params="params.body.list"></call>
    </request>

    <request type="put" name="saveList" url="save-list" description="创建人员信息">
        <params url="saveList.json"></params>
        <runtime type="text/javascript">
            $self.dao.update($params.body.list);
        </runtime>
    </request>

    <request type="delete" name="deleteById" url="{id}" description="创建人员信息">
        <params url="deleteById.json"></params>
        <runtime type="text/javascript">
            $self.dao.update($params.body.list);
        </runtime>
    </request>

    <request type="get" name="findById" url="{id}" description="创建人员信息">
        <params url="findById.json"></params>
        <call ref="dao" method="getById">
            <param name="tenandId" type="string">param.tenandId</param>
            <param id="id" type="string">param.id</param>
        </call>
    </request>

    <request type="get" name="findPaging" url="" description="创建人员信息">
        <params url="findPaging.json"></params>
        <call ref="dao" method="getById">
            <param name="tenandId" type="string">param.tenandId</param>
            <param id="id" type="string">param.id</param>
        </call>
    </request>

</service>
</html>
`

func Test_parseHTML(t *testing.T) {
	server, err := NewService([]byte(html))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", server)
}
