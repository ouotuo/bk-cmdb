<template>
    <div class="transfer-pop" v-show="isShow">
        <div class="transfer-content" ref="drag" v-drag="'#drag'">
            <div class="content-title" id="drag">
                <i class="icon icon-cc-shift mr5"></i>
                机架转移
            </div>
            <div class="content-section clearfix">
                <div class="section-left fl">
                    <div class="section-tree">
                        <v-tree ref="idcTree"
                            :hideRoot="true"
                            :treeData="treeData"
                            :initNode="initNode"
                            @nodeClick="handleNodeClick"
                            @nodeToggle="handleNodeToggle"></v-tree>
                    </div>
                </div>
                <div class="section-right fr">
                    <ul class="selected-list">
                        <li class="selected-item clearfix" v-for="(node, index) in selectedList" :key="index">
                            <span>{{node['bk_inst_name']}}</span>
                            <i class="bk-icon icon-close fr" @click="removeSelected(index)"></i>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="content-footer clearfix">
                <i18n path="Common['已选N项']" tag="div" class="selected-count fl">
                    <span class="color-info" place="N">{{selectedList.length}}</span>
                </i18n>
                <div class="button-group fr">
                    <bk-button type="primary" v-if="isNotModule" v-show="selectedList.length" @click="doTransfer(false)">{{$t('Common[\'确认转移\']')}}</bk-button>
                    <template v-else v-show="selectedList.length">
                        <bk-button type="primary" @click="doTransfer(false)">{{$t('Common[\'覆盖\']')}}</bk-button>
                        <bk-button type="primary" @click="doTransfer(true)">{{$t('Common[\'更新\']')}}</bk-button>
                    </template>
                    <button class="bk-button vice-btn" @click="cancel">{{$t('Common[\'取消\']')}}</button>
                </div>
            </div>
        </div>
    </div>
