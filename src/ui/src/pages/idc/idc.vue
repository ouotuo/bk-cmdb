<template>
    <div class="idc-wrapper clearfix">
        <div class="idc-tree-ctn fl">
            <div class="biz-selector-ctn">
                <v-idc-selector :selected.sync="tree.bkIdcId" @on-selected="handleIdcSelected" :filterable="true"></v-idc-selector>
            </div>
            <div class="idc-options-ctn" hidden>
                <i class="idc-option-del icon-cc-del fr" v-if="isShowOptionDel && Object.keys(tree.treeData).length" @click="deleteNode"></i>
            </div>
            <div class="tree-list-ctn">
                <v-tree ref="idcTree"
                    :treeData="tree.treeData"
                    :initNode="tree.initNode"
                    :model="tree.model"
                    @nodeClick="handleNodeClick"
                    @nodeToggle="handleNodeToggle"
                    @addNode="handleAddNode">
                </v-tree>
            </div>
        </div>
        <div style="position: absolute;top: 30px;right: 2.5%;z-index: 9999">
            <bk-button type="primary" class="bk-button main-btn"  @click="handleAddIdc" >添加机房</bk-button>
        </div>
        <div class="idc-view-ctn">
             <bk-tab :active-name="view.tab.active" @tab-changed="tabChanged" class="idc-view-tab">

                <bk-tabpanel name="host" :title="$t('BusinessIdclogy[\'主机调配\']')">
                    <v-hosts ref="hosts"
                        :outerParams="searchParams"
                        :isShowRefresh="true"
                        :outerLoading="tree.loading"
                        :isShowCrossImport="authority['is_host_cross_biz'] && attributeBkObjId === 'pos'"
                        :tableVisible="view.tab.active === 'host'"
                        @handleCrossImport="handleCrossImport">
                        <div slot="filter"></div>
                    </v-hosts>
                </bk-tabpanel>
                <bk-tabpanel name="attribute" :title="$t('BusinessIdclogy[\'节点属性\']')" :show="isShowAttribute">
                    <v-attribute ref="idcAttribute"
                        :bkObjId="attributeBkObjId"
                        :bkIdcId="tree.bkIdcId"
                        :activeNode="tree.activeNode"
                        :activeParentNode="tree.activeParentNode"
                        :formValues="view.attribute.formValues"
                        :type="view.attribute.type"
                        :active="view.tab.active === 'attribute'"
                        :isLoading="view.attribute.isLoading"
                        :editable="tree.bkIdcName !== '蓝鲸'"
                        @submit="submitNode"
                        @delete="deleteNode"
                        @cancel="cancelCreate"></v-attribute>
                </bk-tabpanel>
            </bk-tab>
        </div>
        <bk-dialog :is-show.sync="view.crossImport.isShow" :quick-close="false" :has-header="false" :has-footer="false" :width="700" :padding="0">
            <v-cross-import  slot="content"
                :is-show.sync="view.crossImport.isShow"
                :idcId="tree.bkIdcId"
                :posId="tree.activeNode['bk_inst_id']"
                @handleCrossImportSuccess="rackSearchParams">
            </v-cross-import>
        </bk-dialog>
    </div>
