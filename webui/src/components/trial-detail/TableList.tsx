import * as React from 'react';
import axios from 'axios';
import ReactEcharts from 'echarts-for-react';
import { connect } from 'react-redux';
import { Row, Table, Button, Modal, Checkbox, Popconfirm, Select, message, Icon, Tabs } from 'antd';
const { TabPane } = Tabs;
const Option = Select.Option;
const CheckboxGroup = Checkbox.Group;
import { MANAGER_IP, trialJobStatus, COLUMN, COLUMN_INDEX, COLUMNPro } from '../../static/const';
import { convertDuration, intermediateGraphOption, filterByStatus } from '../../static/function';
import { TableObj, TrialJob } from '../../static/interface';
import OpenRow from '../public-child/OpenRow';
import Compare from '../Modal/Compare';
import IntermediateVal from '../public-child/IntermediateVal'; // table default metric column
import MonacoHTML from '../public-child/MonacoEditor';
import '../../static/style/search.scss';
require('../../static/style/tableStatus.css');
require('../../static/style/logPath.scss');
require('../../static/style/search.scss');
require('../../static/style/table.scss');
require('../../static/style/button.scss');
require('../../static/style/openRow.scss');
const echarts = require('echarts/lib/echarts');
require('echarts/lib/chart/line');
require('echarts/lib/component/tooltip');
require('echarts/lib/component/title');
echarts.registerTheme('my_theme', {
    color: '#3c8dbc'
});

interface TableListProps {
    entries: number;
    tableSource: Array<TableObj>;
    updateList: Function;
    platform: string;
    logCollection: boolean;
    isMultiPhase: boolean;
    port: string;
    addColumn: Function;
    ExperimentID:string;
}

interface TableListState {
    intermediateOption: object;
    modalVisible: boolean;
    logVisible: boolean;
    isObjFinal: boolean;
    isShowColumn: boolean;
    columnSelected: Array<string>; // user select columnKeys
    selectRows: Array<TableObj>;
    isShowCompareModal: boolean;
    selectedRowKeys: string[] | number[];
    intermediateData: Array<object>; // a trial's intermediate results (include dict)
    intermediateId: string;
    intermediateOtherKeys: Array<string>;
    triallog: string;
    errorlog: string;
    port: string;
}

interface ColumnIndex {
    name: string;
    index: number;
}

class TableList extends React.Component<TableListProps, TableListState> {

    public _isMounted = false;
    public intervalTrialLog = 10;
    public _trialId: string;
    public tables: Table<TableObj> | null;

    constructor(props: TableListProps) {
        super(props);
    
        this.state = {
            intermediateOption: {},
            modalVisible: false,
            logVisible: false,
            isObjFinal: false,
            isShowColumn: false,
            isShowCompareModal: false,
            columnSelected: COLUMN,
            selectRows: [],
            selectedRowKeys: [], // close selected trial message after modal closed
            intermediateData: [],
            intermediateId: '',
            intermediateOtherKeys: [],
            port: this.props.port,
            triallog: '',
            errorlog: ''
        };
    }
    showlogModal = (id: string) => {
        axios(`${MANAGER_IP}/log-data/${id}`, {
            method: 'GET',
            headers:{'upstream': this.state.port, 'ExperimentID': this.props.ExperimentID}
        }).then(res => {
            if (res.status === 200){
                if (this._isMounted){
                    console.log(res)
                    this.setState({
                        triallog:res.data.trial_log,
                        errorlog:res.data.error_log
                    })
                }
            }
        })
        if (this._isMounted) {
            this.setState({
                logVisible: true
            });
        }
    }
    showIntermediateModal = (id: string) => {

        axios(`${MANAGER_IP}/metric-data/${id}`, {
            method: 'GET',
            headers:{'upstream': this.state.port, 'ExperimentID': this.props.ExperimentID}
        })
            .then(res => {
                if (res.status === 200) {
                    const intermediateArr: number[] = [];
                    // support intermediate result is dict because the last intermediate result is
                    // final result in a succeed trial, it may be a dict.
                    // get intermediate result dict keys array
                    let otherkeys: Array<string> = ['default'];
                    if (res.data.length !== 0) {
                        otherkeys = Object.keys(JSON.parse(res.data[0].data));
                    }
                    // intermediateArr just store default val
                    Object.keys(res.data).map(item => {
                        const temp = JSON.parse(res.data[item].data);
                        if (typeof temp === 'object') {
                            intermediateArr.push(temp.default);
                        } else {
                            intermediateArr.push(temp);
                        }
                    });
                    const intermediate = intermediateGraphOption(intermediateArr, id);
                    if (this._isMounted) {
                        this.setState(() => ({
                            intermediateData: res.data, // store origin intermediate data for a trial
                            intermediateOption: intermediate,
                            intermediateOtherKeys: otherkeys,
                            intermediateId: id
                        }));
                    }
                }
            });
        if (this._isMounted) {
            this.setState({
                modalVisible: true
            });
        }
    }

