import 'babel-polyfill';
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import App from './App';
import { Router, Route, browserHistory,  } from 'react-router';
import registerServiceWorker from './registerServiceWorker';
import Overview from './components/Overview';
import { Provider } from 'react-redux';
import TrialsDetail from './components/TrialsDetail';
import Search from "./Search";
import store from './store';
import './index.css';

ReactDOM.render(
    <Provider store={store} >
        <Router history={browserHistory}>
        <Route path="/" component={Search}/>
        <Route path="/project/:id/" component={App}>
                <Route path="/project/:id/oview" component={Overview} />
                <Route path="/project/:id/detail" component={TrialsDetail} />
            </Route>
        </Router>
    </Provider>
    ,
    document.getElementById('root') as HTMLElement
);
registerServiceWorker();