</template>
<script>
    import vIdcSelector from '@/components/common/selector/idc'
    import vTree from '@/components/tree/tree.v2'
    import vHosts from '@/pages/hosts/hosts'
    import vAttribute from './children/attribute'
    import vCrossImport from './children/crossImport'
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                tree: {
                    bkIdcId: -1,
                    bkIdcName: '',
                    treeData: {},
                    model: [],
                    activeNode: {},
                    activeNodeOptions: {},
                    activeParentNode: {},
                    initNode: {},
                    loading: true
                },
                view: {
                    tab: {
                        active: 'host'
                    },
                    attribute: {
                        type: 'update',
                        formValues: {},
                        isLoading: true
                    },
                    crossImport: {
                        isShow: false
                    }
                },
                nodeToggleRecord: {},
                searchParams: null
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            ...mapGetters('navigation', ['authority']),
            /* 获取当前属性表单对应的属性obj_id */
            attributeBkObjId () {
                let bkObjId
                if (this.view.attribute.type === 'create') {
                    console.log(this.view.attribute.IsIdc)
                    if (this.view.attribute.IsIdc) {
                        bkObjId = 'idc'
                    } else {
                        bkObjId = this.tree.model.find(model => {
                            return model['bk_obj_id'] === this.tree.activeNode['bk_obj_id']
                        })['bk_next_obj']
                    }
                } else {
                    bkObjId = this.tree.activeNode['bk_obj_id']
                }
                return bkObjId
            },
            /* 计算是否显示属性修改tab选项卡 */
            isShowAttribute () {
                let isShow = !this.tree.activeNode['default']
                return (isShow || this.view.attribute.type === 'create') && !!Object.keys(this.tree.treeData).length
            },
            /* 计算是否显示添加按钮 */
            isShowOptionAdd () {
                let activeNode = this.tree.activeNode
                return this.tree.model.length && Object.keys(activeNode).length && activeNode['bk_obj_id'] !== 'pos' && !activeNode['default']
            },
            /* 查找当前节点对应的主线拓扑模型 */
            optionModel () {
                return this.tree.model.find(model => {
                    return model['bk_obj_id'] === this.tree.activeNode['bk_obj_id']
                })
            },
            /* 计算是否显示删除按钮 */
            isShowOptionDel () {
                return !this.tree.activeNode['default'] && this.view.tab.active === 'host'
            },
            /* 计算当前树展开的最大层次 */
            maxExpandedLevel () {
                return this.getLevel(this.tree.treeData) - 1
            }
        },
        watch: {
            /* 业务切换，初始化拓扑树 */
            'tree.bkIdcId' (bkIdcId) {
                this.getIdcTree().then(() => {
                    this.tree.initNode = {
                        level: 1,
                        open: true,
                        active: true,
                        bk_inst_id: this.tree.treeData['bk_inst_id']
                    }
                })
            },
            /* 当前节点发生变化且属性修改面板激活时，加载当前节点的具体属性 */
            'tree.activeNode' () {
                if (!this.isShowAttribute) {
                    this.tabChanged('host')
                }
                if (this.view.tab.active === 'attribute') {
                    this.getNodeDetails()
                }
            },
            /* tab选项卡处于切换到属性面板时，加载节点具体属性 */
            'view.tab.active' (activeTab) {
                if (activeTab === 'attribute' && this.view.attribute.type === 'update') {
                    this.getNodeDetails()
                }
            },
            /* 根据当前树节点展开的层级设置横向宽度 */
            maxExpandedLevel (level) {
                let extendLength = 40
                let treeListEl = this.$refs.idcTree.$el
                if (level > 4) {
                    let width = treeListEl.getBoundingClientRect().width
                    treeListEl.style.minWidth = `${width + extendLength}px`
                } else {
                    treeListEl.style.minWidth = 'auto'
                }
            }
        },
        methods: {
            handleIdcSelected (data) {
                this.tree.bkIdcName = data.label
            },
            /* 获取最大展开层级 */
            getLevel (node) {
                let level = node.level
                if (node.isOpen && node.child && node.child.length) {
                    level = Math.max(level, Math.max.apply(null, node.child.map(childNode => this.getLevel(childNode))))
                }
                return level
            },
            /* 获取业务拓扑实例 */
            getIdcInst () {
                return this.$axios.get(`idc/inst/${this.bkSupplierAccount}/${this.tree.bkIdcId}`).then(res => {
                    return res
                }).catch(e => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'您没有当前业务的权限\']'))
                    }
                })
            },
            /* 获取内置业务拓扑 */
            getIdcInternal () {
                return this.$axios.get(`idc/internal/${this.bkSupplierAccount}/${this.tree.bkIdcId}`).then(res => {
                    return res
                }).catch(e => {
                    if (e.response && e.response.status === 403) {
                        this.$alertMsg(this.$t('Common[\'您没有当前业务的权限\']'))
                    }
                })
            },
            /* 获取主线拓扑模型 */
            getIdcModel () {
                this.$axios.get(`idc/model/${this.bkSupplierAccount}`).then(res => {
                    if (res.result) {
                        this.tree.model = res.data
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /* 初始化拓扑树 */
            getIdcTree () {
                this.tree.loading = true
                return this.$Axios.all([this.getIdcInst(), this.getIdcInternal()]).then(this.$Axios.spread((instRes, internalRes) => {
                    if (instRes.result && internalRes.result) {
                        let internalPos = internalRes.data.pos.map(pos => {
                            return {
                                'default': pos['bk_pos_name'] === '空闲机' || pos['bk_pos_name'] === 'idle machine' ? 1 : 2,
                                'bk_obj_id': 'pos',
                                'bk_obj_name': this.$t('Hosts[\'模块\']'),
                                'bk_inst_id': pos['bk_pos_id'],
                                'bk_inst_name': pos['bk_pos_name'],
                                'isFolder': false
                            }
                        })
                       // instRes.data[0]['child'] = internalPos.concat(instRes.data[0]['child'])
                        this.tree.treeData = instRes.data[0]
                    } else {
                        this.$alertMsg(internalRes.result ? instRes.message : internalRes.message)
                    }
                })).then(() => {
                    this.tree.loading = false
                }).catch(() => {
                    this.tree.loading = false
                })
            },
            /* 获取当前节点的具体属性 */
            getNodeDetails () {
                let {
                    bk_inst_id: bkInstId,
                    bk_inst_name: bkInstName,
                    bk_obj_id: bkObjId
                } = this.tree.activeNode
                let url
                let params = {
                    page: {sort: 'id'},
                    fields: [],
                    condition: {}
                }
                if (bkObjId === 'idc') {
                    url = `idc/search/${this.bkSupplierAccount}`
                    params['condition']['bk_idc_id'] = bkInstId
                } else if (bkObjId === 'rack') {
                    url = `rack/search/${this.bkSupplierAccount}/${this.tree.bkIdcId}`
                    params['condition']['bk_rack_id'] = bkInstId
                } else if (bkObjId === 'pos') {
                    url = `pos/search/${this.bkSupplierAccount}/${this.tree.bkIdcId}/${this.tree.activeParentNode['bk_inst_id']}`
                    params['condition']['bk_pos_id'] = bkInstId
                    params['condition']['bk_supplier_account'] = this.bkSupplierAccount
                } else {
                    url = `inst/search/${this.bkSupplierAccount}/${bkObjId}/${bkInstId}`
                }
                this.view.attribute.isLoading = true
                this.$axios.post(url, params).then(res => {
                    if (res.result) {
                        this.view.attribute.formValues = res.data.info[0]
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.view.attribute.isLoading = false
                }).catch(() => {
                    this.view.attribute.isLoading = false
                })
            },
            /* 新增拓扑，切换到属性表单 */
            handleAddNode () {
                this.view.attribute.formValues = {}
                this.view.attribute.isLoading = false
                this.view.attribute.type = 'create'
                this.view.tab.active = 'attribute'
                this.view.attribute.IsIdc = false
            },
            handleAddIdc () {
                this.view.attribute.formValues = {}
                this.view.attribute.isLoading = false
                this.view.attribute.type = 'create'
                this.view.tab.active = 'attribute'
                this.view.attribute.IsIdc = true
            },
            /* 新增拓扑节点/修改拓扑节点 */
            submitNode (formData, originalData) {
                let url
                let method
                let submitType = this.view.attribute.type
                let {
                    bk_inst_id: bkInstId,
                    bk_obj_id: bkObjId
                } = this.tree.activeNode
                if (submitType === 'create') {
                    method = 'post'
                    formData['bk_parent_id'] = bkInstId
                    if (this.attributeBkObjId === 'idc') {
                        delete formData['bk_parent_id']
                        url = `idc/${this.bkSupplierAccount}`
                       // formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else if (this.attributeBkObjId === 'rack') {
                        url = `rack/${this.tree.bkIdcId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else if (this.attributeBkObjId === 'pos') {
                        url = `pos/${this.tree.bkIdcId}/${bkInstId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else {
                        url = `inst/${this.bkSupplierAccount}/${this.attributeBkObjId}`
                        formData['bk_idc_id'] = this.tree.bkIdcId
                    }
                } else if (submitType === 'update') {
                    method = 'put'
                    if (this.attributeBkObjId === 'idc') {
                        url = `idc/${this.bkSupplierAccount}/${this.tree.bkIdcId}`
                       // formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else if (bkObjId === 'rack') {
                        url = `rack/${this.tree.bkIdcId}/${bkInstId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else if (bkObjId === 'pos') {
                        url = `pos/${this.tree.bkIdcId}/${this.tree.activeParentNode['bk_inst_id']}/${bkInstId}`
                        formData['bk_supplier_account'] = this.bkSupplierAccount
                    } else {
                        url = `inst/${this.bkSupplierAccount}/${bkObjId}/${bkInstId}`
                    }
                }
                this.$axios({
                    url: url,
                    method: method,
                    data: formData
                }).then(res => {
                    if (res.result) {
                        this.updateIdcTree(this.view.attribute.type, res.data, formData)
                        this.$alertMsg(submitType === 'create' ? this.$t('Common[\'新建成功\']') : this.$t('Common[\'修改成功\']'), 'success')
                        if (this.view.attribute.type === 'create') {
                            this.view.tab.active = 'host'
                        } else {
                            this.getNodeDetails()
                        }
                        this.view.attribute.type = 'update'
                        this.$refs.idcAttribute.displayType = 'list'
                        if (this.attributeBkObjId === 'idc') {
                            window.location.reload()
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /* 新增、修改拓扑节点成功后更新拓扑树 */
            updateIdcTree (type, response, formData) {
                let node = this.tree.activeNode
                let {
                    bk_next_obj: bkNextObj,
                    bk_next_name: bkNextName,
                    bk_obj_id: bkObjId,
                    bk_obj_name: bkObjName
                } = this.optionModel
                if (type === 'create') {
                    if (node.hasOwnProperty('isFolder')) {
                        node['isFolder'] = true
                    } else {
                        this.$set(node, 'isFolder', true)
                    }
                    node.child = node.child || []
                    node.child.push({
                        'default': 0,
                        'bk_inst_id': bkNextObj === 'rack' ? response['bk_rack_id'] : bkNextObj === 'pos' ? response['bk_pos_id'] : response['bk_inst_id'],
                        'bk_inst_name': bkNextObj === 'rack' ? formData['bk_rack_name'] : bkNextObj === 'pos' ? formData['bk_pos_name'] : formData['bk_inst_name'],
                        'bk_obj_id': bkNextObj,
                        'bk_obj_name': bkNextName,
                        'child': [],
                        'isFolder': false
                    })
                } else if (type === 'update') {
                    node['bk_inst_name'] = bkObjId === 'rack' ? formData['bk_rack_name'] : bkObjId === 'pos' ? formData['bk_pos_name'] : formData['bk_inst_name']
                }
            },
            /* 删除拓扑节点 */
            deleteNode () {
                this.$bkInfo({
                    title: `${this.$t('Common[\'确定删除\']')} ${this.tree.activeNode['bk_inst_name']}?`,
                    content: this.tree.activeNode['bk_obj_id'] === 'pos'
                        ? this.$t('Common["请先转移其下所有的主机"]')
                        : this.$t('Common[\'下属层级都会被删除，请先转移其下所有的主机\']'),
                    confirmFn: () => {
                        let url
                        let {
                            bk_obj_id: bkObjId,
                            bk_inst_id: bkInstId
                        } = this.tree.activeNode
                        if (bkObjId === 'idc') {
                            url = `idc/${this.tree.bkIdcId}`
                        } else if (bkObjId === 'rack') {
                            url = `rack/${this.tree.bkIdcId}/${bkInstId}`
                        } else if (bkObjId === 'pos') {
                            url = `pos/${this.tree.bkIdcId}/${this.tree.activeParentNode['bk_inst_id']}/${bkInstId}`
                        } else {
                            url = `inst/${this.bkSupplierAccount}/${bkObjId}/${bkInstId}`
                        }
                        this.$axios.delete(url).then(res => {
                            if (res.result) {
                                this.view.tab.active = 'host'
                                this.tree.activeParentNode.child.splice(this.tree.activeNodeOptions.index, 1)
                                this.tree.initNode = {
                                    level: 1,
                                    open: true,
                                    active: true,
                                    bk_inst_id: this.tree.treeData['bk_inst_id']
                                }
                                this.$alertMsg(this.$t('Common[\'删除成功\']'), 'success')
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                })
            },
            /* 点击节点，设置查询参数 */
            handleNodeClick (activeNode, nodeOptions) {
                this.$refs.hosts.clearChooseId()
                this.tree.activeNode = activeNode
                this.tree.activeNodeOptions = nodeOptions
                this.tree.activeParentNode = nodeOptions.parent
                this.view.attribute.type = 'update'
                this.rackSearchParams()
            },
            /* node节点展开时，判断是否加载下级节点 */
            handleNodeToggle (isOpen, node, nodeOptions) {
                if (!node.child || !node.child.length) {
                    this.$set(node, 'isLoading', true)
                    this.$axios.get(`idc/inst/child/${this.bkSupplierAccount}/${node['bk_obj_id']}/${this.tree.bkIdcId}/${node['bk_inst_id']}`).then(res => {
                        if (res.result) {
                            let child = res['data'][0]['child']
                            if (Array.isArray(child) && child.length) {
                                node.child = child
                            } else {
                                this.$set(node, 'isFolder', false)
                            }
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                        node.isLoading = false
                    })
                }
            },
            /* 设置主机查询参数 */
            rackSearchParams () {
                let params = {
                    'bk_idc_id': this.tree.bkIdcId,
                    condition: []
                }
                let activeNodeObjId = this.tree.activeNode['bk_obj_id']
                if (activeNodeObjId === 'pos' || activeNodeObjId === 'rack') {
                    params.condition.push({
                        'bk_obj_id': activeNodeObjId,
                        fields: [],
                        condition: [{
                            field: activeNodeObjId === 'pos' ? 'bk_pos_id' : 'bk_rack_id',
                            operator: '$eq',
                            value: this.tree.activeNode['bk_inst_id']
                        }]
                    })
                } else if (activeNodeObjId !== 'idc') {
                    params.condition.push({
                        'bk_obj_id': 'object',
                        fields: [],
                        condition: [{
                            field: 'bk_inst_id',
                            operator: '$eq',
                            value: this.tree.activeNode['bk_inst_id']
                        }]
                    })
                }
                let defaultObj = ['host', 'pos', 'rack', 'idc']
                defaultObj.forEach(id => {
                    if (!params.condition.some(({bk_obj_id: bkObjId}) => bkObjId === id)) {
                        params.condition.push({
                            'bk_obj_id': id,
                            fields: [],
                            condition: []
                        })
                    }
                })
                this.searchParams = params
            },
            handleCrossImport () {
                this.view.crossImport.isShow = true
            },
            tabChanged (active) {
                this.view.tab.active = active
                this.view.attribute.formValues = {}
                if (active === 'host') {
                    this.view.attribute.type = 'update'
                }
            },
            cancelCreate () {
                this.tabChanged('host')
            }
        },
        created () {
            this.getIdcModel()
        },
        components: {
            vIdcSelector,
            vTree,
            vCrossImport,
            vHosts,
            vAttribute
        }
    }
</script>
<style lang="scss" scoped>
    .idc-wrapper{
        height: 100%;
        .idc-tree-ctn{
            width: 280px;
            height: 100%;
            border-right: 1px solid #e7e9ef;
            background: #fafbfd;
        }
        .idc-view-ctn{
            height: 100%;
            overflow: hidden;
            padding: 0 20px;
        }
    }
    .biz-selector-ctn{
        padding: 20px;
    }
    .idc-options-ctn{
        height: 44px;
        line-height: 44px;
        background: #f9f9f9;
        padding: 0 10px;
        .idc-option-add{
            width: 90px;
            height: 24px;
            font-size: 12px;
            line-height: 22px;
            border-radius: 2px;
            background: #fff;
            padding: 0 5px;
            margin: 0 5px;
            @include ellipsis;
        }
        .idc-option-del{
            font-size: 12px;
            margin: 16px 0 0 0;
            cursor: pointer;
        }
    }
    .idc-view-tab{
        height: 100%;
        border: none;
    }
    .tree-list-ctn{
        padding: 0 0 0 20px;
        max-height: calc(100% - 120px);
        overflow: auto;
        @include scrollbar;
    }
</style>
<style lang="scss">
    .bk-tab2.idc-view-tab{
        .bk-tab2-head{
            height: 80px;
            .tab2-nav-item{
                height: 79px;
                line-height: 79px;
            }
        }
        .bk-tab2-content{
            height: calc(100% - 80px);
        }
    }
</style>
