<template>
    <div class="history-wrapper">
        <div class="history-table">
            <v-table
                :loading="table.isLoading"
                :header="table.header"
                :list="table.list"
                :pagination="table.pagination"
                :defaultSort="table.defaultSort"
                :wrapperMinusHeight="250"
                @handlePageChange="setCurrentPage"
                @handleSizeChange="setPageSize"
                >
            </v-table>
        </div>
    </div>
</template>
<script>
    import vTable from '@/components/table/table'
    import vMemberSelector from '@/components/common/selector/member'
    import vHistoryDetails from '@/components/history/details'
    import moment from 'moment'
    export default {
        props: {
            active: {
                type: Boolean,
                default: false
            },
            type: {
                type: String,
                default: 'inst' // inst | host
            },
            innerIP: String,
            instId: Number
        },
        data () {
            return {
                details: {
                    isShow: false,
                    data: null,
                    clickoutside: true
                },
                filter: {
                    date: [],
                    user: ''
                },
                table: {
                    header: [{
                        id: 'bk_host_manageip',
                        name: this.$t("HostResourcePool['交换机IP']")
                    }, {
                        id: 'bk_bind_ip',
                        name: this.$t("HostResourcePool['绑定IP']")
                    }, {
                        id: 'bk_port_name',
                        name: this.$t("HostResourcePool['端口']")
                    }, {
                        id: 'bk_mac_add',
                        name: this.$t("HostResourcePool['物理地址']")
                    }],
                    list: [],
                    pagination: {
                        current: 1,
                        count: 0,
                        size: 10
                    },
                    defaultSort: '-bk_bind_ip',
                    sort: '-bk_bind_ip',
                    isLoading: false
                }
            }
        },
        computed: {
            searchParams () {
                let params = {
                    condition: {
                        bk_host_manageip: this.innerIP
                    },
                    limit: this.table.pagination.size,
                    start: (this.table.pagination.current - 1) * this.table.pagination.size,
                    sort: this.table.sort
                }

                return params
            }
        },
        beforeMount () {
            this.filter.date = [`${this.initDate.start} 00:00:00`, `${this.initDate.end} 23:59:59`]
        },
        watch: {
            active (active) {
                if (active) {
                    this.getHistory()
                } else {
                    let $dateRangePicker = this.$refs.dateRangePicker
                    $dateRangePicker.selectedDateView = `${this.initDate.start} - ${this.initDate.end}`
                    $dateRangePicker.selectedDateRange = [this.initDate.start, this.initDate.end]
                    $dateRangePicker.selectedDateRangeTmp = [this.initDate.start, this.initDate.end]
                    this.filter.date = [`${this.initDate.start} 00:00:00`, `${this.initDate.end} 23:59:59`]
                    this.filter.user = ''
                }
            }
        },
        methods: {
            closeDetails () {
                if (!this.details.clickoutside) {
                    this.details.isShow = false
                    this.details.data = null
                }
            },
            getHistory () {
                this.table.isLoading = true
                this.$axios.post('host/get/port', this.searchParams).then(res => {
                    console.log(res)
                    if (res.result) {
                        this.table.list = res.data.info
                        this.table.pagination.count = res.data.count
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.table.isLoading = false
                }).catch(() => {
                    this.table.isLoading = false
                })
            },
            setFilterDate (oldDate, newDate) {
                if (newDate) {
                    newDate = newDate.split(' - ')
                    newDate[0] = `${newDate[0]} 00:00:00`
                    newDate[1] = `${newDate[1]} 23:59:59`
                    this.filter.date = newDate
                }
            },
            setPageSize (size) {
                this.table.pagination.size = size
                this.getHistory()
            },
            setCurrentPage (current) {
                this.table.pagination.current = current
                this.getHistory()
            }
        },
        components: {
            vTable,
            vMemberSelector,
            vHistoryDetails
        }
    }
</script>

<style lang="scss" scoped>
    .history-wrapper{
        position: relative;
    }
</style>
