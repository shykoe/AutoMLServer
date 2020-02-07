import logger from "redux-logger";
import { createStore, applyMiddleware } from'redux';
import { reducer } from './reducers';
const store = createStore(reducer, applyMiddleware(logger));
export default store;