</template>
<script type="text/javascript">
    import vApplicationSelector from '@/components/common/selector/idc'
    import vTree from '@/components/tree/tree.v2'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            isShow: Boolean,
            chooseId: Array
        },
        data () {
            return {
                bkIdcId: '0',
                treeData: {
                    'default': 0,
                    'bk_obj_id': 'root',
                    'bk_obj_name': 'root',
                    'bk_inst_id': 'root',
                    'bk_inst_name': 'root',
                    'isFolder': true,
                    'child': [{
                        'default': 0,
                        'bk_obj_id': 'source',
                        'bk_obj_name': '空机位',
                        'bk_inst_id': 'source',
                        'bk_inst_name': '空机位',
                        'isFolder': false
                    }]
                },
                initNode: {
                    level: 1,
                    open: true,
                    active: true,
                    'bk_inst_id': 'root'
                },
                activeNode: {},
                activeParentNode: {},
                selectedList: [],
                allowType: ['source', 'pos'],
                isNotModule: true
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            maxExpandedLevel () {
                return this.getLevel(this.treeData) - 1
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow && this.bkIdcId) {
                    this.init()
                } else {
                    this.initNode = {
                        level: 1,
                        open: true,
                        active: true,
                        'bk_inst_id': 'root'
                    }
                    this.$refs.drag.style = ''
                    this.treeData.child.splice(1)
                }
            },
            maxExpandedLevel (level) {
                let treeListEl = this.$refs.idcTree.$el
                if (level > 8) {
                    let width = treeListEl.getBoundingClientRect().width
                    treeListEl.style.minWidth = `${width + 40}px`
                } else {
                    treeListEl.style.minWidth = 'auto'
                }
            }
        },
        methods: {
            handleIdcSelected (data) {
                this.tree.bkIdcName = data.label
            },
            init () {
                this.getTopoTree()
                this.selectedList = []
            },
            getLevel (node) {
                let level = node.level || 1
                if (node.isOpen && node.child && node.child.length) {
                    level = Math.max(level, Math.max.apply(null, node.child.map(childNode => this.getLevel(childNode))))
                }
                return level
            },
            getIdcId (rackNode) {
                for (let i in this.treeData.child) {
                    let idcNode = this.treeData.child[i]
                    for (let j in idcNode.child) {
                        if (rackNode.bk_inst_id === idcNode.child[j].bk_inst_id) {
                            return idcNode
                        }
                    }
                }
            },
            handleNodeClick (activeNode, nodeOptions) {
                this.activeNode = activeNode
                this.activeParentNode = nodeOptions.parent
                if (this.activeNode['bk_obj_id'] === 'pos') {
                    // 获取idcid
                    this.bkIdcId = this.getIdcId(this.activeParentNode).bk_inst_id
                }
                this.checkNode(activeNode)
            },
            handleNodeToggle (isOpen, node, nodeOptions) {
                if (!node.child || !node.child.length) {
                    this.$set(node, 'isLoading', true)
                    this.$axios.get(`idc/inst/child/${this.bkSupplierAccount}/${node['bk_obj_id']}/${this.bkIdcId}/${node['bk_inst_id']}`).then(res => {
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
            checkNode (node) {
                if (this.allowType.indexOf(node['bk_obj_id']) !== -1) {
                    if (node['default'] || node['bk_inst_id'] === 'source') {
                        if (this.selectedList.length && !this.selectedList[0]['default'] && this.selectedList[0]['bk_inst_id'] !== 'source') {
                            this.$bkInfo({
                                title: this.$t('Common[\'转移确认\']', {target: node['bk_inst_name']}),
                                confirmFn: () => {
                                    this.selectedList = [node]
                                }
                            })
                        } else {
                            this.selectedList = [node]
                        }
                        this.isNotModule = true
                    } else {
                        if (this.selectedList.length && (this.selectedList[0]['default'] || this.selectedList[0]['bk_inst_id'] === 'source')) {
                            this.selectedList = []
                        }
                        let isExist = this.selectedList.find(selectedNode => {
                            return selectedNode['bk_obj_id'] === node['bk_obj_id'] && selectedNode['bk_inst_id'] === node['bk_inst_id']
                        })
                        if (!isExist) {
                            this.selectedList.push(node)
                        }
                        this.isNotModule = false
                    }
                }
            },
            doTransfer (type) {
                if (this.selectedList[0]['bk_obj_id'] === 'source') {
                    this.$axios.post('hosts/pos/resource', {
                        'bk_idc_id': this.bkIdcId,
                        'bk_host_id': this.chooseId
                    }).then(res => {
                        if (res.result) {
                            this.$emit('success', res)
                            this.$alertMsg(this.$t('Common[\'转移成功\']'), 'success')
                            this.cancel()
                        } else {
                            if (res.data && res.data['bk_host_id']) {
                                this.$alertMsg(`${res['bk_error_msg']} : ${res.data['bk_host_id']}`)
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        }
                    })
                } else {
                    let posDefault = this.selectedList[0]['default']
                    let transferType = {0: '', 1: 'idle', 2: 'fault'}
                    let url = `hosts/pos/${transferType[posDefault]}`
                    this.$axios.post(url, {
                        'bk_idc_id': this.bkIdcId,
                        'bk_host_id': this.chooseId,
                        'bk_pos_id': this.selectedList.map(node => {
                            return node['bk_inst_id']
                        }),
                        'is_increment': type
                    }).then(res => {
                        if (res.result) {
                            this.$emit('success', res)
                            this.$alertMsg(this.$t('Common[\'转移成功\']'), 'success')
                            this.cancel()
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    }).catch(e => {
                        if (e.response && e.response.status === 403) {
                            this.$alertMsg(this.$t('Common[\'您没有主机转移的权限\']'))
                        }
                    })
                }
            },
            removeSelected (index) {
                this.selectedList.splice(index, 1)
            },
            getTopoInst () {
                return this.$axios.get(`idc/inst/${this.bkSupplierAccount}/0`).then(res => {
                    return res
                })
            },
            getTopoInternal () {
                return this.$axios.get(`idc/internal/${this.bkSupplierAccount}/0`).then(res => {
                    return res
                })
            },
            getTopoTree () {
                return this.$Axios.all([this.getTopoInst()]).then(this.$Axios.spread((instRes) => {
                    if (instRes.result) {
                        for (var i in instRes.data) {
                            this.treeData.child.push(instRes.data[i])
                        }

                        console.log(JSON.stringify(this.treeData))
                    } else {
                        this.$alertMsg(instRes.message)
                    }
                }))
            },
            cancel () {
                this.$emit('update:isShow', false)
            }
        },
        components: {
            vApplicationSelector,
            vTree
        }
    }
</script>
<style lang="scss" scoped>
    .transfer-pop{
        position: fixed;
        width: 100%;
        height: 100%;
        top: 0;
        left: 0;
        z-index: 2000;
        background: rgba(0, 0, 0, 0.6);
    }
    .transfer-content{
        width: 643px;
        height: 588px;
        background: #fff;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        border-radius: 2px;
    }
    .content-title{
        height: 50px;
        background: #f9f9f9;
        color: #333948;
        font-weight: bold;
        line-height: 50px;
        font-size: 14px;
        padding-left: 30px;
        border-bottom: 1px solid #e7e9ef;
        .icon{
            position: relative;
            top: -1px;
        }
    }
    .content-section{
        height: 478px;
        .section-left{
            height: 100%;
            width: 367px;
            border-right: 1px solid #e7e9ef;
        }
        .section-right {
            height: 100%;
            width: 273px;
            padding: 20px 20px 20px 30px;
            overflow: auto;
            @include scrollbar;
        }
    }
    .content-footer{
        height: 60px;
        line-height: 60px;
        background: #f9f9f9;
        border-top: 1px solid #e7e9ef;
    }
    .section-idc{
        height: 64px;
        border-bottom: 1px solid #e7e9ef;
        font-size: 14px;
        padding: 14px 0;
        .idc-label{
            display: inline-block;
            vertical-align: middle;
            padding: 0 20px 0 30px;
        }
        .idc-selector{
            display: inline-block;
            vertical-align: center;
            width: 245px;
        }
    }
    .section-tree{
        height: 413px;
        padding: 10px 5px 0 0;
        overflow: auto;
        @include scrollbar;
    }
    .selected-list{
        color: #3c96ff;
        font-weight: bold;
        font-size: 12px;
        .selected-item{
            padding: 5px 0;
            .icon-close{
                cursor: pointer;
                color: #9196a1;
                &:hover{
                    color: #3c96ff;
                }
            }
        }
    }
    .selected-count{
        padding: 0 0 0 32px;
        .color-info{
            color: #3c96ff;
        }
    }
    .button-group{
        padding: 0 20px 0 0;
        .bk-button{
            margin: 0 5px;
        }
    }
</style>