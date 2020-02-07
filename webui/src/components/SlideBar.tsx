import * as React from 'react';
import { Link } from 'react-router';
import axios from 'axios';
import { connect } from 'react-redux';
import { MANAGER_IP } from '../static/const';
import MediaQuery from 'react-responsive';
import { Row, Col, Menu, Dropdown, Icon, Button, Form, Select } from 'antd';
import { FormComponentProps } from 'antd/lib/form';
import { SetExperimentID } from '../actions/ExperimentAction';
//import { SetExperimentID } from '../actions/ExperimentAction';
//import { NNILOGO } from './stateless-component/NNItabs';
const { SubMenu } = Menu;
const { Option } = Select;
const FormItem = Form.Item;
import LogDrawer from './Modal/LogDrawer';
import ExperimentDrawer from './Modal/ExperimentDrawer';
import '../static/style/slideBar.scss';
import '../static/style/button.scss';

interface SliderState {
    version: string;
    menuVisible: boolean;
    navBarVisible: boolean;
    isdisabledFresh: boolean;
    isvisibleLogDrawer: boolean;
    isvisibleExperimentDrawer: boolean;
    activeKey: string;
    port:string;
    experiments:[];
    nowExperiments:string;
}

interface SliderProps extends FormComponentProps {
    changeInterval: (value: number) => void;
    changeFresh: (value: string) => void;
    port: string;
    ExperimentID:string;
    SetExperimentID: typeof SetExperimentID
}

interface EventPer {
    key: string;
}

class SlideBar extends React.Component<SliderProps, SliderState> {

    public _isMounted = false;
    public divMenu: HTMLDivElement | null;
    public selectHTML: Select | null;

    constructor(props: SliderProps) {
        super(props);
        this.state = {
            version: '',
            menuVisible: false,
            navBarVisible: false,
            isdisabledFresh: false,
            isvisibleLogDrawer: false, // download button (nnimanager·dispatcher) click -> drawer
            isvisibleExperimentDrawer: false,
            activeKey: 'dispatcher',
            port:this.props.port,
            experiments:[],
            nowExperiments:''
        };
    }

    getNNIversion = () => {
        axios(`${MANAGER_IP}/version`, {
            method: 'GET',
            headers:{
                'upstream': this.state.port,
                'ExperimentID': this.props.ExperimentID
            }
        })
            .then(res => {
                if (res.status === 200 && this._isMounted) {
                    this.setState({ version: res.data });
                }
            });
    }

    handleMenuClick = (e: EventPer) => {
        if (this._isMounted) { this.setState({ menuVisible: false }); }
        switch (e.key) {
            // to see & download experiment parameters
            case '1':
                if (this._isMounted === true) {
                    this.setState(() => ({ isvisibleExperimentDrawer: true }));
                }
                break;
            // to see & download nnimanager log
            case '2':
                if (this._isMounted === true) {
                    this.setState(() => ({ activeKey: 'nnimanager', isvisibleLogDrawer: true }));
                }
                break;
            // to see & download dispatcher log
            case '3':
                if (this._isMounted === true) {
                    this.setState(() => ({ isvisibleLogDrawer: true, activeKey: 'dispatcher' }));
                }
                break;
            case 'close':
            case '10':
            case '20':
            case '30':
            case '60':
                this.getInterval(e.key);
                break;
            default:
        }
    }

    handleVisibleChange = ( ) => {
        if (this._isMounted === true) {
            this.setState({ menuVisible: this.state.menuVisible ? false : true });
        }
    }

    getInterval = (value: string) => {

        if (value === 'close') {
            this.props.changeInterval(0);
        } else {
            this.props.changeInterval(parseInt(value, 10));
        }
    }

    menu = () => {
        return (
            <Menu onClick={this.handleMenuClick}>
                <Menu.Item key="1">Experiment Parameters</Menu.Item>
                <Menu.Item key="2">NNImanager Logfile</Menu.Item>
                <Menu.Item key="3">Dispatcher Logfile</Menu.Item>
            </Menu>
        );
    }

    // nav bar
    navigationBar = () => {
        const { version } = this.state;
        const feedBackLink = `https://github.com/Microsoft/nni/issues/new?labels=${version}`;
        return (
            <Menu onClick={this.handleMenuClick} className="menu-list" style={{ width: 216 }}>
                {/* <Menu onClick={this.handleMenuClick} className="menu-list" style={{width: window.innerWidth}}> */}
                <Menu.Item key="feedback">
                    <a href={feedBackLink} target="_blank">Feedback</a>
                </Menu.Item>
                <SubMenu
                    key="download"
                    onTitleClick={this.handleVisibleChange}
                    title={
                        <span>
                            <span>Download</span>
                        </span>
                    }
                >
                    <Menu.Item key="1">Experiment Parameters</Menu.Item>
                    <Menu.Item key="2">NNImanager Logfile</Menu.Item>
                    <Menu.Item key="3">Dispatcher Logfile</Menu.Item>
                </SubMenu>
            </Menu>
        );
    }

