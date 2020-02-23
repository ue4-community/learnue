## 通过Docker安装Solr

参考这个[链接](https://github.com/docker-solr/docker-solr)

```yaml
version: '3'
services:
  solr:
    image: solr:8
    ports:
     - "8983:8983"
    volumes:
      - ./data/solrData:/var/solr
    command:
      - solr-precreate
      - learnue
```

## 关闭托管的Schema定义并使用手动编辑Schema模式
[参考](http://lucene.apache.org/solr/guide/8_4/schema-factory-definition-in-solrconfig.html#switching-from-managed-schema-to-manually-edited-schema-xml)

如果直接把config下的schema.xml和solrconfig.xml复制到solr_data/learnue/conf/就跳过下面的步骤


1.打开solr_data/learnue/conf/solrconfig.xml

2.找到updateRequestProcessorChain,并将autoCreateFields:true设置为false,processors里去掉add-schema-fields

3.启用手动编辑模式
  ```xml
<config>
  ...
  <schemaFactory class="ClassicIndexSchemaFactory"/>
</config>  
```
4.自定义配置,把managed-schema同目录下复制一份取名为schema.xml,然后打开，在根节点内容最后加入下面的内容。
```xml
<!--开始自定义编辑-->
<field name="objid" type="pint" indexed="false" stored="true" required="true" multiValued="false" />
<field name="objtype" type="pint" indexed="true" stored="true" required="true" multiValued="false" />
<field name="title" type="text_general" indexed="true" stored="true" required="true" multiValued="false" />
<field name="author" type="string" indexed="true" stored="true" />
<field name="uid" type="pint" indexed="false" stored="true" />
<field name="pub_time" type="string" indexed="true" stored="true" />
<field name="content" type="text_general" indexed="true" stored="true" multiValued="false"  />
<field name="tags" type="text_general" indexed="true" stored="true" multiValued="false" />
<field name="viewnum" type="pint" indexed="true" stored="true" />
<field name="cmtnum" type="pint" indexed="true" stored="true" />
<field name="likenum" type="pint" indexed="true" stored="true" />
<field name="nid" type="pint" indexed="false" stored="true" />
<field name="lastreplyuid" type="pint" indexed="false" stored="true" />
<field name="lastreplytime" type="string" indexed="false" stored="true" />
<field name="top" type="pint" indexed="true" stored="true" />
<field name="created_at" type="string" indexed="false" stored="true" />
<field name="updated_at" type="string" indexed="false" stored="true" />
<field name="sort_time" type="string" indexed="true" stored="true" />

<copyField source="author" dest="author_str" maxChars="256"/>
<copyField source="title" dest="title_str" maxChars="256"/>
<copyField source="content" dest="content_str" maxChars="256"/>
<copyField source="tags" dest="tags_str" maxChars="256"/>
<!--结束自定义编辑-->
```



