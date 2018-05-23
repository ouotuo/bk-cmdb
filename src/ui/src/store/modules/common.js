/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import Vue from 'vue'
import Cookies from 'js-cookie'
import {$Axios, $axios} from '@/api/axios'

let alertMsg = Vue.prototype.$alertMsg

const state = {
    bkSupplierAccount: '0',
    timezoneList: [],   // 时区列表
    biz: {     // 业务列表
        list: [],
        selected: parseInt(Cookies.get('bk_biz_id')) || -1
    },
    memberList: [],
    isAdmin: window.isAdmin === '1',
    usercustom: {},          // 用户字段配置
    globalLoading: false,
    memberLoading: false,
    language: Cookies.get('blueking_language') || 'zh_CN'
}

const getters = {
    bkSupplierAccount: state => state.bkSupplierAccount,
    bkBizId: state => state.biz.selected,
    bkBizList: state => state.biz.list,
    memberList: state => state.memberList,
    isAdmin: state => state.isAdmin,
    navigation: state => state.navigation,
    timezoneList: state => state.timezoneList,
    usercustom: state => state.usercustom,
    globalLoading: state => state.globalLoading,
    memberLoading: state => state.memberLoading,
    language: state => state.language
}

const actions = {
    getBkBizList ({commit, state}) {
        $axios.post(`biz/search/${state.bkSupplierAccount}`, {fields: ['bk_biz_id', 'bk_biz_name']}).then((res) => {
            if (res.result) {
                if (res.data.info && res.data.info.length) {
                    commit('setBkBizList', res.data.info)
                    if (state.biz.selected === -1) { // 如果未选择过，则选中第一个业务
                        commit('setBkBizId', state['biz']['list'][0]['bk_biz_id'])
                    } else { // 如果已经选择过，则需判断缓存的已选择业务是否被删除
                        let isExist = false
                        state.biz.list.map((biz) => {
                            if (state.biz.selected === biz['bk_biz_id']) {
                                isExist = true
                            }
                        })
                        if (!isExist) {
                            commit('setBkBizId', state.biz.list[0]['bk_biz_id'])
                        }
                    }
                } else {
                    commit('setBkBizList', [])
                    commit('setBkBizId', -1)
                }
            } else {
                alertMsg(res['bk_error_msg'])
            }
        })
    },
    getMemberList ({commit, state}, type) {
        state.memberLoading = true
        let baseURL = $axios.defaults.baseURL
        $axios.defaults.baseURL = window.siteUrl
        $axios.get(`/user/list?_t=${(new Date()).getTime()}`, { type }).then((res) => {
            if (res.result) {
                commit('setMemberList', res.data)
            } else {
                alertMsg(res['bk_error_msg'])
            }
            state.memberLoading = false
        }).catch(() => {
            state.memberLoading = false
        })
        $axios.defaults.baseURL = baseURL
    }
}

const mutations = {
    setLang (state, language) {
        state.language = language
    },
    setBkBizList (state, list) {
        state.biz.list = list
    },
    setBkBizId (state, selected) {
        Cookies.set('bk_biz_id', selected, { expires: 30, path: '' })
        state.biz.selected = selected
    },
    deleteApplication (state, appId) {
        let applicationList = state.application.list
        for (let i = 0; i < applicationList.length; i++) {
            if (applicationList[i]['ApplicationID'] === appId) {
                applicationList.splice(i, 1)
                if (state.application.selected === appId) {
                    // 如果删除的业务是已选中的，则重新设置当前选中的业务
                    if (applicationList.length) {
                        Cookies.set('selectedApplicationId', applicationList[0]['ApplicationID'], { expires: 30, path: '' })
                        state.application.selected = applicationList[0]['ApplicationID']
                    } else {
                        Cookies.set('selectedApplicationId', '', { expires: 30, path: '' })
                        state.application.selected = ''
                    }
                }
                break
            }
        }
    },
    setMemberList (state, memberList) {
        state.memberList = memberList
    },
    setTimezoneList (state, timezoneList) {
        state.timezoneList = timezoneList
    },
    setUsercustom (state, usercustom) {
        state.usercustom = usercustom
    },
    setGlobalLoading (state, isLoading) {
        state.globalLoading = isLoading
    }
}

export default {
    state,
    getters,
    actions,
    mutations
}