    mobileTabs = () => {
        const { port } = this.state;

        return (
            // <Menu className="menuModal" style={{width: 880, position: 'fixed', left: 0, top: 56}}>
            <Menu className="menuModal" style={{ padding: '0 10px' }}>
                <Menu.Item key="overview"><Link to={`/project/${port}/oview`}>Overview</Link></Menu.Item>
                <Menu.Item key="detail"><Link to={`/project/${port}/detail`}>Trials detail</Link></Menu.Item>
            </Menu>
        );
    }

    refreshInterval = () => {
        const {
            form: { getFieldDecorator },
            // form: { getFieldDecorator, getFieldValue },
        } = this.props;
        return (
            <Form>
                <FormItem style={{ marginBottom: 0 }}>
                    {getFieldDecorator('interval', {
                        initialValue: 'Refresh every 10s',
                    })(
                        <Select onSelect={this.getInterval}>
                            <Option value="close">Disable Auto Refresh</Option>
                            <Option value="10">Refresh every 10s</Option>
                            <Option value="20">Refresh every 20s</Option>
                            <Option value="30">Refresh every 30s</Option>
                            <Option value="60">Refresh every 1min</Option>
                        </Select>,
                    )}
                </FormItem>
            </Form>
        );
    }

    select = () => {
        const { isdisabledFresh } = this.state;

        return (
            <div className="interval">
                {this.refreshInterval()}
                <Button
                    className="fresh"
                    onClick={this.fresh}
                    type="ghost"
                    disabled={isdisabledFresh}
                >
                    <Icon type="sync" /><span>Refresh</span>
                </Button>
            </div>
        );
    }
    handleChange = (value : string) =>{
        this.props.SetExperimentID(value);
        //console.log(`selected ${value}`);
        // axios(`${MANAGER_IP}/get-experiments/${value}`, {
        //     method: 'PUT',
        //     headers:{'upstream': this.state.port}
        // }).then(
        //     res => {
        //         if (res.status === 200) {
        //             this.setState({nowExperiments:value})
        //         }
        //     }

        // )
    }
    fresh = (event: React.SyntheticEvent<EventTarget>) => {
        event.preventDefault();
        event.stopPropagation();
        if (this._isMounted) {
            this.setState({ isdisabledFresh: true }, () => {
                const whichPage = window.location.pathname;
                this.props.changeFresh(whichPage);
                setTimeout(() => { this.setState(() => ({ isdisabledFresh: false })); }, 1000);
            });
        }
    }

    desktopHTML = () => {
        const { menuVisible, port } = this.state;
       
        return (
            <Row className="nav">
                <Col span={8}>
                    <span className="left-right-margin">
                    <Link to={`/project/${port}/oview`} activeClassName="high-light" className="common-tabs">
                        Overview
                    </Link>
                    </span>
                    <span>
                    <Link to={`/project/${port}/detail`} activeClassName="high-light" className="common-tabs">
                        Trials detail
                    </Link>
                    </span>
                </Col>


                <Col span={16} className="desktop-right">
                {this.state.nowExperiments !== '' && <Select defaultValue={this.state.nowExperiments} style={{ width: 300 }} onChange={this.handleChange}>
                    {this.state.experiments.map((item)=>(<Option value={item[0]}>{`ID:${item[0]} Time: ${item[1]}`} </Option>))}
                    {/* <Option value="jack">Jack</Option>
                    <Option value="lucy">Lucy</Option>
                    <Option value="Yiminghe">yiminghe</Option> */}
                </Select>}
                    <span>{this.select()}</span>
                    <span>
                        <Dropdown
                            className="dropdown"
                            overlay={this.menu()}
                            onVisibleChange={this.handleVisibleChange}
                            visible={menuVisible}
                            trigger={['click']}
                        >
                            <a className="ant-dropdown-link" href="#">
                                <Icon type="download" className="down-icon" />
                                <span>Download</span>
                                {
                                    menuVisible
                                    ?
                                    <Icon type="up" className="margin-icon"/>
                                    :
                                    <Icon type="down" className="margin-icon"/>
                                }
                            </a>
                        </Dropdown>
                    </span>

                </Col>
            </Row>
        );
    }