    selectOtherKeys = (value: string) => {

        const isShowDefault: boolean = value === 'default' ? true : false;
        const { intermediateData, intermediateId } = this.state;
        const intermediateArr: number[] = [];
        // just watch default key-val
        if (isShowDefault === true) {
            Object.keys(intermediateData).map(item => {
                const temp = JSON.parse(intermediateData[item].data);
                if (typeof temp === 'object') {
                    intermediateArr.push(temp[value]);
                } else {
                    intermediateArr.push(temp);
                }
            });
        } else {
            Object.keys(intermediateData).map(item => {
                const temp = JSON.parse(intermediateData[item].data);
                if (typeof temp === 'object') {
                    intermediateArr.push(temp[value]);
                }
            });
        }
        const intermediate = intermediateGraphOption(intermediateArr, intermediateId);
        // re-render
        if (this._isMounted) {
            this.setState(() => ({
                intermediateOption: intermediate
            }));
        }
    }

    hideIntermediateModal = () => {
        if (this._isMounted) {
            this.setState({
                modalVisible: false
            });
        }
    }
    hidelogModal = () =>{
        if (this._isMounted) {
            this.setState({
                logVisible: false
            });
        }
    }
    hideShowColumnModal = () => {
        if (this._isMounted) {
            this.setState({
                isShowColumn: false
            });
        }
    }
    trialLogHTML = () => {
        return (
            <div>
                <span>Trial Log</span>
            </div>
        );
    }
    errorLogHTML = () => {
        return (
            <div>
                <span>Error Log</span>
            </div>
        );
    }
    // click add column btn, just show the modal of addcolumn
    addColumn = () => {
        // show user select check button
        if (this._isMounted) {
            this.setState({
                isShowColumn: true
            });
        }
    }

    // checkbox for coloumn
    selectedColumn = (checkedValues: Array<string>) => {
        // 7: because have seven common column, "Intermediate count" is not shown by default
        let count = 7;
        const want: Array<object> = [];
        const finalKeys: Array<string> = [];
        const wantResult: Array<string> = [];
        Object.keys(checkedValues).map(m => {
            switch (checkedValues[m]) {
                case 'Trial No.':
                case 'ID':
                case 'Duration':
                case 'Status':
                case 'Operation':
                case 'Default':
                case 'Intermeidate count':
                    break;
                default:
                    finalKeys.push(checkedValues[m]);
            }
        });

        Object.keys(finalKeys).map(n => {
            want.push({
                name: finalKeys[n],
                index: count++
            });
        });

        Object.keys(checkedValues).map(item => {
            const temp = checkedValues[item];
            Object.keys(COLUMN_INDEX).map(key => {
                const index = COLUMN_INDEX[key];
                if (index.name === temp) {
                    want.push(index);
                }
            });
        });

        want.sort((a: ColumnIndex, b: ColumnIndex) => {
            return a.index - b.index;
        });

        Object.keys(want).map(i => {
            wantResult.push(want[i].name);
        });

        if (this._isMounted) {
            this.setState(() => ({ columnSelected: wantResult }));
        }
    }

