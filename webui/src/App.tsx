import * as React from 'react';
import { connect } from 'react-redux';
import { Row, Col } from 'antd';
import { RouteComponentProps } from 'react-router';
import './App.css';
import SlideBar from './components/SlideBar';
import { SetPort } from './actions/PortAction';
interface AppState {
  interval: number;
  whichPageToFresh: string;
}
interface APPProps extends RouteComponentProps<any,any,{}>{
  SetPort: typeof SetPort
}
class App extends React.Component<APPProps, AppState> {
  public _isMounted: boolean;
  constructor(props:APPProps) {
    super(props);
    this.state = {
      interval: 10, // sendons
      whichPageToFresh: ''
    };
    const id = this.props[`params`].id
    this.props.SetPort(id);
    console.log(this.props);
    // this.props.router.push(`/${id}/oview`);
  }

  changeInterval = (interval: number) => {
    if (this._isMounted === true) {
      this.setState(() => ({ interval: interval }));
    }
  }

  changeFresh = (fresh: string) => {
    // interval * 1000 
    if (this._isMounted === true) {
      this.setState(() => ({ whichPageToFresh: fresh }));
    }
  }

  componentDidMount() {
    this._isMounted = true;
  }

  componentWillUnmount() {
    this._isMounted = false;
  }
  render() {
    const { interval, whichPageToFresh } = this.state;
    const reactPropsChildren = React.Children.map(this.props.children, child =>
      // tslint:disable-next-line:no-any
      React.cloneElement(child as React.ReactElement<any>, { interval, whichPageToFresh })
    );
    return (
      <Row className="nni" style={{ minHeight: window.innerHeight }}>
        <Row className="header">
          <Col span={1} />
          <Col className="headerCon" span={22}>
            <SlideBar changeInterval={this.changeInterval} changeFresh={this.changeFresh}/>
          </Col>
          <Col span={1} />
        </Row>
        <Row className="contentBox">
          <Row className="content">
            {reactPropsChildren}
          </Row>
        </Row>
      </Row>
    );
  }
}

export default connect<any, any, any>((state,props)=>({port:state.PortReducer}), {SetPort: SetPort})(App);
