{%- macro createField(prop, ui) %}
    {%- if ui.Widget=='input' %}
        <ui5-input id="{{prop.Name}}" class="w100" ></ui5-input>
    {%- elif ui.Widget=='select' %}
        <ui5-select  id="{{prop.Name}}" class="w100">
            <ui5-option icon="iphone">Phone</ui5-option>
            <ui5-option icon="ipad">Tablet</ui5-option>
            <ui5-option icon="laptop" selected="" focused="">Desktop</ui5-option>
        </ui5-select>
    {%- elif ui.Widget=='textarea' %}
        <ui5-textarea id="{{prop.Name}}" class="w100 h100"></ui5-textarea>
    {%- elif ui.Widget=='date' %}
        <ui5-date-picker  id="{{prop.Name}}" class="w100"></ui5-date-picker>
    {%- elif ui.Widget=='daterange' %}
        <ui5-daterange-picker  id="{{prop.Name}}" class="w100"></ui5-daterange-picker> 
    {%- elif ui.Widget=='datetime' %}
        <ui5-datetime-picker  id="{{prop.Name}}" class="w100"></ui5-datetime-picker>
    {%- elif ui.Widget=='color' %}
        <ui5-color-picker  id="{{prop.Name}}" class="w100"></ui5-color-picker>
    {%- elif ui.Widget=='check' %}
        <ui5-checkbox  id="{{prop.Name}}"  text="Basic" class="w100"></ui5-checkbox>
    {%- elif ui.Widget=='link' %}
        <ui5-link  id="{{prop.Name}}" class="w100" href="https://www.sap.com" target="_blank">
            Standard Link
        </ui5-link>
    {%- elif ui.Widget=='rating' %}
        <ui5-rating-indicator  id="{{prop.Name}}" value="5" max="7" required="" class="w100">
        </ui5-rating-indicator>
    {%- elif ui.Widget=='radio' %}
        <ui5-radio-button id="{{prop.Name}}" name="{{prop.Name}}" text="Option 1" checked>
        <ui5-radio-button></ui5-radio-button></ui5-radio-button>
        <ui5-radio-button name="GroupA" text="Option 2">
        <ui5-radio-button></ui5-radio-button></ui5-radio-button>
        <ui5-radio-button name="GroupA" text="Option 2">
        <ui5-radio-button></ui5-radio-button></ui5-radio-button>
    {%- elif ui.Widget=='multi-select' %}
    <ui5-multi-combobox id="{{prop.Name}}" >
        <ui5-mcb-item text="Albania"></ui5-mcb-item>
        <ui5-mcb-item selected="" text="Argentina"></ui5-mcb-item>
        <ui5-mcb-item text="Bulgaria"></ui5-mcb-item>
        <ui5-mcb-item text="Denmark"></ui5-mcb-item>
        <ui5-mcb-item text="England"></ui5-mcb-item>
        <ui5-mcb-item text="Germany"></ui5-mcb-item>
        <ui5-mcb-item text="Philippines"></ui5-mcb-item>
        <ui5-mcb-item text="Portugal"></ui5-mcb-item>
        <ui5-mcb-item text="The United Kingdom of Great Britain and Northern Ireland" ></ui5-mcb-item>
    </ui5-multi-combobox>
    {%- else %}
        {{ui.Widget}}
    {%- endif %}
{%- endmacro %}

{%- macro colStyle(col) %}
    style="flex:0 0 calc({{col.Span*100/24}}% - 10px);"
{%- endmacro %}

{%- macro rowStyle(row) %}
    {%- if row.Height %}
    style="min-height:{{row.Height}};"
    {%- endif %}
{%- endmacro %}

<style>
    .grid {
        display: flex;
        flex-direction: column;
        padding: 0px;
        min-width: 100%;
    }
    .row {
        display: flex;
        flex-direction: row;
        flex: 1;
        border: 0px solid #de0e0e;
        padding: 0px;
    }
    .col {
        display: flex;
        padding: 0px;
        margin: 5px;
    }
    .label {
        max-width:{{grid.GetTitleWidth()}};
        min-width:{{grid.GetTitleWidth()}};
    }
    .field {
        width: calc(100% - {{grid.GetTitleWidth()}});
    }
    .w100 {
        width: 100%;
    }
    .h100{
        height: 100%;
    }
</style>
<script>
   var schema = {{schema | JSON}}
</script>
<ui5-panel>
    <div slot="header" class="header">
        {{schema.Title}}
    </div>
    <div class="grid">
        {%- for row in grid.Rows %}
        <div class="row" {{rowStyle(row)}}>
            {%- for col in row.Cols %}
            <div class="col" {{colStyle(col)}}>
                <div class="label">
                    <ui5-label id="{{col.Property.Name}}-ref" required="{{col.Required}}" for="{{col.Property.Name}}">
                        {{col.Property.GetTitle()}}
                    </ui5-label>
                </div>
                <div class="field" field="{{col.Property.Name}}">
                    {{createField(col.Property, col.UiProperty)}}
                </div>
            </div>
            {%- endfor %}
        </div>
        {%- endfor %}
    </div>
</ui5-panel>
    