    openRow = (record: TableObj) => {
        const { platform, logCollection, isMultiPhase } = this.props;
        return (
            <OpenRow
                trainingPlatform={platform}
                record={record}
                logCollection={logCollection}
                multiphase={isMultiPhase}
            />
        );
    }

    fillSelectedRowsTostate = (selected: number[] | string[], selectedRows: Array<TableObj>) => {
        if (this._isMounted === true) {
            this.setState(() => ({ selectRows: selectedRows, selectedRowKeys: selected }));
        }
    }
    killJob = (id: string) => {
        axios(`${MANAGER_IP}/kill-jobs/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json;charset=utf-8',
                'upstream': this.state.port,
                'ExperimentID': this.props.ExperimentID
            }
        }) 
        .then(res => {
            if (res.status === 200) {
                message.destroy();
                message.success('Cancel the job successfully');
                // render the table
            } else {
                message.error('fail to cancel the job');
            }
        })
        .catch(error => {
            if (error.response.status === 500) {
                if (error.response.data.error) {
                    message.error(error.response.data.error);
                } else {
                    message.error('500 error, fail to cancel the job');
                }
            }
        });
    }
    // open Compare-modal
    compareBtn = () => {

        const { selectRows } = this.state;
        if (selectRows.length === 0) {
            alert('Please select datas you want to compare!');
        } else {
            if (this._isMounted === true) {
                this.setState({ isShowCompareModal: true });
            }
        }
    }
    // close Compare-modal
    hideCompareModal = () => {
        // close modal. clear select rows data, clear selected track
        if (this._isMounted) {
            this.setState({ isShowCompareModal: false, selectedRowKeys: [], selectRows: [] });
        }
    }

    componentDidMount() {
        this._isMounted = true;
    }

    componentWillUnmount() {
        this._isMounted = false;
    }

    render() {

        const { entries, tableSource } = this.props;
        const { intermediateOption, modalVisible, isShowColumn, columnSelected, logVisible, 
            selectRows, isShowCompareModal, selectedRowKeys, intermediateOtherKeys } = this.state;
        const rowSelection = {
            selectedRowKeys: selectedRowKeys,
            onChange: (selected: string[] | number[], selectedRows: Array<TableObj>) => {
                this.fillSelectedRowsTostate(selected, selectedRows);
            }
        };
        let showTitle = COLUMNPro;
        let bgColor = '';
        const trialJob: Array<TrialJob> = [];
        const showColumn: Array<object> = [];
        const heights: number = window.innerHeight - 48; // padding top and bottom
        // only succeed trials have final keys
        if (tableSource.filter(filterByStatus).length >= 1) {
            const temp = tableSource.filter(filterByStatus)[0].acc;
            if (temp !== undefined && typeof temp === 'object') {
                if (this._isMounted) {
                    // concat default column and finalkeys
                    const item = Object.keys(temp);
                    // item: ['default', 'other-keys', 'maybe loss']
                    if (item.length > 1) {
                        const want: Array<string> = [];
                        item.forEach(value => {
                            if (value !== 'default') {
                                want.push(value);
                            }
                        });
                        showTitle = COLUMNPro.concat(want);
                    }
                }
            }
        }
        trialJobStatus.map(item => {
            trialJob.push({
                text: item,
                value: item
            });
        });
        Object.keys(columnSelected).map(key => {
            const item = columnSelected[key];
            switch (item) {
                case 'Trial No.':
                    showColumn.push({
                        title: 'Trial No.',
                        dataIndex: 'sequenceId',
                        key: 'sequenceId',
                        width: 120,
                        className: 'tableHead',
                        sorter: (a: TableObj, b: TableObj) => (a.sequenceId as number) - (b.sequenceId as number)
                    });
                    break;
                case 'ID':
                    showColumn.push({
                        title: 'ID',
                        dataIndex: 'id',
                        key: 'id',
                        width: 60,
                        className: 'tableHead leftTitle',
                        // the sort of string
                        sorter: (a: TableObj, b: TableObj): number => a.id.localeCompare(b.id),
                        render: (text: string, record: TableObj) => {
                            return (
                                <div>{record.id}</div>
                            );
                        }
                    });
                    break;
                case 'Duration':
                    showColumn.push({
                        title: 'Duration',
                        dataIndex: 'duration',
                        key: 'duration',
                        width: 100,
                        // the sort of number
                        sorter: (a: TableObj, b: TableObj) => (a.duration as number) - (b.duration as number),
                        render: (text: string, record: TableObj) => {
                            let duration;
                            if (record.duration !== undefined) {
                                // duration is nagative number(-1) & 0-1
                                if (record.duration > 0 && record.duration < 1 || record.duration < 0) {
                                    duration = `${record.duration}s`;
                                } else {
                                    duration = convertDuration(record.duration);
                                }
                            } else {
                                duration = 0;
                            }
                            return (
                                <div className="durationsty"><div>{duration}</div></div>
                            );
                        },
                    });
                    break;
                case 'Status':
                    showColumn.push({
                        title: 'Status',
                        dataIndex: 'status',
                        key: 'status',
                        width: 150,
                        className: 'tableStatus',
                        render: (text: string, record: TableObj) => {
                            bgColor = record.status;
                            return (
                                <span className={`${bgColor} commonStyle`}>{record.status}</span>
                            );
                        },
                        filters: trialJob,
                        onFilter: (value: string, record: TableObj) => {
                            return record.status.indexOf(value) === 0;
                        },
                        // onFilter: (value: string, record: TableObj) => record.status.indexOf(value) === 0,
                        sorter: (a: TableObj, b: TableObj): number => a.status.localeCompare(b.status)
                    });
                    break;
                case 'Intermeidate count':
                    showColumn.push({
                        title: 'Intermediate count',
                        dataIndex: 'progress',
                        key: 'progress',
                        width: 86,
                        render: (text: string, record: TableObj) => {
                            return (
                                <span>{`#${record.description.intermediate.length}`}</span>
                            );
                        },
                    });
                    break;
                case 'Default':
                    showColumn.push({
                        title: 'Default metric',
                        className: 'leftTitle',
                        dataIndex: 'acc',
                        key: 'acc',
                        width: 120,
                        sorter: (a: TableObj, b: TableObj) => {
                            const aa = a.description.intermediate;
                            const bb = b.description.intermediate;
                            if (aa !== undefined && bb !== undefined) {
                                return aa[aa.length - 1] - bb[bb.length - 1];
                            } else {
                                return NaN;
                            }
                        },
                        render: (text: string, record: TableObj) => {
                            return (
                                <IntermediateVal record={record} />
                            );
                        }
                    });
                    break;
                case 'Operation':
                    showColumn.push({
                        title: 'Operation',
                        dataIndex: 'operation',
                        key: 'operation',
                        width: 120,
                        render: (text: string, record: TableObj) => {
                            let trialStatus = record.status;
                            const flag: boolean = (trialStatus === 'RUNNING') ? false : true;
                            return (
                                <Row id="detail-button">
                                    {/* see intermediate result graph */}
                                    <Button
                                        type="primary"
                                        className="common-style"
                                        onClick={this.showIntermediateModal.bind(this, record.id)}
                                        title="Intermediate"
                                    >
                                        <Icon type="line-chart" />
                                    </Button>
                                    <Button
                                        type="primary"
                                        className="common-style"
                                        onClick={this.showlogModal.bind(this, record.id)}
                                        title="log"
                                    >
                                        <Icon type="file-text" theme="twoTone" />
                                    </Button>
                                    {/* kill job */}
                                    <Popconfirm
                                        title="Are you sure to cancel this trial?"
                                        onConfirm={() => this.killJob(record.id)}
                                    >
                                        <Button
                                            type="default"
                                            disabled={flag}
                                            className="margin-mediate special"
                                            title="kill"
                                        >
                                            <Icon type="stop" />
                                        </Button>
                                    </Popconfirm>
                                </Row>
                            );
                        },
                    });
                    break;

                case 'Intermediate result':
                    showColumn.push({
                        title: 'Intermediate result',
                        dataIndex: 'intermediate',
                        key: 'intermediate',
                        width: '16%',
                        render: (text: string, record: TableObj) => {
                            return (
                                <Button
                                    type="primary"
                                    className="tableButton"
                                    onClick={this.showIntermediateModal.bind(this, record.id)}
                                >
                                    Intermediate
                                </Button>
                            );
                        },
                    });
                    break;
                default:
                    showColumn.push({
                        title: item,
                        dataIndex: item,
                        key: item,
                        width: 150,
                        render: (text: string, record: TableObj) => {
                            const temp = record.acc;
                            let decimals = 0;
                            let other = '';
                            if (temp !== undefined) {
                                if (temp[item].toString().indexOf('.') !== -1) {
                                    decimals = temp[item].toString().length - temp[item].toString().indexOf('.') - 1;
                                    if (decimals > 6) {
                                        other = `${temp[item].toFixed(6)}`;
                                    } else {
                                        other = temp[item].toString();
                                    }
                                }
                            } else {
                                other = '--';
                            }
                            return (
                                <div>{other}</div>
                            );
                        }
                    });
            }
        });
        
        return (
            <Row className="tableList">
                <div id="tableList">
                    <Table
                        ref={(table: Table<TableObj> | null) => this.tables = table}
                        columns={showColumn}
                        rowSelection={rowSelection}
                        expandedRowRender={this.openRow}
                        dataSource={tableSource}
                        className="commonTableStyle"
                        pagination={{ pageSize: entries }}
                    />
                    <Modal
                        title="Trial log"
                        visible={logVisible}
                        width="80%"
                        destroyOnClose={true}
                        onCancel={this.hidelogModal}
                        footer={null}
                    >
                        <Tabs type="line" defaultActiveKey={"triallog"} style={{ height: heights}} >
                            <TabPane tab={this.trialLogHTML()}  key="triallog" >
                                <MonacoHTML content={this.state.triallog} loading={false} />
                            </TabPane>
                            <TabPane tab={this.errorLogHTML()} key="errorlog" >
                                <MonacoHTML content={this.state.errorlog} loading={false} />
                            </TabPane>
                        </Tabs>
                    </Modal>
                    {/* Intermediate Result Modal */}
                    <Modal
                        title="Intermediate result"
                        visible={modalVisible}
                        onCancel={this.hideIntermediateModal}
                        footer={null}
                        destroyOnClose={true}
                        width="80%"
                    >
                        {
                            intermediateOtherKeys.length > 1
                                ?
                                <Row className="selectKeys">
                                    <Select
                                        className="select"
                                        defaultValue="default"
                                        onSelect={this.selectOtherKeys}
                                    >
                                        {
                                            Object.keys(intermediateOtherKeys).map(item => {
                                                const keys = intermediateOtherKeys[item];
                                                return <Option value={keys} key={item}>{keys}</Option>;
                                            })
                                        }
                                    </Select>

                                </Row>
                                :
                                <div />
                        }
                        <ReactEcharts
                            option={intermediateOption}
                            style={{
                                width: '100%',
                                height: 0.7 * window.innerHeight
                            }}
                            theme="my_theme"
                        />
                    </Modal>
                </div>
                {/* Add Column Modal */}
                <Modal
                    title="Table Title"
                    visible={isShowColumn}
                    onCancel={this.hideShowColumnModal}
                    footer={null}
                    destroyOnClose={true}
                    width="40%"
                >
                    <CheckboxGroup
                        options={showTitle}
                        defaultValue={columnSelected}
                        onChange={this.selectedColumn}
                        className="titleColumn"
                    />
                </Modal>
                <Compare compareRows={selectRows} visible={isShowCompareModal} cancelFunc={this.hideCompareModal} />
            </Row>
        );
    }
}
//const Tablels:React.ComponentClass<TableListProps> = connect<any, any, TableListProps>((state,props)=>({port:state.PortReducer}) )(TableList);
export default connect<any, any, any>((state,props)=>({port:state.PortReducer, ExperimentID: state.ExperimentReducer}) )(TableList);
//export default Tablels;