    tabeltHTML = () => {
        const { port } = this.state;
        return (
            <Row className="nav">
                <Col className="tabelt-left" span={14}>
                    <span>
                        <Dropdown overlay={this.navigationBar()} trigger={['click']}>
                            <Icon type="unordered-list" className="more" />
                        </Dropdown>
                    </span>
                    <span className="left-right-margin tabelt-img">
                        <Link to={`/project/${port}/oview`}>
                            <img
                                src={require('../static/img/logo2.png')}
                                alt="NNI logo"
                                style={{height: 40}}
                            />
                        </Link>
                    </span>
                    <span>
                    <Link to={`/project/${port}/oview`} activeClassName="high-light" className="common-tabs">
                        Overview
                    </Link>
                    </span>
                    <span className="left-margin">
                    <Link to={`/project/${port}/detail`} activeClassName="high-light" className="common-tabs">
                        Trials detail
                    </Link>
                    </span>
                </Col>
                <Col className="tabelt-right" span={10}>
                    {this.select()}
                </Col>
            </Row>
        );
    }

    mobileHTML = () => {
        const { isdisabledFresh } = this.state;
        return (
            <Row className="nav">
                <Col className="left" span={8}>
                    <span>
                        <Dropdown className="more-mobile" overlay={this.navigationBar()} trigger={['click']}>
                            <Icon type="unordered-list" className="more" />
                        </Dropdown>
                    </span>
                    <span>
                        <Dropdown overlay={this.mobileTabs()} trigger={['click']}>
                            <a className="ant-dropdown-link" href="#">
                                <span>NNI <Icon type="down" /></span>
                            </a>
                        </Dropdown>
                    </span>
                </Col>
                <Col className="center" span={8}>
                    <img
                        src={require('../static/img/logo2.png')}
                        alt="NNI logo"
                    />
                </Col>
                {/* the class interval have other style ! */}
                <Col className="right interval" span={8}>
                    <Button
                        className="fresh"
                        onClick={this.fresh}
                        type="ghost"
                        disabled={isdisabledFresh}
                    >
                        <Icon type="sync" /><span>Refresh</span>
                    </Button>
                </Col>
            </Row>
        );
    }
    // close log drawer (nnimanager.dispatcher)
    closeLogDrawer = () => {
        if (this._isMounted === true) {
            this.setState(() => ({ isvisibleLogDrawer: false, activeKey: '' }));
        }
    }

    // close download experiment parameters drawer
    closeExpDrawer = () => {
        if (this._isMounted === true) {
            this.setState(() => ({ isvisibleExperimentDrawer: false }));
        }
    }
    getExperiments = () => {
        axios(`${MANAGER_IP}/get-experiments/`, {
            method: 'GET',
            headers:{
                'upstream': this.state.port,
                'ExperimentID': this.props.ExperimentID
            }
        }).then(res => {
            if (res.status === 200) {
                const currentExperimentID = res.data.now;
                this.setState({nowExperiments:currentExperimentID})
                this.props.SetExperimentID(currentExperimentID);
                this.setState({experiments:res.data.data})
            }
           
        })
    }
    componentWillMount() {
        this._isMounted = true;
        this.getExperiments();
        this.getNNIversion();
        
    }

    componentWillUnmount() {
        this._isMounted = false;
    }

    render() {
        const mobile = (<MediaQuery maxWidth={884}>{this.mobileHTML()}</MediaQuery>);
        const tablet = (<MediaQuery minWidth={885} maxWidth={1241}>{this.tabeltHTML()}</MediaQuery>);
        const desktop = (<MediaQuery minWidth={1242}>{this.desktopHTML()}</MediaQuery>);
        const { isvisibleLogDrawer, activeKey, isvisibleExperimentDrawer } = this.state;
        return (
            <div>
                {mobile}
                {tablet}
                {desktop}
                {/* the drawer for dispatcher & nnimanager log message */}
                <LogDrawer
                    isVisble={isvisibleLogDrawer}
                    closeDrawer={this.closeLogDrawer}
                    activeTab={activeKey}
                />
                <ExperimentDrawer
                    isVisble={isvisibleExperimentDrawer}
                    closeExpDrawer={this.closeExpDrawer}
                />
            </div>
        );
    }
}

// export default Form.create<FormComponentProps>()(SlideBar);
export default connect<any, any, any>((state,props)=>({port:state.PortReducer, ExperimentID: state.ExperimentReducer }) ,
{SetExperimentID: SetExperimentID})(Form.create<FormComponentProps>()(SlideBar));
