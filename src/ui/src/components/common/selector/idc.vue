/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <bk-select
        :disabled="disabled"
        :filterFn="filterFn"
        :filterable="filterable"
        :list="bkIdcList"
        :multiple="multiple"
        :placeholder="placeholder"
        :selected.sync="curSelected"
        :valueKey="valueKey"
        @on-selected="setSelectedData">
        <bk-select-option v-for="(idc, index) in bkIdcList"
            :key="idc['bk_idc_id']"
            :value="idc['bk_idc_id']"
            :label="idc['bk_idc_name']">
        </bk-select-option>
    </bk-select>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            disabled: {
                type: Boolean,
                required: false,
                default: false
            },
            filterFn: {
                type: Function,
                required: false
            },
            filterable: {
                type: Boolean,
                required: false,
                default: false
            },
            multiple: {
                type: Boolean,
                required: false,
                default: false
            },
            placeholder: {
                type: String,
                required: false,
                default: ''
            },
            valueKey: {
                type: String,
                required: false,
                default: 'value'
            },
            selected: {
                required: false
            }
        },
        computed: {
            ...mapGetters(['bkIdcList', 'bkIdcId']),
            curSelected: {
                get () {
                    return this.bkIdcId
                },
                set (val) {
                    this.$store.commit('setBkIdcId', val)
                }
            }
        },
        data () {
            return {
                selectedData: {
                    label: '',
                    value: ''
                },
                selectedIndex: 0
            }
        },
        watch: {
            bkIdcList () {
                this.setSelectedData()
                this.setHeader()
                this.$nextTick(() => {
                    this.$emit('update:selected', this.curSelected)
                    this.$emit('on-selected', this.selectedData, this.selectedIndex)
                })
            },
            curSelected (val) {
                this.setHeader()
                this.$nextTick(() => {
                    this.$emit('update:selected', val)
                    this.$emit('on-selected', this.selectedData, this.selectedIndex)
                })
            }
        },
        methods: {
            ...mapActions(['getBkIdcList']),
            setSelectedData (data, index) {
                console.log('***********')
                console.log(data)
                if (data) {
                    this.selectedData = data
                    this.selectedIndex = index
                } else {
                    /* 用于默认选择时向父组件派发on-selected事件 */
                    let label = ''
                    for (var i = 0; i < this.bkIdcList.length; i++) {
                        if (this.bkIdcList[i]['bk_idc_id'] === this.curSelected) {
                            label = this.bkIdcList[i]['bk_idc_name']
                            index = i
                            break
                        }
                    }
                    this.selectedData = {
                        label: label,
                        value: this.curSelected
                    }
                    this.selectedIndex = index
                }
            },
            setHeader () {
                if (this.$route.meta.setBkIdcId) {
                    this.$axios.defaults.headers.bk_idc_id = this.curSelected
                } else {
                    delete this.$axios.defaults.headers.bk_idc_id
                }
            }
        },
        beforeCreate () {
            delete this.$axios.defaults.headers.bk_idc_id
        },
        created () {
            this.getBkIdcList()
        },
        beforeDestroy () {
            delete this.$axios.defaults.headers.bk_idc_id
        }
    }
</script>