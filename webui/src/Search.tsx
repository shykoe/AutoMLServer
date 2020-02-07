import * as React from 'react';
import { connect } from 'react-redux';
// import { Row, Col } from 'antd';
import axios from 'axios';
import { Input, DatePicker, Button, Table } from 'antd';
import { RouteComponentProps, Link } from 'react-router';
import { SetPort } from './actions/PortAction';
import { RangePickerValue } from 'antd/lib/date-picker/interface';

const {  RangePicker } = DatePicker;
interface TableObj {
    id: number;
    expName: string;
    runner: string;
    searchSpace: string;
    startTime: string;
    endTime: string;
    trialConcurrency: number;
    maxTrialNum: number;
    algorithmType: string;
    status: string;
    optimizeParam: string;
}
interface SearchState {
    interval: number;
    whichPageToFresh: string;
    startTime:string;
    endTime:string;
    tableData: Array<TableObj>
  }
interface SearchProps extends RouteComponentProps<any,any,{}>{
    SetPort: typeof SetPort
    DataRow: Array<TableObj>;
  }
class Search extends React.Component<SearchProps, SearchState> {
    public Input: Input|null;

    constructor(props:SearchProps) {
        super(props);
        this.state = {
          interval: 10, // sendons
          whichPageToFresh: '',
          startTime: "",
          endTime: "",
          tableData: [],
        };
      }
    componentDidMount() {
        axios(`/listExp`, {
            method: 'GET',
        }).then(res =>{
            if(res.status == 200){
                const tableData = res.data;
                this.setState({tableData:tableData});
            }
        })
    }
    searchClick = ()=>{
        if(this.Input!=null){
            console.log(this.state.startTime, this.state.endTime);
            console.log(this.Input.input.value);
        }

    }
    onChange = (date: RangePickerValue, dateString:[string, string]) => {
            console.log(date, dateString);
            this.setState({startTime:dateString[0], endTime:dateString[1]});
      }
      render(){
        const columns = [
        {
            title: 'Expriment_id',
            dataIndex: 'ExprimentId',
            key: 'ExprimentId',
            width: 140,
        },
        {
            title: 'User_name',
            dataIndex: 'UserName',
            key: 'UserName',
            width: 140,
        },
        {
            title: 'Status',
            dataIndex: 'Status',
            key: 'Status',
            width: 140,
        },
        {
            title: 'Action',
            dataIndex: 'Action',
            key: 'Action',
            width: 140,
            render:(text:string, record:TableObj)=>(
                <span>
                    <Link to={`/project/${record.id}/oview`} > View</Link>
                </span>
            )  
        }];
        // const mockData = [{
        //     id:12,
        //     expName:"test1",
        //     runner:"kwinsheng",
        //     status:"RUNNING"
        // }]
       // const {DataRow} = this.props;
        return (
            <div style={{ marginTop: 16, marginLeft: 16 }}>
                <div style={{display:"flex"}}>
                <div style={{ marginBottom: 16, width:"20em"}}>
                    <Input addonBefore="UserName" ref={input => this.Input = input} />
                </div>
                <div style={{ marginLeft: 10}}>
                    <RangePicker onChange={this.onChange} />
                </div>
                <div style={{ marginLeft: 10}}>
                    <Button icon="search" onClick={this.searchClick} >Search</Button>
                </div>
                </div>
                <div style={{ marginTop: 50}}>
                    <Table dataSource={this.state.tableData} columns={columns} />
                </div>               
            </div>

        )
      }
}
export default connect<any, any, any>((state,props)=>({port:state.PortReducer}), {SetPort: SetPort})(Search);
